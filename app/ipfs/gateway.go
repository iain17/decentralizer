package ipfs

import (
	"fmt"
	manet "gx/ipfs/QmV6FjemM1K8oXjrvuq3wuVWWoU2TLDPmNnKrxHzY3v6Ai/go-multiaddr-net"
	ma "gx/ipfs/QmYmsdtJ3HsodkePE3eU3TsCaP2YvPZJ4LoXnNkDE5Tpt7/go-multiaddr"
	cmds "gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/commands"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/core"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/core/corehttp"
)

//Based on: https://sourcegraph.com/github.com/ipfs/go-ipfs@ce22b83f24f72f18318c8649ff1bed3d3e96768e/-/blob/cmd/ipfs/daemon.go#L566
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
