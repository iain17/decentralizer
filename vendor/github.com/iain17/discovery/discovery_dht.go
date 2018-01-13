package discovery

import (
	"github.com/nictuku/dht"
	"context"
	"time"
	"net"
	"github.com/iain17/logger"
	"strconv"
	"encoding/hex"
)

type DiscoveryDHT struct {
	node      *dht.DHT
	localNode *LocalNode
	context context.Context
	ih dht.InfoHash

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

	cfg := dht.NewConfig()
	cfg.Port = d.localNode.port + PORT_RANGE
	cfg.NumTargetPeers = d.localNode.discovery.max
	cfg.MaxNodes = 100
	d.node, err = dht.New(cfg)

	ih := d.localNode.discovery.network.InfoHash()
	d.ih, err = dht.DecodeInfoHash(hex.EncodeToString(ih[:]))
	if err != nil {
		return
	}
	err = d.node.Start()
	if err != nil {
		return
	}
	return
}

func (d *DiscoveryDHT) Serve(ctx context.Context) {
	defer d.Stop()
	if err := d.init(ctx); err != nil {
		d.localNode.lastError = err
		panic(err)
	}
	d.localNode.waitTilReady()
	if d.node == nil {
		panic("Can't initiate DHT.")
	}
	ticker := time.Tick(HEARTBEAT_DELAY * time.Second)
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
		case r, _ := <-d.node.PeersRequestResults:
			if !d.localNode.netTableService.isEnoughPeers() {
				for _, peers := range r {
					for _, x := range peers {
						host, rawPort, err := net.SplitHostPort(dht.DecodePeerAddress(x))
						if err != nil {
							d.logger.Debug(err)
							continue
						}
						port, err := strconv.Atoi(rawPort)
						if err != nil {
							d.logger.Debug(err)
							continue
						}
						addr := &net.UDPAddr{
							IP:   net.ParseIP(host),
							Port: int(port-PORT_RANGE),
						}
						go d.localNode.netTableService.Discovered(addr)
					}
				}
			}
			break
		}
	}
}

func (d *DiscoveryDHT) Stop() {
	if d.node != nil {
		d.node.Stop()
	}
}

func (d *DiscoveryDHT) request() {
	d.logger.Debugf("sending request '%s'", d.ih.String())
	d.node.PeersRequest(string(d.ih), !d.localNode.netTableService.isEnoughPeers())
}