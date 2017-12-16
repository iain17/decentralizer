package api

import (
	"github.com/iain17/decentralizer/pb"
	"context"
	"github.com/iain17/logger"
	"github.com/iain17/decentralizer/app"
	"github.com/hashicorp/go-version"
	"errors"
)

//
// Platform
//
func (s *Server) GetHealth(ctx context.Context, in *pb.RPCHealthRequest) (*pb.RPCHealthReply, error) {
	logger.Info("Getting health..")
	ready, err := s.app.Health()
	var error string
	if err != nil {
		error = err.Error()
	}
	return &pb.RPCHealthReply{
		Ready: ready,
		Message: error,
	}, nil
}

func (s *Server) SetNetwork(ctx context.Context, request *pb.RPCSetNetworkRequest) (*pb.RPCSetNetworkResponse, error) {
	if s.app != nil {
		return nil, errors.New("network already set")
	}
	clientVersion, err := version.NewVersion(request.ClientVersion)
	if err != nil {
		return nil, err
	}
	versionMismatch := pb.CONSTRAINT.Check(clientVersion)
	if !versionMismatch {
		return nil, errors.New("please update your client")
	}
	s.app, err = app.New(s.ctx, request.NetworkKey, request.PrivateKey)
	if err != nil {
		return nil, err
	}
	return &pb.RPCSetNetworkResponse{
		Version: pb.VERSION.String(),
	}, nil
}