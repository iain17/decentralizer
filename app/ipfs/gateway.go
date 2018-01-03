package ipfs

import (
	"fmt"
	cmds "gx/ipfs/QmYHpXQEWuhwgRFBnrf4Ua6AZhcqXCYa7Biv65SLGgTgq5/go-ipfs/commands"
	manet "gx/ipfs/QmX3U3YXCQ6UYBxq2LVWF8dARS1hPUTEYLrSx654Qyxyw6/go-multiaddr-net"
	"gx/ipfs/QmYHpXQEWuhwgRFBnrf4Ua6AZhcqXCYa7Biv65SLGgTgq5/go-ipfs/core"
	"gx/ipfs/QmYHpXQEWuhwgRFBnrf4Ua6AZhcqXCYa7Biv65SLGgTgq5/go-ipfs/core/corehttp"
	ma "gx/ipfs/QmXY77cVe7rVRQXZZQRioukUM7aRW3BTcAgJe12MCtb3Ji/go-multiaddr"
)

//Based on: https://sourcegraph.com/gx/ipfs/QmYHpXQEWuhwgRFBnrf4Ua6AZhcqXCYa7Biv65SLGgTgq5/go-ipfs@ce22b83f24f72f18318c8649ff1bed3d3e96768e/-/blob/cmd/ipfs/daemon.go#L566
func serveHTTPGateway(node *core.IpfsNode) (error, <-chan error) {

	cfg, err := node.Repo.Config()
	if err != nil {
		return fmt.Errorf("serveHTTPGateway: GetConfig() failed: %s", err), nil
	}

	gatewayMaddr, err := ma.NewMultiaddr(cfg.Addresses.Gateway)
	if err != nil {
		return fmt.Errorf("serveHTTPGateway: invalid gateway address: %q (err: %s)", cfg.Addresses.Gateway, err), nil
	}

	gwLis, err := manet.Listen(gatewayMaddr)
	if err != nil {
		return fmt.Errorf("serveHTTPGateway: manet.Listen(%s) failed: %s", gatewayMaddr, err), nil
	}
	req, err := cmds.NewEmptyRequest()
	if err != nil {
		return err, nil
	}
	context := req.InvocContext()
	context.ConstructNode = func() (*core.IpfsNode, error) {
		return node, nil
	}

	var opts = []corehttp.ServeOption{
		corehttp.MetricsCollectionOption("gateway"),
		corehttp.CommandsROOption(*context),
		corehttp.VersionOption(),
		corehttp.IPNSHostnameOption(),
		corehttp.GatewayOption(true, "/ipfs", "/ipns"),
	}

	if len(cfg.Gateway.RootRedirect) > 0 {
		opts = append(opts, corehttp.RedirectOption("", cfg.Gateway.RootRedirect))
	}

	errc := make(chan error)
	go func() {
		errc <- corehttp.Serve(node, gwLis.NetListener(), opts...)
		close(errc)
	}()
	return nil, errc
}
