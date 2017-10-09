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
	ip string//Gets filled in by stun service.
	port int
	//Services
	listenerService ListenerService
	upNpService UPnPService
	netTableService NetTableService
	StunService StunService
	//Peer discoveries
	discoveryDHT  DiscoveryDHT
	discoveryIRC  DiscoveryIRC
}

func newLocalNode(discovery *Discovery) (*LocalNode, error) {
	instance := &LocalNode{
		Node: Node{
			logger: logging.MustGetLogger("LocalNode"),
			info:   map[string]string{},
		},
		discovery: discovery,
		port:      freeport.GetUDPPort(),
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

	return instance, nil
}

func (ln *LocalNode) sendPeerInfo(w io.Writer) error {
	peerInfo, err := proto.Marshal(&pb.Message{
		Version: env.VERSION,
		Msg: &pb.Message_PeerInfo{
			PeerInfo: &pb.PeerInfo{
				Network: string(ln.discovery.network.ExportPublicKey()),
				Info: ln.info,
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

//Will trigger updating the clients I'm connected to
func (ln *LocalNode) SetInfo(key string, value string) {
	ln.info[key] = value
	for _, peer := range ln.netTableService.GetPeers() {
		go ln.sendPeerInfo(peer.conn)
	}
}