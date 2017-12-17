package api

import (
	"fmt"
	"google.golang.org/grpc"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"net/http"
	"github.com/iain17/decentralizer/pb"
	"github.com/iain17/logger"
)

func (s *Server) initHTTP(port int) error {
	mux := runtime.NewServeMux()
	address := fmt.Sprintf(":%d", port)
	logger.Infof("Serving HTTP API on: %s", address)
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := pb.RegisterDecentralizerHandlerFromEndpoint(s.ctx, mux, s.endpoint, opts)
	if err != nil {
		return err
	}
	return http.ListenAndServe(address, mux)
}
