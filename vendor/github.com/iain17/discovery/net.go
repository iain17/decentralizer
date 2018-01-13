package discovery

import (
	"github.com/iain17/discovery/pb"
	"net"
	"time"
)

//Initiate a handshake procedure.
//See (l *ListenerService) process(c net.Conn) error for the receiving side.
func connect(h *net.UDPAddr, ln *LocalNode) (*RemoteNode, error) {
	accepted := false

	conn, errDial := ln.listenerService.socket.DialContext(ln.discovery.ctx, h.String())
	defer func() {
		//s.Close()
		if !accepted && conn != nil{
			conn.Close()
		}
	}()
	if errDial != nil {
		return nil, errDial
	}
	conn.SetDeadline(time.Now().Add(300 * time.Millisecond))

	rn := NewRemoteNode(conn, ln)

	//Handshake dance.
	rn.logger.Debug("Sending our peer info")
	ln.sendPeerInfo(rn.conn)

	//They will respond by sending their peer info
	rn.logger.Debug("Waiting for their peer info...")
	peerInfo, err := pb.DecodePeerInfo(rn.conn, string(ln.discovery.network.ExportPublicKey()))
	if err != nil {
		rn.logger.Debug(err)
		return nil, err
	}
	rn.logger.Debug("Received peer info...")
	rn.Initialize(peerInfo)

	rn.logger.Info("connected!")
	accepted = true
	return rn, nil
}