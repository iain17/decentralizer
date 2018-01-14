package app

import (
	"github.com/iain17/logger"
	"gx/ipfs/QmXYjuNuxVzXKJCfWasQk1RqkhVLDM9jtUKhqc2WPQmFSB/go-libp2p-peer"
	inet "gx/ipfs/QmNa31VPzC561NWwRsJLE7nGYZYuuD2QfpK2b1q9BK54J1/go-libp2p-net"
	"github.com/iain17/framed"
	"github.com/iain17/decentralizer/pb"
	gogoProto "gx/ipfs/QmZ4Qi3GaRbjcx28Sme5eMH7RQjGkt8wHxt2a65oLaeFEV/gogo-protobuf/proto"
	"encoding/binary"
	"bytes"
	"fmt"
	"io"
)

type sessionRequest struct{
	id   peer.ID
	search *search
}

//On top of DHT we will query our providers for a larger list. This is done with streams
//This is because DHT will only allow one value for one key. hashmap, duh. But this means a new user won't receive the whole list
//in one go as we want. By querying the peers that are giving these session from DHT (already connected etc) we can fetch the bigger list.
func (d *Decentralizer) initMatchmakingStream() {
	d.sessionQueries 			= make(chan sessionRequest, CONCURRENT_SESSION_REQUEST)
	d.i.PeerHost.SetStreamHandler(GET_SESSION_REQ, d.getSessionResponse)
	//Spawn some workers
	logger.Debugf("Running %d session request workers", CONCURRENT_SESSION_REQUEST)
	for i := 0; i < CONCURRENT_SESSION_REQUEST; i++ {
		go d.processSessionRequest()
	}
}

func (d *Decentralizer) processSessionRequest() {
	for {
		select {
		case <-d.ctx.Done():
			return
		case req, ok := <-d.sessionQueries:
			logger.Infof("Querying %s for sessions of type: %d", req.id.Pretty(), req.search.sessionType)
			if !ok {
				return
			}
			d.getSessionsRequest(req.id, req.search)
		}
	}
}

//Receive a request to give our sessions for a certain type
func (d *Decentralizer) getSessionResponse(stream inet.Stream) {
	defer stream.Close()
	//Receive session type
	resData, err := framed.Read(stream)
	var sessionType uint64
	binary.Read(bytes.NewReader(resData), binary.LittleEndian, &sessionType)

	if !d.hasSessionStorage(sessionType) {
		err = fmt.Errorf("we don't have session type %d", sessionType)
		logger.Warning(err)
		return
	}

	logger.Infof("%s requested sessions of type %d...", stream.Conn().RemotePeer().Pretty(), sessionType)
	sessionsStorage := d.getSessionStorage(sessionType)
	rawIds := sessionsStorage.SessionIds()
	logger.Infof("Sending %d sessions back of type %d", len(rawIds), sessionType)
	buf := new(bytes.Buffer)//TODO: Allocate already the max buffer size. bytes.NewBuffer(make([]byte, 8 * len(rawIds) - 1))
	for _, sessionId := range rawIds {
		id := sessionId.(uint64)
		logger.Debugf("Boasting that we've got %d to %s", id, stream.Conn().RemotePeer().Pretty())
		err = binary.Write(buf, binary.LittleEndian, id)
		if err != nil {
			logger.Error(err)
			return
		}
	}
	framed.Write(stream, buf.Bytes())

	//Wait for a session request. This is done by sending one uint64
	var sessionId uint64
	for resData, err = framed.Read(stream); err == nil; {
		binary.Read(bytes.NewReader(resData), binary.LittleEndian, &sessionId)
		session, err := sessionsStorage.FindSessionId(sessionId)
		if err != nil {
			err = fmt.Errorf("%s requested session %d which we didn't have: %s", stream.Conn().RemotePeer().Pretty(), sessionId, err.Error())
			logger.Warning(err)
			continue
		}
		response, err := gogoProto.Marshal(session)
		err = framed.Write(stream, response)
		if err != nil {
			logger.Error(err)
			return
		}
		logger.Debugf("Responded back to %s with session %d", stream.Conn().RemotePeer().Pretty(), sessionId)
	}
}

// Get in contact with a peer and ask it for session of a certain type
func (d *Decentralizer) getSessionsRequest(peer peer.ID, search *search) error {
	stream, err := d.i.PeerHost.NewStream(d.i.Context(), peer, GET_SESSION_REQ)
	if err != nil {
		return err
	}
	defer stream.Close()
	logger.Debugf("Requesting %s for any sessions", peer.Pretty())
	//Request available session ids
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.LittleEndian, search.sessionType)
	if err != nil {
		return err
	}
	err = framed.Write(stream, buf.Bytes())
	if err != nil {
		return err
	}

	//Response
	resData, err := framed.Read(stream)
	if err != nil {
		return err
	}
	idsResponse := bytes.NewReader(resData)

	//Request missing session ids
	store := search.storage
	seen := map[uint64]bool{}
	for err == nil {
		var sessionId uint64
		binary.Read(idsResponse, binary.LittleEndian, &sessionId)
		if sessionId == 0 {
			logger.Warning("Received stop sign.")
			break
		}
		if seen[sessionId] {
			logger.Warningf("%d sessionId dup boasting detected.", sessionId)
			break
		}
		seen[sessionId] = true
		if !store.Contains(sessionId) {
			logger.Debugf("Missing session %d. Asking %s for it", sessionId, stream.Conn().RemotePeer().Pretty())
			session, err := d.requestSessionId(stream, sessionId)
			if err != nil {
				err = fmt.Errorf("failed to receive %d from %s: %s", sessionId, stream.Conn().RemotePeer().Pretty(), err.Error())
				logger.Warning(err)
				break
			}
			if session.PId == d.i.Identity.Pretty() {
				err = fmt.Errorf("we will not allow another peer to dictate our sessions. skipping")
				logger.Warning(err)
				continue
			}
			d.setSessionIdToType(session.SessionId, session.Type)
			search.storage.Insert(session)
		}
	}
	return nil
}

func (d *Decentralizer) requestSessionId(w io.ReadWriter, sessionId uint64) (*pb.Session, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, sessionId)
	if err != nil {
		return nil, err
	}
	err = framed.Write(w, buf.Bytes())
	if err != nil {
		return nil, err
	}
	resData, err := framed.Read(w)
	if err != nil {
		return nil, err
	}
	var session pb.Session
	err = d.unmarshal(resData, &session)
	return &session, err
}