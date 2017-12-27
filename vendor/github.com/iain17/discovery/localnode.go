package discovery

import (
	"github.com/iain17/discovery/env"
	"github.com/iain17/discovery/pb"
	"github.com/iain17/freeport"
	"github.com/iain17/logger"
	"github.com/rs/xid"
	"github.com/golang/protobuf/proto"
	"io"
	"github.com/iain17/framed"
)

type LocalNode struct {
	Node
	discovery *Discovery
	ip        string //Gets filled in by stun service.
	port      int
	//Services
	listenerService ListenerService
	upNpService     UPnPService
	netTableService NetTableService
	StunService     StunService
	//Peer discoveries
	discoveryDHT  DiscoveryDHT
	discoveryIRC  DiscoveryIRC
	discoveryMDNS DiscoveryMDNS
}

func newLocalNode(discovery *Discovery) (*LocalNode, error) {
	instance := &LocalNode{
		Node: Node{
			id:        xid.New().String(),
			logger: logger.New("LocalNode"),
			info:   map[string]string{},
		},
		discovery: discovery,
		port:      freeport.GetPortRange("udp", 10),
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

	err = instance.StunService.Init(discovery.ctx, instance)
	if err != nil {
		return nil, err
	}

	err = instance.discoveryDHT.Init(discovery.ctx, instance)
	if err != nil {
		return nil, err
	}

	err = instance.discoveryIRC.Init(discovery.ctx, instance)
	if err != nil {
		return nil, err
	}

	err = instance.discoveryMDNS.Init(discovery.ctx, instance)
	if err != nil {
		return nil, err
	}

	return instance, nil
}

func (ln *LocalNode) sendPeerInfo(w io.Writer) error {
	peerInfo, err := proto.Marshal(&pb.Message{
		Version: env.VERSION,
		Msg: &pb.Message_PeerInfo{
			PeerInfo: &pb.DPeerInfo{
				Network: string(ln.discovery.network.ExportPublicKey()),
				Id:      ln.id,
				Info:    ln.info,
			},
		},
	})
	if err != nil {
		return err
	}
	return framed.Write(w, peerInfo)
}

func (ln *LocalNode) String() string {
	return "Local node."
}

//Will trigger updating the clients I'm connected to
func (ln *LocalNode) SetInfo(key string, value string) {
	if ln.info[key] == value {
		return
	}
	ln.info[key] = value
	for _, peer := range ln.netTableService.GetPeers() {
		go ln.sendPeerInfo(peer.conn)
	}
}
