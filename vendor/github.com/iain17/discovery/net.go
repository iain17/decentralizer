package discovery

import (
	"github.com/iain17/discovery/pb"
	"net"
	"time"
	"fmt"
)

//Initiate a handshake procedure.
//See (l *ListenerService) process(c net.Conn) error for the receiving side.
func connect(h *net.UDPAddr, ln *LocalNode) (*RemoteNode, error) {
	accepted := false

	conn, errDial := ln.listenerService.socket.DialContext(ln.discovery.ctx, "udp4", h.String())
	defer func() {
		//s.Close()
		if !accepted && conn != nil{
			conn.Close()
		}
	}()
	if errDial != nil {
		return nil, fmt.Errorf("error dialing %s: %s", h.String(), errDial.Error())
	}
	conn.SetDeadline(time.Now().Add(1 * time.Second))

	rn := NewRemoteNode(conn, ln)

	//Handshake dance.
	//Wait for them to accept and send their peer info
	rn.logger.Debug("Waiting for their peer info...")
	peerInfo, err := pb.DecodePeerInfo(rn.conn, string(ln.discovery.network.ExportPublicKey()))
	if err != nil {
		return nil, fmt.Errorf("error at waiting for their peer info: %s", err.Error())
	}
	rn.logger.Debug("Received peer info...")
	rn.Initialize(peerInfo)

	conn.SetDeadline(time.Now().Add(5 * time.Second))
	rn.logger.Debug("Sending our peer info")
	err = ln.sendPeerInfo(rn.conn)
	if err != nil {
		return nil, fmt.Errorf("error sending our peer info: %s", err.Error())
	}

	rn.logger.Info("connected!")
	accepted = true
	return rn, nil
}