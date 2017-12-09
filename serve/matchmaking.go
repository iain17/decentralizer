package serve

import (
	"github.com/iain17/decentralizer/pb"
	"github.com/iain17/logger"
	"time"
)

func (s *Serve) handleUpsertSessionRequest(msg *pb.RPCMessage) (*pb.RPCMessage, error) {
	request := msg.GetUpsertSessionRequest()
	if request.Info.Details == nil {
		request.Info.Details = make(map[string]string)
	}
	request.Info.Details["updated"] = time.Now().String()
	sessId, err := s.app.UpsertSession(request.Info.Type, request.Info.Name, request.Info.Port, request.Info.Details)
	return &pb.RPCMessage{
		Msg: &pb.RPCMessage_UpsertSessionResponse{
			UpsertSessionResponse: &pb.RPCUpsertSessionResponse{
				Result: err == nil,
				SessionId: sessId,
			},
		},
	}, err
}

func (s *Serve) handleDeleteSessionRequest(msg *pb.RPCMessage) (*pb.RPCMessage, error) {
	request := msg.GetDeleteSessionRequest()
	err := s.app.DeleteSession(request.SessionId)
	return &pb.RPCMessage{
		Msg: &pb.RPCMessage_DeleteSessionResponse{
			DeleteSessionResponse: &pb.RPCDeleteSessionResponse{
				Result: err == nil,
			},
		},
	}, err
}

func (s *Serve) handleRefreshSessionsRequest(msg *pb.RPCMessage) (*pb.RPCMessage, error) {
	logger.Info("handleRefreshSessionsRequest")
	return nil, nil
}

func (s *Serve) handleSessionIdsRequest(msg *pb.RPCMessage) (*pb.RPCMessage, error) {
	logger.Info("handleSessionIdsRequest")
	return nil, nil
}

func (s *Serve) handleGetSessionRequest(msg *pb.RPCMessage) (*pb.RPCMessage, error) {
	logger.Info("handleGetSessionRequest")
	return nil, nil
}