package discovery

import (
	"context"
	"github.com/grandcat/zeroconf"
	"github.com/iain17/logger"
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

	d.server, err = zeroconf.Register(d.localNode.id, SERVICE, "local.", d.localNode.port, []string{}, nil)
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

	entriesCh := make(chan *zeroconf.ServiceEntry)
	resolver.Browse(d.context, SERVICE, "", entriesCh)
	for {
		select {
		case <-d.context.Done():
			return
		case entry, ok := <-entriesCh:
			if !ok {
				return
			}
			if entry == nil {
				return
			}
			var ip net.IP
			if len(entry.AddrIPv4) > 0 {
				ip = entry.AddrIPv4[0]
			}
			if ip == nil && len(entry.AddrIPv6) > 0  {
				ip = entry.AddrIPv6[0]
			}
			addr := &net.UDPAddr{
				IP:   ip,
				Port: entry.Port,
			}
			d.logger.Debugf("Discovered local peer %s %v", entry.HostName, addr)
			go d.localNode.netTableService.Discovered(addr)
		}
	}
}