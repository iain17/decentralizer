package api

import (
	"github.com/iain17/decentralizer/app"
	"google.golang.org/grpc"
	"github.com/iain17/logger"
	"context"
	"fmt"
	"sync"
)

type Server struct {
	ctx context.Context
	app *app.Decentralizer
	grpc *grpc.Server
	endpoint string
	mutex sync.Mutex
	listeningChannels map[uint32]bool//To keep track if a client is already listening for direct messages on this channel.
}

func New(ctx context.Context, port int) (*Server, error) {
	i := &Server {
		ctx: ctx,
		endpoint: fmt.Sprintf(":%d", port),
		listeningChannels: make(map[uint32]bool),
	}
	go func() {
		err := i.initGRPC(port)
		if err != nil {
			logger.Errorf("GRPC API error: %s", err)
		}
	}()
	go func() {
		err := i.initHTTP(port + 1)
		if err != nil {
			logger.Errorf("HTTP API error: %s", err)
		}
	}()

	return i, nil
}


func (s *Server) Stop() {
	s.grpc.Stop()
	s.app.Stop()
}