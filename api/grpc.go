package api

import (
	"github.com/iain17/decentralizer/pb"
	"github.com/iain17/decentralizer/vars"
	"net"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"github.com/iain17/logger"
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"errors"
	"google.golang.org/grpc/metadata"
	"time"
)

func (s *Server) initGRPC(port int) error {
	lis, err := net.Listen("tcp", s.endpoint)
	if err != nil {
		return err
	}
	logger.Infof("Serving GRPC API on: %s", s.endpoint)
	s.grpc = grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_auth.StreamServerInterceptor(s.auth),
			grpc_recovery.StreamServerInterceptor(),
			//nil,
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_auth.UnaryServerInterceptor(s.auth),
			grpc_recovery.UnaryServerInterceptor(),
			//nil,
		)),
	)
	pb.RegisterDecentralizerServer(s.grpc, s)
	// Register reflection service on gRPC server.
	reflection.Register(s.grpc)
	if err := s.grpc.Serve(lis); err != nil {
		return err
	}
	return nil
}

func (s *Server) auth(ctx context.Context) (context.Context, error) {
	s.LastCall = time.Now()

	var clientVersion, networkKey string
	var isPrivateKey, limitedConnection bool
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx, errors.New("set context pls")
	}
	if len(meta["cver"]) != 0 {
		clientVersion = meta["cver"][0]
	}
	if len(meta["netkey"]) != 0 {
		networkKey = meta["netkey"][0]
	}
	if len(meta["privkey"]) != 0 {
		isPrivateKey = meta["privkey"][0] == "1"
	}
	if len(meta["limited"]) != 0 {
		limitedConnection = meta["limited"][0] == "1"
	}
	if s.App == nil && networkKey != "" {
		logger.Info("Joining network...")
		err := s.SetNetwork(clientVersion, networkKey, isPrivateKey, limitedConnection)
		if err != nil {
			logger.Warningf("Failed to join network: %v", err)
			vars.DEFAULT_PORT = vars.DEFAULT_PORT + 1//Try incrementing the default port. To fix Windows editions
			return ctx, err
		} else {
			logger.Info("Joined network.")
		}
	}
	if s.App == nil {
		return ctx, errors.New("network is not set. Please set the network first")
	}

	//Check health
	if len(meta["health"]) == 0 {
		ready, _, err := s.App.Health(false)
		if err != nil {
			return ctx, err
		}
		if !ready {
			return ctx, errors.New("not ready yet. check health check")
		}
	}

	return ctx, nil
}

