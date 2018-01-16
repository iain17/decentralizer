package discovery

import (
	"github.com/iain17/discovery/env"
	"github.com/iain17/discovery/pb"
	"github.com/iain17/logger"
	"github.com/rs/xid"
	"github.com/golang/protobuf/proto"
	"io"
	"github.com/iain17/framed"
	"cirello.io/supervisor"
	"sync"
)

type LocalNode struct {
	Node
	discovery *Discovery
	ip        string //Gets filled in by stun service.
	port      int
	outgoingPort      int
	wg		  *sync.WaitGroup
	lastError error
	supervisor supervisor.Supervisor
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
	i := &LocalNode{
		Node: Node{
			id:        xid.New().String(),
			logger: logger.New("LocalNode"),
			info:   map[string]string{},
		},
		discovery: discovery,
		wg: &sync.WaitGroup{},
	}
	i.supervisor.Log = func(s interface{}) {
		logger.Debugf("[supervisor]: %s", s)
	}
	i.upNpService.localNode = i
	i.supervisor.AddService(&i.upNpService, supervisor.Temporary)
	if !i.discovery.limited {
		i.discoveryDHT.localNode = i
		i.supervisor.AddService(&i.discoveryDHT, supervisor.Permanent)
	}
	i.discoveryIRC.localNode = i
	i.supervisor.AddService(&i.discoveryIRC, supervisor.Permanent)
	i.discoveryMDNS.localNode = i
	i.supervisor.AddService(&i.discoveryMDNS, supervisor.Permanent)

	i.netTableService.localNode = i
	i.supervisor.AddService(&i.netTableService, supervisor.Transient)

	//Listener and stun service
	i.StunService.localNode = i
	i.listenerService.localNode = i
	i.supervisor.AddService(&i.listenerService, supervisor.Permanent)

	numServices := len(i.supervisor.Services())
	i.wg.Add(numServices)
	go i.supervisor.Serve(discovery.ctx)
	return i, i.lastError
}

//Hangs until all servers have at least initialized once
func (ln *LocalNode) waitTilReady() {
	if ln.wg != nil {
		ln.wg.Wait()
		ln.wg = nil
	}
}

func (ln *LocalNode) sendPeerInfo(w io.Writer) error {
	ln.infoMutex.Lock()
	defer ln.infoMutex.Unlock()
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
	ln.infoMutex.Lock()
	defer ln.infoMutex.Unlock()

	if ln.wg == nil {
		ln.info[key] = value
		for _, peer := range ln.netTableService.GetPeers() {
			go ln.sendPeerInfo(peer.conn)
		}
	}
}
