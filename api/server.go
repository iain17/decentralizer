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
)

type Server struct {
	ctx context.Context
	app *app.Decentralizer
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
	logger.Infof("Serving GRPC API on: %s", address)
	s := grpc.NewServer()
	pb.RegisterDecentralizerServer(s, i)
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		return nil, err
	}

	return i, nil
}