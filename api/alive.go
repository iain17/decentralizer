//Takes care of asserting if services are still using the this background service. If its not being used, close.
package api

import (
	"google.golang.org/grpc"
	"context"
	"time"
	"os"
)

func (s *Server) AliveStreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		s.wg.Add(1)
		go func() {
			<- stream.Context().Done()
			s.wg.Done()
		}()
		return handler(srv, stream)
	}
}

func (s *Server) AliveUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		s.wg.Add(1)
		defer s.wg.Done()
		return handler(ctx, req)
	}
}

func (s *Server) RunAlive() {
	var free time.Time
	for {
		time.Sleep(MAX_IDLE_TIME)
		s.wg.Wait()
		free = time.Now()
		s.wg.Wait()
		if free.Add(MAX_IDLE_TIME).After(time.Now()) {
			s.Stop()
			os.Exit(0)
		}
	}
}