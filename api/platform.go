package api

import (
	"github.com/iain17/decentralizer/pb"
	"context"
	"github.com/iain17/logger"
	"github.com/iain17/decentralizer/app"
	"github.com/hashicorp/go-version"
	"errors"
	"strings"
)

//
// Platform
//
func (s *Server) GetHealth(ctx context.Context, in *pb.RPCHealthRequest) (*pb.RPCHealthReply, error) {
	logger.Info("Getting health..")
	ready, numConns, err := s.App.Health(in.WaitForMinConnections)
	var error string
	if err != nil {
		error = err.Error()
	}
	return &pb.RPCHealthReply{
		Ready: ready,
		Message: error,
		BasePath: app.Base.Path,
		NumConnections: uint32(numConns),
	}, nil
}

func (s *Server) SetNetwork(clientVersion string, networkKey string, isPrivateKey bool, limitedConnection bool) error {
	s.mutex.Lock()
	defer func () {
		s.mutex.Unlock()
	}()
	if s.App != nil {
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
	make:
	s.App, err = app.New(s.ctx, networkKey, isPrivateKey, limitedConnection, s.profile)
	if err != nil && strings.Contains(err.Error(), "corrupted") {
		logger.Warningf("%s: Resetting...", err)
		app.Reset()
		goto make
	}
	return err
}