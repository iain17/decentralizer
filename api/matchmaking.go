package api

import (
	"github.com/iain17/decentralizer/pb"
	"time"
	"context"
	"github.com/iain17/logger"
)

//
// Matchmaking
//
// Create or update a session. Takes session info, returns session id.
func (s *Server) UpsertSession(ctx context.Context, request *pb.RPCUpsertSessionRequest) (*pb.RPCUpsertSessionResponse, error) {
	logger.Infof("Upsert session request received")
	if request.Session.Details == nil {
		request.Session.Details = make(map[string]string)
	}
	request.Session.Details["updated"] = time.Now().String()
	sessId, err := s.app.UpsertSession(request.Session.Type, request.Session.Name, request.Session.Port, request.Session.Details)
	return &pb.RPCUpsertSessionResponse{
		SessionId: sessId,
	}, err
}

// Delete a session. Takes session id, returns bool informing if the deletion was a success
func (s *Server) DeleteSession(ctx context.Context, request *pb.RPCDeleteSessionRequest) (*pb.RPCDeleteSessionResponse, error) {
	logger.Infof("Delete session request received")
	err := s.app.DeleteSession(request.SessionId)
	return &pb.RPCDeleteSessionResponse{
		Result: err == nil,
	}, err
}

// Get session ids. Takes session type, and a key and value to filter the sessions by details. If left empty this filter will not apply  and all will be fetched.
func (s *Server) GetSessionIds(ctx context.Context, request *pb.RPCGetSessionIdsRequest) (*pb.RPCGetSessionIdsResponse, error) {
	logger.Infof("Get session ids request received")
	var sessions []*pb.Session
	var err error
	if request.Key == "" && request.Value == "" {
		sessions, err = s.app.GetSessions(request.Type)
	} else {
		sessions, err = s.app.GetSessionsByDetails(request.Type, request.Key, request.Value)
	}
	var sessionIds []uint64
	for _, session := range sessions {
		sessionIds = append(sessionIds, session.SessionId)
	}
	return &pb.RPCGetSessionIdsResponse{
		SessionIds: sessionIds,
	}, err
}

// Get an individual session. Takes session id and returns session info.
func (s *Server) GetSession(ctx context.Context, request *pb.RPCGetSessionRequest) (*pb.RPCGetSessionResponse, error) {
	logger.Infof("Get session request received")
	session, err := s.app.GetSession(request.SessionId)
	return &pb.RPCGetSessionResponse{
		Session: session,
	}, err
}