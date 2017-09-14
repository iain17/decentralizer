package discovery

import (
	"github.com/nictuku/dht"
	"sync"
	"time"
	"github.com/op/go-logging"
	"context"
)

type DiscoveryDHT struct {
	node      *dht.DHT
	ih        dht.InfoHash
	localNode *LocalNode
	context context.Context

	lastPeers map[string]bool
	mutex     sync.Mutex

	logger *logging.Logger
}

func (d *DiscoveryDHT) Init(ctx context.Context, ln *LocalNode) (err error) {
	d.logger = logging.MustGetLogger("DiscoveryDHT")
	d.localNode = ln
	d.context = ctx
	d.lastPeers = map[string]bool{}

	d.ih, err = dht.DecodeInfoHash(ln.network.InfoHash())
	if err != nil {
		return
	}
	config := dht.NewConfig()
	config.Port = ln.port
	d.node, err = dht.New(config)
	go d.Run()
	return err
}

func (d *DiscoveryDHT) Stop() {
	d.node.Stop()
}

func (d *DiscoveryDHT) Run() {
	ticker := time.NewTicker(10 * time.Second)
	defer func() {
		d.Stop()
		ticker.Stop()
	}()
	d.request()

	for {
		select {
		case <-d.context.Done():
			return
		case <-ticker.C:
			d.request()
		case r := <-d.node.PeersRequestResults:
			for _, peers := range r {
				if len(peers) == 0 {
					d.logger.Debug("No peers received.")
				}
				for _, x := range peers {
					host := dht.DecodePeerAddress(x)
					d.addPeer(host)
				}
			}
		}
	}
}

func (d *DiscoveryDHT) request() {
	d.logger.Debug("sending request...")
	d.node.PeersRequest(string(d.ih), true)
}

func (d *DiscoveryDHT) exists(peer string) bool {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	return d.lastPeers[peer]
}

func (d *DiscoveryDHT) addPeer(peer string) {
	d.logger.Debugf("addPeer %s", peer)
	if d.exists(peer) {
		return
	}
	d.logger.Debugf("new peer %q received", peer)
	d.localNode.netTableService.GetDHTInChannel() <- peer
}