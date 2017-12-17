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

func (s *Server) setNetwork(clientVersion string, networkKey string, isPrivateKey bool) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.app != nil {
		return errors.New("network already set")
	}
	v, err := version.NewVersion(clientVersion)
	if err != nil {
		return err
	}
	versionMismatch := pb.CONSTRAINT.Check(v)
	if !versionMismatch {
		return errors.New("please update your client")
	}
	logger.Info("Starting IPFS...")
	s.app, err = app.New(s.ctx, networkKey, isPrivateKey)
	if err == nil {
		logger.Info("IPFS started")
	}
	return err
}