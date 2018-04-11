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
	"github.com/iain17/framed"
	"sync"
)

type RemoteNode struct {
	Node
	conn          net.Conn
	closed		  bool
	ln            *LocalNode
	lastHeartbeat time.Time
	mutex		  sync.Mutex
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

func (rn *RemoteNode) GetIp() net.IP {
	parts := strings.Split(rn.conn.RemoteAddr().String(), ":")
	if len(parts) != 2 {
		return net.IP{}
	}
	return net.ParseIP(parts[0])
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
	return framed.Write(rn.conn, heartbeat)
}

func (rn *RemoteNode) Send(data []byte) error {
	rn.logger.Debug("sending data...")
	transfer, err := proto.Marshal(&pb.Message{
		Version: env.VERSION,
		Msg: &pb.Message_Transfer{
			Transfer: &pb.Transfer{
				Data: data,
			},
		},
	})
	if err != nil {
		return err
	}
	return framed.Write(rn.conn, transfer)
}

func (rn *RemoteNode) Close() error {
	rn.mutex.Lock()
	defer rn.mutex.Unlock()
	if rn.closed {
		return nil
	}
	defer rn.conn.Close()
	rn.closed = true
	rn.logger.Debug("Closing")
	rn.ln.netTableService.seen.Remove(rn.conn.RemoteAddr().String())
	transfer, err := proto.Marshal(&pb.Message{
		Version: env.VERSION,
		Msg: &pb.Message_Shutdown{
			Shutdown: &pb.Shutdown{},
		},
	})
	if err != nil {
		return err
	}
	return framed.Write(rn.conn, transfer)
}

func (rn *RemoteNode) listen() {
	defer func() {
		if r := recover(); r != nil {
			rn.logger.Errorf("[%s]: %s", rn.id, r)
		}
		rn.logger.Debug("Stopping with listening.")
		rn.conn.Close()
		rn.ln.netTableService.RemoveRemoteNode(rn)
	}()
	rn.SharePeers()

	rn.logger.Debug("listening...")
	i := 0
	for {
		rn.conn.SetDeadline(time.Now().Add((HEARTBEAT_DELAY * 1.5) * time.Second))
		packet, err := pb.Decode(rn.conn)
		if err != nil {
			rn.logger.Debugf("decode error, %v", err)
			if err == io.EOF || err.Error() == "no packet read timeout" || err.Error() == "timed out waiting for ack" || err.Error() == "i/o timeout" || err.Error() == "closed" {
				break
			}
			if i > 10 {
				break
			}
			i++
			continue
		}
		//rn.logger.Debugf("received, %+v", packet)

		switch packet.GetMsg().(type) {
		case *pb.Message_Heartbeat:
			rn.logger.Debug("heart beat received")
			rn.lastHeartbeat = time.Now()
			break
		case *pb.Message_Shutdown:
			rn.conn.Close()
			return
		case *pb.Message_Peers:
			msg := packet.GetMsg().(*pb.Message_Peers).Peers
			rn.receiveSharedPeers(msg.Peers)
			break
		case *pb.Message_PeerInfo:
			msg := packet.GetMsg().(*pb.Message_PeerInfo).PeerInfo
			rn.SetMapInfo(msg.Info)
			break
		case *pb.Message_Transfer:
			data := packet.GetMsg().(*pb.Message_Transfer).Transfer.Data
			rn.ln.ReceivedMessage <- ExternalMessage{
				from: rn,
				data: data,
			}
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
	return framed.Write(rn.conn, peerShare)
}

func (rn *RemoteNode) String() string {
	return fmt.Sprintf("Remote node(%s) with info: %#v", rn.conn.RemoteAddr().String(), rn.info)
}
