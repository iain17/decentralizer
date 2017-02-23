package decentralizer

import (
	"net"
	logger "github.com/Sirupsen/logrus"
)

/*
- The introductory server is a UDP server that is used to exchange the details of a service.
 */
func (d *decentralizer) listenIntroServer() error {
	conn, host, err := getUdpConn()
	if err != nil {
		return err
	}
	port := conn.LocalAddr().(*net.UDPAddr).Port
	d.introPort = uint16(port)
	d.ip = host.IP()
	logger.Infof("Intro server listening at %d", port)

	//todo: Do something with the udp server?!?!
	return nil
}