package serve

import (
	"net"
	"github.com/iain17/logger"
	"fmt"
	"github.com/iain17/decentralizer/serve/pb"
	"reflect"
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

	for {
		packet, err := pb.Decode(conn)
		if err != nil {
			logger.Warning(err)
			continue
		}

		handler := s.handlers[reflect.TypeOf(packet.GetMsg())]
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