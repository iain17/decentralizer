package serve

import (
	"github.com/iain17/decentralizer/pb"
	"github.com/iain17/logger"
)

func (s *Serve) handleUpsertSessionRequest(msg *pb.RPCMessage) (*pb.RPCMessage, error) {
	logger.Info("handleUpsertSessionRequest")
	return nil, nil
}

func (s *Serve) handleDeleteSessionRequest(msg *pb.RPCMessage) (*pb.RPCMessage, error) {
	logger.Info("handleDeleteSessionRequest")
	return nil, nil
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