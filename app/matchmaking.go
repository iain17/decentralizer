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
	inet "gx/ipfs/QmU4vCDZTPLDqSDKguWbHCiUe46mZUtmM2g2suBZ9NE8ko/go-libp2p-net"
	peer "gx/ipfs/QmWNY7dV54ZDYmTA1ykVdwNCqC11mpU4zSUp6XDpLTH9eG/go-libp2p-peer"
	"time"
	"encoding/hex"
	"github.com/iain17/timeout"
	"context"
)

func (d *Decentralizer) getMatchmakingKey(sessionType uint64) string {
	ih := d.n.InfoHash()
	return fmt.Sprintf("%s_MATCHMAKING_%d", hex.EncodeToString(ih[:]), sessionType)
}

func (d *Decentralizer) initMatchmaking() {
	d.i.PeerHost.SetStreamHandler(GET_SESSION_REQ, d.getSessionResponse)
}

func (d *Decentralizer) getSessionStorage(sessionType uint64) *sessionstore.Store {
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

func (d *Decentralizer) getSessionSearch(sessionType uint64) *search {
	if d.searches[sessionType] == nil {
		var err error
		d.searches[sessionType], err = d.newSearch(d.i.Context(), sessionType)
		if err != nil {
			logger.Warningf("Could not start session search: %v", err)
			return nil
		}
	}
	return d.searches[sessionType]
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
	for _, search := range d.searches {
		storage := search.fetch()
		peers, err := storage.FindByPeerId(peerId)
		if err != nil {
			logger.Warning(err)
			continue
		}
		result = append(result, peers...)
	}
	return result, nil
}

func (d *Decentralizer) getSessionResponse(stream inet.Stream) {
	reqData, err := Read(stream)
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
	err = Write(stream, response)
	if err != nil {
		logger.Error(err)
		return
	}
}

func (d *Decentralizer) getSessionsRequest(peer peer.ID, sessionType uint64) ([]*pb.Session, error) {
	stream, err := d.i.PeerHost.NewStream(d.i.Context(), peer, GET_SESSION_REQ)
	if err != nil {
		return nil, err
	}
	stream.SetDeadline(time.Now().Add(300 * time.Millisecond))
	defer stream.Close()
	//Request
	reqData, err := proto.Marshal(&pb.DNSessionRequest{
		Type: sessionType,
	})
	if err != nil {
		return nil, err
	}
	err = Write(stream, reqData)
	if err != nil {
		return nil, err
	}

	//Response
	resData, err := Read(stream)
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
