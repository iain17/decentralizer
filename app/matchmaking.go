package app

import (
	"fmt"
	"github.com/iain17/decentralizer/app/pb"
	"github.com/iain17/decentralizer/utils"
	"errors"
	"time"
	"github.com/iain17/decentralizer/app/sessionstore"
	"github.com/giantswarm/retry-go"
	inet "gx/ipfs/QmahYsGWry85Y7WUe2SX5G4JkH2zifEQAUtJVLZ24aC9DF/go-libp2p-net"
	peer "gx/ipfs/QmXYjuNuxVzXKJCfWasQk1RqkhVLDM9jtUKhqc2WPQmFSB/go-libp2p-peer"
	"sync"
	"github.com/gogo/protobuf/proto"
	"io"
	"github.com/iain17/logger"
)

func getKey(sessionType uint64) string {
	return fmt.Sprintf("%d", sessionType)
}

func (d *Decentralizer) initMatchmaking() {
	d.i.PeerHost.SetStreamHandler(GET_SESSION_REQ, d.getSessionResponse)
}

func (d *Decentralizer) getSessionStorage(sessionType uint64) *sessionstore.Store {
	if d.sessions[sessionType] == nil {
		var err error
		d.sessions[sessionType], err = sessionstore.New(MAX_SESSIONS, time.Duration((EXPIRE_TIME_SESSION * 1.5) * time.Second))
		if err != nil {
			return nil
		}
	}
	return d.sessions[sessionType]
}

func (d *Decentralizer) UpsertSession(sessionType uint64, name string, port uint32, details map[string]string) (uint64, error) {
	sessions := d.getSessionStorage(sessionType)
	pId, dId := pb.GetPeer(d.i.Identity)
	info := &pb.SessionInfo{
		DId: dId,
		PId: pId,
		Type: sessionType,
		Name: name,
		Address: utils.Inet_aton(d.d.GetIP()),
		Port: port,
		Details: details,
	}
	sessionId, err := sessions.Insert(info)
	if err != nil {
		return 0, err
	}
	d.sessionIdToSessionType[sessionId] = sessionType

	err = d.b.Provide(getKey(sessionType))
	if err != nil {
		return 0, err
	}
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

func (d *Decentralizer) GetSessions(sessionType uint64, key, value string) ([]*pb.SessionInfo, error) {
	retry.Do(func() error {
		d.refreshSessions(sessionType)
		return nil
	}, retry.Timeout(10 * time.Second))
	sessions := d.getSessionStorage(sessionType)
	return sessions.FindByDetails(key, value)
}

func (d *Decentralizer) refreshSessions(sessionType uint64) {
	var wg sync.WaitGroup
	sessionsStorage := d.getSessionStorage(sessionType)
	for peer := range d.b.Find(getKey(sessionType), MAX_SESSIONS) {
		wg.Add(1)
		go func() {
			sessions, err := d.getSessions(peer, sessionType)
			if err != nil {
				logger.Error(err)
				return
			}
			for _, session := range sessions {
				sessionId, err := sessionsStorage.Insert(session)
				if err != nil {
					return
				}
				d.sessionIdToSessionType[sessionId] = sessionType
			}
		}()
	}
	wg.Wait()
}

func (d *Decentralizer) getSessionResponse(stream inet.Stream)  {
	reqData := make([]byte, MAX_SIZE)
	_, err := io.ReadFull(stream, reqData)
	if err != nil {
		logger.Error(err)
		return
	}
	var request pb.SessionRequest
	err = proto.Unmarshal(reqData, &request)
	if err != nil {
		logger.Error(err)
		return
	}
	sessionsStorage := d.getSessionStorage(request.Type)
	sessions, err := sessionsStorage.FindByPeerId(d.i.Identity.Pretty())
	if err != nil {
		logger.Error(err)
		return
	}

	//Response
	response, err := proto.Marshal(&pb.SessionResponse{
		Results: sessions,
	})
	_, err = stream.Write(response)
	if err != nil {
		logger.Error(err)
		return
	}
}

func (d *Decentralizer) getSessions(peer peer.ID, sessionType uint64) ([]*pb.SessionInfo, error) {
	stream, err := d.i.PeerHost.NewStream(d.i.Context(), peer, GET_SESSION_REQ)
	if err != nil {
		return nil, err
	}
	//Request
	reqData, err := proto.Marshal(&pb.SessionRequest{
		Type: sessionType,
	})
	if err != nil {
		return nil, err
	}
	_, err = stream.Write(reqData)
	if err != nil {
		return nil, err
	}

	//Response
	resData := make([]byte, MAX_SIZE)
	_, err = io.ReadFull(stream, resData)
	if err != nil {
		return nil, err
	}
	var response pb.SessionResponse
	err = proto.Unmarshal(resData, &response)
	if err != nil {
		return nil, err
	}
	return response.Results, nil
}