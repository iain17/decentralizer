package discovery

import (
	"github.com/iain17/dht"
	"context"
	"time"
	"net"
	"github.com/iain17/logger"
)

type DiscoveryDHT struct {
	node      *dht.Server
	announce *dht.Announce
	localNode *LocalNode
	context context.Context

	logger *logger.Logger
}

func (d *DiscoveryDHT) String() string {
	return "DiscoveryDHT"
}

func (d *DiscoveryDHT) init(ctx context.Context) (err error) {
	defer func() {
		if d.localNode.wg != nil {
			d.localNode.wg.Done()
		}
	}()
	d.logger = logger.New(d.String())
	d.context = ctx
	conn, err := net.ListenPacket("udp", ":0")
	if err != nil {
		return err
	}
	d.node, err = dht.NewServer(&dht.ServerConfig{
		Conn: conn,
		StartingNodes: dht.GlobalBootstrapAddrs,
	})
	if err != nil {
		return err
	}
	return
}

func (d *DiscoveryDHT) Serve(ctx context.Context) {
	defer d.Stop()
	d.localNode.waitTilCoreReady()
	if err := d.init(ctx); err != nil {
		d.localNode.lastError = err
		panic(err)
	}
	d.localNode.WaitTilReady()
	if d.node == nil {
		panic("Can't initiate DHT.")
	}
	ticker := time.Tick(HEARTBEAT_DELAY * time.Second)
	d.request()
	for {
		select {
		case <-d.context.Done():
			return
		case _, ok := <-ticker:
			if !ok {
				break
			}
			d.request()
			break
		case v, ok := <-d.announce.Peers:
			if !ok {
				break
			}
			if !d.localNode.netTableService.isEnoughPeers() {
				for _, peer := range v.Peers {
					addr := &net.UDPAddr{
						IP:   peer.IP[:],
						Port: int(peer.Port),
					}
					d.localNode.netTableService.Discovered(addr)
				}
			}
			break
		}
	}
}

func (d *DiscoveryDHT) Stop() {
	if d.announce != nil {
		d.announce.Close()
	}
	if d.node != nil {
		d.node.Close()
	}
}

func (d *DiscoveryDHT) request() {
	ih := d.localNode.discovery.network.InfoHash()
	d.logger.Debugf("sending request '%x'", ih)

	if d.announce != nil {
		d.announce.Close()
	}
	var err error
	d.announce, err = d.node.Announce(ih, d.localNode.port, false)
	if err != nil {
		logger.Warning(err)
	}
}