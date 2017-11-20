package app

import (
	"fmt"
	"github.com/iain17/decentralizer/app/pb"
	"github.com/iain17/decentralizer/utils"
	"github.com/iain17/logger"
	"errors"
	"time"
	"github.com/iain17/decentralizer/app/sessionstore"
)

func getKey(sessionType uint64) string {
	return fmt.Sprintf("%d", sessionType)
}

func (d *Decentralizer) getSessionStorage(sessionType uint64) *sessionstore.Store {
	if d.sessions[sessionType] == nil {
		var err error
		d.sessions[sessionType], err = sessionstore.New(1000, time.Duration((EXPIRE_TIME_SESSION * 1.5) * time.Second))
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
	d.i.Routing.PutValue(d.i.Context(), getKey(info.Type), []byte("Everything in this world is magic... except to the magician."))
	sessionId, err := sessions.Insert(info)
	d.sessionIdToSessionType[sessionId] = sessionType
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
	sessions := d.getSessionStorage(sessionType)
	return sessions.FindByDetails(key, value)
}

func (d *Decentralizer) refreshSessions(sessionType uint64) {
	answers, err := d.i.Routing.GetValues(d.i.Context(), getKey(sessionType), 1)
	if err != nil {
		logger.Warning(err)
		return
	}
	for _, answer := range answers {
		logger.Infof("Peer %s has possible sessions", answer.From)
	}
}


/*
rawSessions, err := d.i.Routing.GetValues(d.i.Context(), strconv.FormatUint(sessionType, 16), 0)
	if err != nil {
		return nil, err
	}
	var result []*pb.SessionInfo
	for _, rawSession := range rawSessions {
		var session pb.SessionInfo
		err = proto.Unmarshal(rawSession.Val, &session)
		if err != nil {
			logger.Debug(err)
			continue
		}
		result = append(result, &session)
	}
	return result, nil
 */