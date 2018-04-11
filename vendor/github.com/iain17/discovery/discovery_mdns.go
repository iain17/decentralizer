package discovery

import (
	"context"
	"github.com/grandcat/zeroconf"
	"github.com/iain17/logger"
	"net"
	"fmt"
)

type DiscoveryMDNS struct {
	server      *zeroconf.Server
	resolver *zeroconf.Resolver
	localNode *LocalNode
	context context.Context
	logger *logger.Logger
}

func (d *DiscoveryMDNS) String() string {
	return "DiscoveryMDNS"
}

func (d *DiscoveryMDNS) init(ctx context.Context) (err error) {
	defer func() {
		if d.localNode.wg != nil {
			d.localNode.wg.Done()
		}
	}()
	d.logger = logger.New(d.String())
	d.context = ctx
	if err != nil {
		return err
	}
	d.resolver, err = zeroconf.NewResolver(nil)
	if err != nil {
		return fmt.Errorf("failed to initialize resolver: %s", err.Error())
	}
	d.server, err = zeroconf.Register(d.localNode.id, SERVICE, "local.", d.localNode.port, []string{d.localNode.id}, nil)
	if err != nil {
		return err
	}
	return err
}

func (d *DiscoveryMDNS) Serve(ctx context.Context) {
	defer d.Stop()
	d.localNode.waitTilCoreReady()

	if err := d.init(ctx); err != nil {
		d.localNode.lastError = err
		panic(err)
	}
	d.localNode.WaitTilReady()
	entriesCh := make(chan *zeroconf.ServiceEntry)
	err := d.resolver.Browse(d.context, SERVICE, "", entriesCh)
	if err != nil {
		panic(err)
	}
	for {
		select {
		case <-d.context.Done():
			return
		case entry, ok := <-entriesCh:
			if !ok {
				continue
			}
			if entry == nil {
				continue
			}
			if len(entry.Text) != 1 || entry.Text[0] == d.localNode.id {
				continue
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
			d.localNode.netTableService.Discovered(addr)
		}
	}
}

func (d *DiscoveryMDNS) Stop() {
	if d.server != nil {
		d.server.Shutdown()
	}
}