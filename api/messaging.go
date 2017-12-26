package api

import (
	"context"
	"github.com/iain17/decentralizer/pb"
	"time"
	"errors"
)

//
// Messaging
//
// Send another peer a direct message. Takes a peer id and the data it should send
func (s *Server) SendDirectMessage(ctx context.Context, request *pb.RPCDirectMessage) (*pb.Empty, error) {
	err := s.app.SendMessage(request.Channel, request.PId, request.Message)
	return &pb.Empty{}, err
}

func (s *Server) ReceiveDirectMessage(request *pb.RPCReceiveDirectMessageRequest, stream pb.Decentralizer_ReceiveDirectMessageServer) error {
	if s.listeningChannels[request.Channel] {
		return errors.New("another instance is already listening on this channel")
	}
	s.listeningChannels[request.Channel] = true
	messageChannel := s.app.GetMessagingChan(request.Channel)
	for {
		select {
		case <-stream.Context().Done():
			return nil
		case msg, _ := <- messageChannel:
			stream.Send(&pb.RPCDirectMessage{
				PId: msg.PId,
				Message: msg.Message,
			})
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}