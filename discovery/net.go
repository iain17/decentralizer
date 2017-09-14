package discovery

import (
	"time"
	"github.com/anacrolix/utp"
	"github.com/iain17/decentralizer/discovery/pb"
)

//Initiate a hand sake procedure.
//See (l *ListenerService) process(c net.Conn) error for the receiving side.
func connect(h string, ln *LocalNode) (*RemoteNode, error) {
	s, errSocket := utp.NewSocket("udp4", ":0")
	if errSocket != nil {
		return nil, errSocket
	}

	conn, errDial := s.DialTimeout(h, 10 * time.Second)
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