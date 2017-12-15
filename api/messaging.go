package api

import (
	"context"
	"github.com/iain17/decentralizer/pb"
)

//
// Messaging
//
// Send another peer a direct message. Takes a peer id and the data it should send
func (s *Server) SendDirectMessage(ctx context.Context, request *pb.RPCDirectMessageRequest) (*pb.RPCDirectMessageResponse, error) {
	err := s.app.SendMessage(request.PId, request.Message)
	return &pb.RPCDirectMessageResponse{}, err
}