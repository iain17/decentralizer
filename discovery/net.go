package discovery

import (
	"github.com/anacrolix/utp"
	"github.com/iain17/decentralizer/discovery/pb"
	"net"
)

//Initiate a hand sake procedure.
//See (l *ListenerService) process(c net.Conn) error for the receiving side.
func connect(h *net.UDPAddr, ln *LocalNode) (*RemoteNode, error) {
	s, errSocket := utp.NewSocket("udp4", ":0")
	if errSocket != nil {
		return nil, errSocket
	}

	conn, errDial := s.Dial(h.String())
	if errDial != nil {
		return nil, errDial
	}
	rn := NewRemoteNode(conn)

	//We start by sending our heartbeat.
	rn.sendHeartBeat()

	//They will respond by sending their peer info
	peerInfo, errPeerInfo := pb.DecodePeerInfo(rn.conn)
	if errPeerInfo != nil {
		return nil, errPeerInfo
	}
	rn.info = peerInfo.Info

	//We send our peer peer info back
	ln.sendPeerInfo(rn.conn)

	rn.logger.Info("connected!")
	return rn, nil
}