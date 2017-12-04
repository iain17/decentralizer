package serve

import (
	"github.com/iain17/decentralizer/serve/pb"
)

func (s *Serve) handleHealthRequest(msg *pb.RPCMessage) (*pb.RPCMessage, error) {
	ready, err := s.app.Health()
	return &pb.RPCMessage{
		Id: msg.Id,
		Msg: &pb.RPCMessage_HealthReply{
			HealthReply: &pb.HealthReply{
				Ready: ready,
				Message: err.Error(),
			},
		},
	}, nil
}