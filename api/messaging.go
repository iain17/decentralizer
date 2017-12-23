package api

import (
	"context"
	"github.com/iain17/decentralizer/pb"
	"time"
)

//
// Messaging
//
// Send another peer a direct message. Takes a peer id and the data it should send
func (s *Server) SendDirectMessage(ctx context.Context, request *pb.RPCDirectMessage) (*pb.Empty, error) {
	err := s.app.SendMessage(request.PId, request.Message)
	return &pb.Empty{}, err
}

func (s *Server) ReceiveDirectMessage(request *pb.Empty, stream pb.Decentralizer_ReceiveDirectMessageServer) error {
	for {
		select {
		case <-stream.Context().Done():
			return nil
		case msg, _ := <- s.app.DirectMessage:
			stream.Send(&pb.RPCDirectMessage{
				PId: msg.PId,
				Message: msg.Message,
			})
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}