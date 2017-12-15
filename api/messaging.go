package api

import (
	"context"
	"github.com/iain17/decentralizer/pb"
	"errors"
)

//
// Messaging
//
// Send another peer a direct message. Takes a peer id and the data it should send
func (s *Server) SendDirectMessage(context.Context, *pb.RPCDirectMessageRequest) (*pb.RPCDirectMessageResponse, error) {
	return nil, errors.New("Unimplemented")
}