package discovery

import (
	"fmt"
	"github.com/iain17/discovery/env"
	"github.com/iain17/discovery/pb"
	"github.com/iain17/logger"
	"github.com/golang/protobuf/proto"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

type RemoteNode struct {
	Node
	conn          net.Conn
	ln            *LocalNode
	lastHeartbeat time.Time
}

func NewRemoteNode(conn net.Conn, ln *LocalNode) *RemoteNode {
	return &RemoteNode{
		Node: Node{
			logger: logger.New(fmt.Sprintf("RemoteNode(%s)", conn.RemoteAddr().String())),
		},
		ln:            ln,
		conn:          conn,
		lastHeartbeat: time.Now(),
	}
}

func (rn *RemoteNode) sendHeartBeat() error {
	rn.logger.Debug("sending heartbeat...")
	heartbeat, err := proto.Marshal(&pb.Message{
		Version: env.VERSION,
		Msg: &pb.Message_Heartbeat{
			Heartbeat: &pb.Hearbeat{
				Message: "",
			},
		},
	})
	if err != nil {
		return err
	}
	return pb.Write(rn.conn, heartbeat)
}

func (rn *RemoteNode) Send(message string) error {
	rn.logger.Debug("sending data...")
	transfer, err := proto.Marshal(&pb.Message{
		Version: env.VERSION,
		Msg: &pb.Message_Transfer{
			Transfer: &pb.Transfer{
				Data: message,
			},
		},
	})
	if err != nil {
		return err
	}
	return pb.Write(rn.conn, transfer)
}

func (rn *RemoteNode) Close() {
	defer rn.conn.Close()
	rn.logger.Debug("closing...")
}

func (rn *RemoteNode) listen(ln *LocalNode) {
	defer func() {
		if r := recover(); r != nil {
			rn.logger.Errorf("panic: %s", r)
		}
		rn.logger.Debug("Stopping with listening.")
		rn.conn.Close()
		ln.netTableService.RemoveRemoteNode(rn.conn.RemoteAddr())
	}()
	rn.SharePeers()

	rn.logger.Debug("listening...")
	for {
		packet, err := pb.Decode(rn.conn)
		if err != nil {
			if err == io.EOF || err.Error() == "no packet read timeout" || err.Error() == "timed out waiting for ack" || err.Error() == "i/o timeout" || err.Error() == "closed" {
				break
			}
			rn.logger.Debugf("decode error, %v", err)
			continue
		}
		//rn.logger.Debugf("received, %+v", packet)

		switch packet.GetMsg().(type) {
		case *pb.Message_Heartbeat:
			rn.logger.Debug("heart beat received")
			rn.lastHeartbeat = time.Now()
			break
		case *pb.Message_Peers:
			msg := packet.GetMsg().(*pb.Message_Peers).Peers
			rn.receiveSharedPeers(msg.Peers)
			break
		case *pb.Message_PeerInfo:
			msg := packet.GetMsg().(*pb.Message_PeerInfo).PeerInfo
			rn.info = msg.Info
			break
		}
	}
}

func (rn *RemoteNode) receiveSharedPeers(peers []*pb.DPeer) {
	if len(peers) > rn.ln.discovery.max {
		rn.logger.Debug("Sent too many peers")
		return
	}
	for _, peer := range peers {
		rn.ln.netTableService.Discovered(&net.UDPAddr{
			IP:   net.ParseIP(peer.Ip),
			Port: int(peer.Port),
		})
	}
}

func (rn *RemoteNode) SharePeers() error {
	var peers []*pb.DPeer
	for _, peer := range rn.ln.netTableService.GetPeers() {
		//Don't share himself to him.
		if rn.conn.RemoteAddr().String() == peer.conn.RemoteAddr().String() {
			continue
		}
		ipPort := strings.Split(peer.conn.RemoteAddr().String(), ":")
		if len(ipPort) != 2 {
			rn.logger.Warning("Weird peer found %s", peer.conn.RemoteAddr().String())
			continue
		}
		port, err := strconv.Atoi(ipPort[1])
		if err != nil {
			rn.logger.Warning(err)
			continue
		}
		peers = append(peers, &pb.DPeer{
			Ip:   ipPort[0],
			Port: int32(port),
		})
		//Enough?
		if len(peers) >= rn.ln.discovery.max {
			break
		}
	}

	peerShare, err := proto.Marshal(&pb.Message{
		Version: env.VERSION,
		Msg: &pb.Message_Peers{
			Peers: &pb.Peers{
				Peers: peers,
			},
		},
	})
	if err != nil {
		return err
	}
	return pb.Write(rn.conn, peerShare)
}

func (rn *RemoteNode) String() string {
	return fmt.Sprintf("Remote node(%s) with info: %#v", rn.conn.RemoteAddr().String(), rn.info)
}
