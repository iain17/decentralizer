package decentralizer

import (
	"net"
	logger "github.com/Sirupsen/logrus"
	"github.com/iain17/dht-hello/decentralizer/pb"
	"github.com/gogo/protobuf/proto"
)

/*
- The introductory server is a UDP server that is used to exchange the details of a service.
 */
func (d *decentralizer) setupIntroServer() error {
	conn, host, err := getUdpConn()
	if err != nil {
		return err
	}
	port := conn.LocalAddr().(*net.UDPAddr).Port
	d.introPort = uint16(port)
	d.ip = host.IP()
	logger.Infof("Intro server listening at %d", port)
	d.introConn = conn
	go d.listenIntroServer(d.introConn)
	return nil
}

func (d *decentralizer) listenIntroServer(conn *net.UDPConn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		n, addr, err := conn.ReadFromUDP(buf)
		request := &pb.IntroductionRequest{}
		abc := buf[0:n]
		logger.Infof("%s", abc)
		err = proto.Unmarshal(abc, request)
		if err != nil {
			logger.Warn(err)
			continue
		}
		logger.Infof("Received request %s from %s", request.Hash, addr)
	}
}