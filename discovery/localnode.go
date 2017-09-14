package discovery

import (
	"github.com/iain17/decentralizer/network"
	"context"
	"github.com/iain17/freeport"
	"io"
	"github.com/golang/protobuf/proto"
	"github.com/iain17/decentralizer/discovery/pb"
	"github.com/op/go-logging"
)

type LocalNode struct {
	Node
	network *network.Network
	port int
	//Services
	listenerService ListenerService
	upNpService UPnPService
	netTableService NetTableService
	//Peer discoveries
	discoveryDHT  DiscoveryDHT
}

func NewLocalNode(ctx context.Context, network *network.Network) (*LocalNode, error) {
	instance := &LocalNode{
		Node: Node{
			logger: logging.MustGetLogger("LocalNode"),
		},
		network: network,
		port: freeport.GetUDPPort(),
	}
	err := instance.listenerService.Init(ctx, instance)
	if err != nil {
		return nil, err
	}

	err = instance.netTableService.Init(ctx, instance)
	if err != nil {
		return nil, err
	}

	err = instance.upNpService.Init(ctx, instance)
	if err != nil {
		return nil, err
	}

	err = instance.discoveryDHT.Init(ctx, instance)
	if err != nil {
		return nil, err
	}

	return instance, nil
}

func (ln *LocalNode) sendPeerInfo(w io.Writer) error {
	peerInfo, err := proto.Marshal(&pb.PeerInfo{
		Info: ln.info,
	})
	if err != nil {
		return err
	}
	return ln.write(w, pb.TransferMessage, peerInfo)
}

func (rn *LocalNode) write(w io.Writer, messageType pb.MessageType, data []byte) error {
	rn.logger.Debug("sending message...")
	packet := pb.NewPacket(messageType, data)
	err := packet.Write(w)
	if err != nil {
		return err
	}
	rn.logger.Debug("message sent")
	return nil
}