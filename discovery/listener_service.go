package discovery

import (
	"github.com/anacrolix/utp"
	"net"
	"github.com/op/go-logging"
	"context"
	"fmt"
	"github.com/iain17/decentralizer/discovery/pb"
)

type ListenerService struct {
	localNode *LocalNode
	context context.Context
	socket    *utp.Socket

	logger *logging.Logger
}

func (l *ListenerService) Init(ctx context.Context, ln *LocalNode) error {
	l.logger = logging.MustGetLogger("listener")
	l.localNode = ln
	l.context = ctx

	//Open a new socket on a free UDP port.
	var err error
	l.logger.Infof("listening on %d", ln.port)
	l.socket, err = utp.NewSocket("udp4", fmt.Sprintf("0.0.0.0:%d", ln.port))
	go l.Run()
	return err
}

func (l *ListenerService) Run() {
	defer l.Stop()
	for {
		select {
		case <-l.context.Done():
			return
		default:
			conn, err := l.socket.Accept()
			if err != nil {
				break
			}

			l.logger.Debugf("new connection from %q", conn.RemoteAddr().String())

			if err = l.process(conn); err != nil {
				l.logger.Errorf("error on process, %v", err)
			}
		}
	}
}

func (l *ListenerService) Stop() {
	l.socket.Close()
}

//We receive a connection from a possible new peer.
func (l *ListenerService) process(c net.Conn) error {
	rn := NewRemoteNode(c)
	err := rn.sendHeartBeat()
	if err != nil {
		return err
	}

	rn.logger.Debug("Waiting for heartbeat...")
	err = pb.DecodeHeartBeat(c)
	if err != nil {
		return err
	}
	rn.logger.Debug("Received the heartbeat...")

	rn.logger.Debug("Sending our peer info...")
	err = l.localNode.sendPeerInfo(c)
	if err != nil {
		return err
	}
	rn.logger.Debug("Sent our peer info...")

	rn.logger.Debug("Waiting for peer info...")
	peerInfo, err := pb.DecodePeerInfo(c, string(l.localNode.discovery.network.ExportPublicKey()))
	if err != nil {
		return err
	}
	rn.logger.Debug("Received peer info...")

	rn.Info = peerInfo.Info
	l.localNode.netTableService.AddRemoteNode(rn)
	return nil
}
