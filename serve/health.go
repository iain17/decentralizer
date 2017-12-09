package serve

import (
	"github.com/iain17/decentralizer/serve/pb"
	"github.com/iain17/logger"
)

func (s *Serve) handleHealthRequest(msg *pb.RPCMessage) (*pb.RPCMessage, error) {
	logger.Info("handleHealthRequest:", msg.GetHealthRequest().Msg)
	ready, err := s.app.Health()
	return &pb.RPCMessage{
		Msg: &pb.RPCMessage_HealthReply{
			HealthReply: &pb.HealthReply{
				Ready: ready,
				Message: err.Error(),
			},
		},
	}, nil
}