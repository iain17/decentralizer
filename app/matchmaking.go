package app

import (
	"fmt"
	"github.com/iain17/decentralizer/app/pb"
	"github.com/iain17/decentralizer/utils"
	"github.com/iain17/logger"
	"errors"
	"github.com/golang/protobuf/proto"
	"strconv"
)

func getSessionId(info *pb.SessionInfo) uint64 {
	return info.Type+uint64(info.Port)+info.DId
}

func getKey(sessionType uint64) string {
	return fmt.Sprintf("MATCHMAKING_%d", sessionType)
}

func (d *Decentralizer) UpsertSession(sessionType uint64, name string, port uint32, details map[string]string) (uint64, error) {
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
	info.SessionId = getSessionId(info)
	d.sessions[info.SessionId] = info
	return info.SessionId, d.setSession(info)
}

func (d *Decentralizer) setSession(info *pb.SessionInfo) error {
	msg, err := proto.Marshal(info)
	if err != nil {
		return err
	}
	return d.i.Routing.PutValue(d.i.Context(), getKey(info.Type), msg)
}

func (d *Decentralizer) DeleteSession(sessionId uint64) error {
	if d.sessions[sessionId] == nil {
		return errors.New("no such session exists")
	}
	delete(d.sessions, sessionId)
	return d.i.Routing.PutValue(d.i.Context(), getKey(d.sessions[sessionId].Type), nil)
}

func (d *Decentralizer) GetSession(sessionType uint64) (*pb.SessionInfo, error) {
	data, err := d.i.Routing.GetValue(d.i.Context(), getKey(sessionType))
	if err != nil {
		return nil, err
	}
	var session pb.SessionInfo
	err = proto.Unmarshal(data, &session)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (d *Decentralizer) GetSessions(sessionType uint64) ([]*pb.SessionInfo, error) {
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
}