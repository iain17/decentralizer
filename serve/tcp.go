package serve

import (
	"net"
	"github.com/iain17/logger"
	"fmt"
	"github.com/iain17/decentralizer/serve/pb"
)

func (s *Serve) ListenTCP(port int) {
	address := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		logger.Error(err)
		return
	}
	defer listener.Close()
	logger.Infof("TCP API serving on %s", address)

	for {
		conn, err := listener.Accept()
		logger.Infof("New connection: %s", conn.RemoteAddr())
		if err != nil {
			logger.Warning(err)
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *Serve) handleConnection(conn net.Conn) {
	defer func() {
		if error := recover(); error != nil {
			logger.Errorf("Recover error: %s", error)
		}

		conn.Close()
	}()

	demo1(conn)
	demo2(conn)

	for {
		packet, err := pb.Decode(conn)
		if err != nil {
			logger.Warning(err)
			break
		}

		handler := s.handlers[pb.MessageType(packet.Type)]
		if handler != nil {
			res, err := handler(packet)
			if err != nil {
				logger.Warning(err)
				continue
			}
			if res != nil {
				err = pb.Write(conn, res)
				if err != nil {
					logger.Warning(err)
					continue
				}
			}
		}
	}
}

func demo1(conn net.Conn) {
	pb.Write(conn, &pb.RPCMessage{
		Id: 1337,
		Msg: &pb.RPCMessage_HealthReply{
			HealthReply: &pb.HealthReply{
				Ready: true,
				Message: "very nice...",
			},
		},
	})
	logger.Error("Sent a demo 1 message: ")
}

func demo2(conn net.Conn) {
	pb.Write(conn, &pb.RPCMessage{
		Id: 1338,
		Msg: &pb.RPCMessage_HealthReply{
			HealthReply: &pb.HealthReply{
				Ready: true,
				Message: "very nice...",
			},
		},
	})
	logger.Error("Sent a demo 2 message: ")
}