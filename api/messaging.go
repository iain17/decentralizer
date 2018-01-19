package api

import (
	"context"
	"github.com/iain17/decentralizer/pb"
	"errors"
	"github.com/iain17/logger"
)

//
// Messaging
//
// Send another peer a direct message. Takes a peer id and the data it should send
func (s *Server) SendDirectMessage(ctx context.Context, request *pb.RPCDirectMessage) (*pb.Empty, error) {
	err := s.App.SendMessage(request.Channel, request.PId, request.Message)
	return &pb.Empty{}, err
}

func (s *Server) ReceiveDirectMessage(request *pb.RPCReceiveDirectMessageRequest, stream pb.Decentralizer_ReceiveDirectMessageServer) error {
	s.mutex.Lock()
	if s.listeningChannels[request.Channel] {
		err := errors.New("another instance is already listening on this channel")
		logger.Warning(err)
		return err
	}
	s.listeningChannels[request.Channel] = true
	s.mutex.Unlock()

	defer func() {
		s.mutex.Lock()
		s.listeningChannels[request.Channel] = false
		s.mutex.Unlock()
	}()

	messageChannel := s.App.GetMessagingChan(request.Channel)
	for {
		select {
		case <-stream.Context().Done():
			return nil
		case msg, _ := <- messageChannel:
			stream.Send(&pb.RPCDirectMessage{
				PId: msg.PId,
				Message: msg.Message,
			})
		}
	}
}