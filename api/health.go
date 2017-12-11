package api

//import (
//	"github.com/iain17/decentralizer/pb"
//	"github.com/iain17/logger"
//)
//
//func (s *Serve) handleHealthRequest(msg *pb.RPCMessage) (*pb.RPCMessage, error) {
//	logger.Info("handleHealthRequest")
//	ready, err := s.app.Health()
//	logger.Info("handleHealthRequest done")
//	var error string
//	if err != nil {
//		error = err.Error()
//	}
//	return &pb.RPCMessage{
//		Msg: &pb.RPCMessage_HealthReply{
//			HealthReply: &pb.RPCHealthReply{
//				Ready: ready,
//				Message: error,
//			},
//		},
//	}, nil
//}