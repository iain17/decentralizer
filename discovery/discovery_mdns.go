package discovery

import (
	"context"
	"github.com/grandcat/zeroconf"
	"github.com/iain17/logger"
	"os"
	"net"
)

type DiscoveryMDNS struct {
	server      *zeroconf.Server
	localNode *LocalNode
	context context.Context
	logger *logger.Logger
}

func (d *DiscoveryMDNS) Init(ctx context.Context, ln *LocalNode) (err error) {
	d.logger = logger.New("DiscoveryMDNS")
	d.localNode = ln
	d.context = ctx
	infoHash := d.localNode.discovery.network.InfoHash()
	host, _ := os.Hostname()

	d.server, err = zeroconf.Register(host, SERVICE, "local.", d.localNode.port, []string{string(infoHash[:])}, nil)
	if err != nil {
		return err
	}
	go d.Run()
	return err
}

func (d *DiscoveryMDNS) Stop() {
	if d.server != nil {
		d.server.Shutdown()
	}
}

func (d *DiscoveryMDNS) Run() {
	defer func () {
		d.logger.Info("Stopping...")
		d.Stop()
	}()

	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		d.logger.Errorf("Failed to initialize resolver:", err.Error())
	}

	entriesCh := make(chan *zeroconf.ServiceEntry, 4)
	resolver.Browse(d.context, SERVICE, "local.", entriesCh)
	for {
		select {
		case <-d.context.Done():
			return
		case entry, ok := <-entriesCh:
			if !ok {
				return
			}
			ip := entry.AddrIPv4[0]
			if ip == nil {
				ip = entry.AddrIPv6[0]
			}
			go d.localNode.netTableService.Discovered(&net.UDPAddr{
				IP:   ip,
				Port: entry.Port,
			})
			d.logger.Debugf("found entry %v", entry)
		}
	}
}