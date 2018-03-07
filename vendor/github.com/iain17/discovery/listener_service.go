package discovery

import (
	"github.com/anacrolix/utp"
	"net"
	"context"
	"github.com/iain17/discovery/pb"
	"github.com/iain17/logger"
	"errors"
	"fmt"
	"os"
	"strings"
)

type ListenerService struct {
	localNode *LocalNode
	context context.Context
	listener	  net.PacketConn
	socket    *utp.Socket

	logger *logger.Logger
}

func (d *ListenerService) String() string {
	return "listener"
}

func (l *ListenerService) init(ctx context.Context) error {
	defer func() {
		if l.localNode.wg != nil {
			l.localNode.wg.Done()
		}
		if l.localNode.coreWg != nil {
			l.localNode.coreWg.Done()
			l.localNode.coreWg = nil
		}
	}()
	l.logger = logger.New(l.String())
	l.context = ctx

	var err error
	l.socket, err = utp.NewSocket("udp4", ":0")
	if err != nil {
		return fmt.Errorf("could not listen: %s", err.Error())
	}
	addr := l.socket.Addr().(*net.UDPAddr)
	l.localNode.port = addr.Port
	l.logger.Infof("listening on %d", l.localNode.port)

	//Stun Disabled: Just fucks with the initial connections.
	//go func() {
	//	stunErr := l.localNode.StunService.Serve(ctx)
	//	if stunErr != nil {
	//		logger.Warningf("Stun error: %s", stunErr)
	//	}
	//}()
	return err
}

func (l *ListenerService) Serve(ctx context.Context) {
	defer l.Stop()
	if err := l.init(ctx); err != nil {
		l.localNode.lastError = err
		panic(err)
	}

	for {
		select {
		case <-l.context.Done():
			return
		default:
			conn, err := l.socket.Accept()
			if err != nil {
				logger.Warning(err)
				if opErr, ok := err.(*net.OpError); ok {
					if scErr, ok := opErr.Err.(*os.SyscallError); ok && strings.Contains(scErr.Error(), "keep-alive") {
						return
					}
				}
				break
			}
			key := conn.RemoteAddr().String()
			if _, ok := l.localNode.netTableService.blackList.Get(key); ok {
				conn.Close()
				return
			}

			go func(conn net.Conn) {
				l.logger.Debugf("new connection from %s", conn.RemoteAddr().String())

				if err = l.process(conn); err != nil {
					conn.Close()
					if err.Error() == "peer reset" || err.Error() == "we can't add ourselves" {
						return
					}
					l.logger.Warningf("[%s] error on processing new connection, %s", conn.RemoteAddr().String(), err)
				}
			}(conn)
		}
	}
}

func (l *ListenerService) Stop() {
	l.logger.Info("Stopping...")
	l.socket.Close()
}

//We receive a connection from a possible new peer.
func (l *ListenerService) process(c net.Conn) error {
	rn := NewRemoteNode(c, l.localNode)

	rn.logger.Debug("Sending our peer info...")
	err := l.localNode.sendPeerInfo(c)
	if err != nil {
		return err
	}
	rn.logger.Debug("Sent our peer info...")

	rn.logger.Debug("Waiting for peer info...")
	peerInfo, err := pb.DecodePeerInfo(c, string(l.localNode.discovery.network.ExportPublicKey()))
	if err != nil {
		return err
	}
	if peerInfo.Id == l.localNode.id {
		return errors.New("we can't add ourselves")
	}
	if l.localNode.netTableService.isConnected(peerInfo.Id) {
		logger.Debugf("We are already connected to %s", peerInfo.Id)
		return nil
	}
	rn.logger.Debug("Received peer info...")

	rn.Initialize(peerInfo)
	l.localNode.netTableService.AddRemoteNode(rn)
	return nil
}
