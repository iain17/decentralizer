package app

import (
	"errors"
	"fmt"
	"github.com/iain17/decentralizer/app/peerstore"
	"github.com/iain17/decentralizer/app/sessionstore"
	"github.com/iain17/decentralizer/pb"
	"github.com/iain17/decentralizer/utils"
	"github.com/iain17/logger"
	"gx/ipfs/QmT6n4mspWYEya864BhCUJEgyxiRfmiSY9ruQwTUNpRKaM/protobuf/proto"
	inet "gx/ipfs/QmNa31VPzC561NWwRsJLE7nGYZYuuD2QfpK2b1q9BK54J1/go-libp2p-net"
	peer "gx/ipfs/QmXYjuNuxVzXKJCfWasQk1RqkhVLDM9jtUKhqc2WPQmFSB/go-libp2p-peer"
	"time"
	"encoding/hex"
	"github.com/iain17/timeout"
	"context"
	"github.com/iain17/framed"
	"github.com/Akagi201/kvcache/ttlru"
)

type sessionRequest struct{
	peer peer.ID
	sessionType uint64
}

func (d *Decentralizer) getMatchmakingKey(sessionType uint64) string {
	ih := d.n.InfoHash()
	return fmt.Sprintf("%s_MATCHMAKING_%d", hex.EncodeToString(ih[:]), sessionType)
}

func (d *Decentralizer) initMatchmaking() {
	d.sessionQueries 			= make(chan sessionRequest, CONCURRENT_SESSION_REQUEST)
	d.sessions 					= make(map[uint64]*sessionstore.Store)
	d.sessionIdToSessionType	= make(map[uint64]uint64)
	var err error
	d.searches, err = lru.NewTTL(MAX_SESSION_SEARCHES)
	if err != nil {
		logger.Fatal(err)
	}
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
			logger.Info("Received query request")
			if !ok {
				return
			}
			search := d.getSessionSearch(req.sessionType)
			sessions, err := d.getSessionsRequest(req.peer, req.sessionType)
			if err != nil {
				logger.Debug(err)
				if err.Error() == "i/o deadline reached" {
					continue
				}
				if err.Error() == "protocol not supported" {
					d.ignore.Add(req.peer, true)
					continue
				}
				logger.Debugf("Failed to get sessions from %s: %v", req.peer, err)
			} else {
				search.add(sessions, req.peer)
			}
		}
	}
}

func (d *Decentralizer) getSessionStorage(sessionType uint64) *sessionstore.Store {
	d.matchmakingMutex.Lock()
	defer d.matchmakingMutex.Unlock()
	if d.sessions[sessionType] == nil {
		var err error
		d.sessions[sessionType], err = sessionstore.New(MAX_SESSIONS, time.Duration((EXPIRE_TIME_SESSION*1.5)*time.Second), d.i.Identity)
		if err != nil {
			logger.Warningf("Could not get session storage: %v", err)
			return nil
		}
	}
	return d.sessions[sessionType]
}

func (d *Decentralizer) getSessionSearch(sessionType uint64) (result *search) {
	d.searchMutex.Lock()
	defer d.searchMutex.Unlock()
	if !d.searches.Contains(sessionType) {
		var err error
		result, err = d.newSearch(d.i.Context(), sessionType)
		d.searches.Add(sessionType, result)
		if err != nil {
			logger.Warningf("Could not start session search: %v", err)
			return nil
		}
	} else {
		value, ok := d.searches.Get(sessionType)
		if !ok {
			return nil
		}
		result = value.(*search)
	}
	return result
}

func (d *Decentralizer) UpsertSession(sessionType uint64, name string, port uint32, details map[string]string) (uint64, error) {
	sessions := d.getSessionStorage(sessionType)
	pId, dId := peerstore.PeerToDnId(d.i.Identity)
	info := &pb.Session{
		DnId:    dId,
		PId:     pId,
		Type:    sessionType,
		Name:    name,
		Address: uint32(utils.Inet_aton(d.GetIP())),
		Port:    port,
		Details: details,
	}
	//Also look for other sessions. So we can become a provider for more than just ourselves.
	go d.getSessionSearch(sessionType)
	sessionId, err := sessions.Insert(info)
	if err != nil {
		return 0, err
	}
	d.sessionIdToSessionType[sessionId] = sessionType
	timeout.Do(func(ctx context.Context) {
		d.b.Provide(d.getMatchmakingKey(sessionType))
	}, 3*time.Second)
	return sessionId, err
}

