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

func (d *DiscoveryDHT) Init(ctx context.Context, ln *LocalNode) (err error) {
	d.logger = logger.New("DiscoveryDHT")
	d.localNode = ln
	d.context = ctx

	cfg := dht.NewConfig()
	cfg.Port = d.localNode.port+10
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
	go d.process()
	go d.Run()
	return
}

func (d *DiscoveryDHT) Stop() {
	if d.node != nil {
		d.node.Stop()
	}
}

func (d *DiscoveryDHT) process() {
	for {
		select {
		case <-d.context.Done():
			return
		default:
			time.Sleep(10 * time.Second)//Give us time to find peers on other ways.
			d.request()
		}
	}
}

func (d *DiscoveryDHT) Run() {
	defer d.Stop()
	if d.node == nil {
		d.logger.Error("Can't initiate DHT.")
		return
	}

	for {
		select {
		case <-d.context.Done():
			return
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
							Port: int(port-10),
						}
						d.logger.Debugf("Discovered: %v", addr)
						go d.localNode.netTableService.Discovered(addr)
					}
				}
			}
		}
	}
}

func (d *DiscoveryDHT) request() {
	d.logger.Debugf("sending request '%s'", d.ih.String())
	d.node.PeersRequest(string(d.ih), !d.localNode.netTableService.isEnoughPeers())
}