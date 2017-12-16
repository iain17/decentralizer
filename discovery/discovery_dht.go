package discovery

import (
	"github.com/nictuku/dht"
	"context"
	"time"
	"net"
	"github.com/iain17/logger"
	"strconv"
)

type DiscoveryDHT struct {
	node      *dht.DHT
	localNode *LocalNode
	context context.Context
	ih dht.InfoHash

	logger *logger.Logger
}

func (d *DiscoveryDHT) Init(ctx context.Context, ln *LocalNode) (err error) {
	d.logger = logger.New("DiscoveryDHT")
	d.localNode = ln
	d.context = ctx

	cfg := dht.NewConfig()
	cfg.Port = d.localNode.port
	d.node, err = dht.New(cfg)

	ih := d.localNode.discovery.network.InfoHash()
	d.ih, err = dht.DecodeInfoHash(string(ih[:]))
	if err != nil {
		return
	}
	err = d.node.Start()
	if err != nil {
		return
	}
	go d.Run()
	return
}

func (d *DiscoveryDHT) Stop() {
	if d.node != nil {
		d.node.Stop()
	}
}

func (d *DiscoveryDHT) Run() {
	defer d.Stop()
	d.request()
	if d.node == nil {
		d.logger.Error("Can't initiate DHT.")
		return
	}

	for {
		select {
		case <-d.context.Done():
			return
		case r, ok := <-d.node.PeersRequestResults:
			if !ok {
				time.Sleep(30 * time.Second)
				d.request()
				continue
			}
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
						go d.localNode.netTableService.Discovered(&net.UDPAddr{
							IP:   net.ParseIP(host),
							Port: int(port),
						})
					}
				}
			}
		}
	}
}

func (d *DiscoveryDHT) request() {
	d.logger.Debugf("sending request '%s'", string(d.ih))
	d.node.PeersRequest(string(d.ih), false)
}