func (d *Decentralizer) setSessionIdToType(sessionId uint64, sessionType uint64) {
	d.matchmakingMutex.Lock()
	defer d.matchmakingMutex.Unlock()
	d.sessionIdToSessionType[sessionId] = sessionType
}

func (d *Decentralizer) DeleteSession(sessionId uint64) error {
	if d.sessionIdToSessionType[sessionId] == 0 {
		return errors.New("no such session exists")
	}
	sessionType := d.sessionIdToSessionType[sessionId]
	sessions := d.getSessionStorage(sessionType)
	return sessions.Remove(sessionId)
}

func (d *Decentralizer) GetSession(sessionId uint64) (*pb.Session, error) {
	if d.sessionIdToSessionType[sessionId] == 0 {
		return nil, errors.New("no such session exists")
	}
	sessionType := d.sessionIdToSessionType[sessionId]
	sessions := d.getSessionStorage(sessionType)
	return sessions.FindSessionId(sessionId)
}

func (d *Decentralizer) GetSessions(sessionType uint64) ([]*pb.Session, error) {
	search := d.getSessionSearch(sessionType)
	if search != nil {
		storage := search.fetch()
		return storage.FindAll()
	}
	return nil, errors.New("could not get session search")
}

func (d *Decentralizer) GetSessionsByDetails(sessionType uint64, key, value string) ([]*pb.Session, error) {
	search := d.getSessionSearch(sessionType)
	if search != nil {
		storage := search.fetch()
		return storage.FindByDetails(key, value)
	}
	return nil, errors.New("could not get session search")
}

func (d *Decentralizer) GetSessionsByPeer(peerId string) ([]*pb.Session, error) {
	var result []*pb.Session
	for _, key := range d.searches.Keys() {
		value, ok := d.searches.Get(key)
		if ! ok {
			continue
		}
		storage := value.(*search).fetch()
		peers, err := storage.FindByPeerId(peerId)
		if err != nil {
			logger.Warning(err)
			continue
		}
		result = append(result, peers...)
	}
	return result, nil
}

//Receive a request to give our sessions for a certain type
func (d *Decentralizer) getSessionResponse(stream inet.Stream) {
	reqData, err := framed.Read(stream)
	if err != nil {
		logger.Error(err)
		return
	}
	var request pb.DNSessionRequest
	err = proto.Unmarshal(reqData, &request)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Infof("Someone requested our sessions of type %d...", request.Type)
	sessionsStorage := d.getSessionStorage(request.Type)
	sessions, err := sessionsStorage.FindByPeerId(d.i.Identity.Pretty())
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Infof("Sending %d sessions back of type %d", len(sessions), request.Type)
	//Response
	response, err := proto.Marshal(&pb.DNSessionResponse{
		Results: sessions,
	})
	if err != nil {
		logger.Error(err)
		return
	}
	err = framed.Write(stream, response)
	if err != nil {
		logger.Error(err)
		return
	}
}

// Get in contact with a peer and ask it for session of a certain type
func (d *Decentralizer) getSessionsRequest(peer peer.ID, sessionType uint64) ([]*pb.Session, error) {
	stream, err := d.i.PeerHost.NewStream(d.i.Context(), peer, GET_SESSION_REQ)
	if err != nil {
		return nil, err
	}
	stream.SetDeadline(time.Now().Add(1 * time.Second))
	defer stream.Close()
	logger.Debugf("Requesting %s for any sessions", peer.Pretty())
	//Request
	reqData, err := proto.Marshal(&pb.DNSessionRequest{
		Type: sessionType,
	})
	if err != nil {
		return nil, err
	}
	err = framed.Write(stream, reqData)
	if err != nil {
		return nil, err
	}

	//Response
	resData, err := framed.Read(stream)
	if err != nil {
		return nil, err
	}
	var response pb.DNSessionResponse
	err = proto.Unmarshal(resData, &response)
	if err != nil {
		return nil, err
	}
	return response.Results, nil
}
