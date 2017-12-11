package api

//import (
//	"github.com/iain17/decentralizer/pb"
//	"time"
//)

//func (s *Serve) handleUpsertSessionRequest(msg *pb.RPCMessage) (*pb.RPCMessage, error) {
//	request := msg.GetUpsertSessionRequest()
//	if request.Info.Details == nil {
//		request.Info.Details = make(map[string]string)
//	}
//	request.Info.Details["updated"] = time.Now().String()
//	sessId, err := s.app.UpsertSession(request.Info.Type, request.Info.Name, request.Info.Port, request.Info.Details)
//	return &pb.RPCMessage{
//		Msg: &pb.RPCMessage_UpsertSessionResponse{
//			UpsertSessionResponse: &pb.RPCUpsertSessionResponse{
//				Result: err == nil,
//				SessionId: sessId,
//			},
//		},
//	}, err
//}
//
//func (s *Serve) handleDeleteSessionRequest(msg *pb.RPCMessage) (*pb.RPCMessage, error) {
//	request := msg.GetDeleteSessionRequest()
//	err := s.app.DeleteSession(request.SessionId)
//	return &pb.RPCMessage{
//		Msg: &pb.RPCMessage_DeleteSessionResponse{
//			DeleteSessionResponse: &pb.RPCDeleteSessionResponse{
//				Result: err == nil,
//			},
//		},
//	}, err
//}
//
//func (s *Serve) handleSessionIdsRequest(msg *pb.RPCMessage) (*pb.RPCMessage, error) {
//	request := msg.GetSessionIdsRequest()
//	var sessions []*pb.SessionInfo
//	var err error
//	if request.Key == "" && request.Value == "" {
//		sessions, err = s.app.GetSessions(request.Type)
//	} else {
//		sessions, err = s.app.GetSessionsByDetails(request.Type, request.Key, request.Value)
//	}
//	var sessionIds []uint64
//	for _, session := range sessions {
//		sessionIds = append(sessionIds, session.SessionId)
//	}
//	return &pb.RPCMessage{
//		Msg: &pb.RPCMessage_SessionIdsResponse{
//			SessionIdsResponse: &pb.RPCSessionIdsResponse{
//				SessionIds: sessionIds,
//			},
//		},
//	}, err
//}
//
//func (s *Serve) handleGetSessionRequest(msg *pb.RPCMessage) (*pb.RPCMessage, error) {
//	request := msg.GetGetSessionRequest()
//	session, err := s.app.GetSession(request.SessionId)
//	return &pb.RPCMessage{
//		Msg: &pb.RPCMessage_GetSessionResponse{
//			GetSessionResponse: &pb.RPCGetSessionResponse{
//				Found: err != nil,
//				Result: session,
//			},
//		},
//	}, err
//}