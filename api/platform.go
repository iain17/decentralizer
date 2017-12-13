package api

import (
	"github.com/iain17/decentralizer/pb"
	"context"
	"github.com/iain17/logger"
)

//
// Platform
//
func (s *Server) GetHealth(ctx context.Context, in *pb.RPCHealthRequest) (*pb.RPCHealthReply, error) {
	logger.Info("Getting health..")
	return &pb.RPCHealthReply{
		Ready: true,
		Message: "cool",
	}, nil
	/*
	ready, err := s.app.Health()
	var error string
	if err != nil {
		error = err.Error()
	}
	return &pb.RPCHealthReply{
		Ready: ready,
		Message: error,
	}, nil
	*/
}