package api

import (
	"github.com/iain17/decentralizer/app"
	"github.com/iain17/decentralizer/pb"
	"net"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	app *app.Decentralizer
}

func New(app *app.Decentralizer, port int) (*Server, error) {
	i := &Server {
		app: app,
	}
	address := fmt.Sprintf(":%d", port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}
	s := grpc.NewServer()
	pb.RegisterDecentralizerServer(s, i)
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		return nil, err
	}

	return i, nil
}