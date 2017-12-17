package api

import (
	"github.com/iain17/decentralizer/app"
	"github.com/iain17/decentralizer/pb"
	"net"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"github.com/iain17/logger"
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"errors"
	"google.golang.org/grpc/metadata"
)

type Server struct {
	ctx context.Context
	app *app.Decentralizer
	grpc *grpc.Server
}

func New(ctx context.Context, port int) (*Server, error) {
	i := &Server {
		ctx: ctx,
	}
	address := fmt.Sprintf(":%d", port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}
	logger.Infof("Serving API on: %s", address)
	i.grpc = grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_auth.StreamServerInterceptor(i.auth),
			grpc_recovery.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_auth.UnaryServerInterceptor(i.auth),
			grpc_recovery.UnaryServerInterceptor(),
		)),
	)
	pb.RegisterDecentralizerServer(i.grpc, i)
	// Register reflection service on gRPC server.
	reflection.Register(i.grpc)
	if err := i.grpc.Serve(lis); err != nil {
		return nil, err
	}

	return i, nil
}

func (s *Server) Stop() {
	s.grpc.GracefulStop()
}

func (s *Server) auth(ctx context.Context) (context.Context, error) {
	//Check health
	_, err := s.app.Health()
	if err != nil {
		return ctx, err
	}

	var clientVersion string
	var networkKey string
	var isPrivateKey bool
	meta, ok := metadata.FromIncomingContext(ctx)
	if ! ok {
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
	if s.app == nil && networkKey != "" {
		err := s.SetNetwork(clientVersion, networkKey, isPrivateKey)
		if err != nil {
			return ctx, err
		}
	}
	if s.app == nil {
		return ctx, errors.New("network is not set. Please set the network first")
	}
	return ctx, nil
}

