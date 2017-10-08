package discovery

import (
	"github.com/iain17/freeport"
	"io"
	"github.com/golang/protobuf/proto"
	"github.com/iain17/decentralizer/discovery/pb"
	"github.com/op/go-logging"
	"github.com/iain17/decentralizer/discovery/env"
)

type LocalNode struct {
	Node
	discovery *Discovery
	port int
	//Services
	listenerService ListenerService
	upNpService UPnPService
	netTableService NetTableService
	//Peer discoveries
	discoveryDHT  DiscoveryDHT
}

func newLocalNode(discovery *Discovery) (*LocalNode, error) {
	instance := &LocalNode{
		Node: Node{
			logger: logging.MustGetLogger("LocalNode"),
			Info: map[string]string{},
		},
		discovery: discovery,
		port: freeport.GetUDPPort(),
	}
	err := instance.listenerService.Init(discovery.ctx, instance)
	if err != nil {
		return nil, err
	}

	err = instance.netTableService.Init(discovery.ctx, instance)
	if err != nil {
		return nil, err
	}

	err = instance.upNpService.Init(discovery.ctx, instance)
	if err != nil {
		return nil, err
	}

	err = instance.discoveryDHT.Init(discovery.ctx, instance)
	if err != nil {
		return nil, err
	}

	return instance, nil
}

func (ln *LocalNode) sendPeerInfo(w io.Writer) error {
	peerInfo, err := proto.Marshal(&pb.Message{
		Version: env.VERSION,
		Msg: &pb.Message_PeerInfo{
			PeerInfo: &pb.PeerInfo{
				Network: string(ln.discovery.network.ExportPublicKey()),
				Info: ln.Info,
			},
		},
	})
	if err != nil {
		return err
	}
	return pb.Write(w, peerInfo)
}

func (ln *LocalNode) String() string {
	return "Local node."
}