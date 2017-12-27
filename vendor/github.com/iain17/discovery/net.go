package discovery

import (
	"github.com/anacrolix/utp"
	"github.com/iain17/discovery/pb"
	"net"
)

//Initiate a handshake procedure.
//See (l *ListenerService) process(c net.Conn) error for the receiving side.
func connect(h *net.UDPAddr, ln *LocalNode) (*RemoteNode, error) {
	s, errSocket := utp.NewSocket("udp", ":0")
	if errSocket != nil {
		return nil, errSocket
	}
	conn, errDial := s.Dial(h.String())
	if errDial != nil {
		return nil, errDial
	}
	rn := NewRemoteNode(conn, ln)

	//Handshake dance.
	rn.logger.Debug("Sending our peer info")
	ln.sendPeerInfo(rn.conn)

	//They will respond by sending their peer info
	rn.logger.Debug("Waiting for their peer info...")
	peerInfo, err := pb.DecodePeerInfo(rn.conn, string(ln.discovery.network.ExportPublicKey()))
	if err != nil {
		rn.logger.Debug(err)
		conn.Close()
		return nil, err
	}
	rn.logger.Debug("Received peer info...")
	rn.Initialize(peerInfo)

	rn.logger.Info("connected!")
	return rn, nil
}