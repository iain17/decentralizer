package ipfs

import (
	"fmt"
	cmds "gx/ipfs/QmUvjLCSYy7t4msRzrxfsfj99wboPhTUy7WktCv2LxS7BT/go-ipfs/commands"
	manet "gx/ipfs/QmRK2LxanhK2gZq6k6R7vk5ZoYZk8ULSSTB7FzDsMUX6CB/go-multiaddr-net"
	"gx/ipfs/QmUvjLCSYy7t4msRzrxfsfj99wboPhTUy7WktCv2LxS7BT/go-ipfs/core"
	"gx/ipfs/QmUvjLCSYy7t4msRzrxfsfj99wboPhTUy7WktCv2LxS7BT/go-ipfs/core/corehttp"
	ma "gx/ipfs/QmWWQ2Txc2c6tqjsBpzg5Ar652cHPGNsQQp2SejkNmkUMb/go-multiaddr"
	"net"
	"net/http"
)

//Based on: https://sourcegraph.com/gx/ipfs/QmUvjLCSYy7t4msRzrxfsfj99wboPhTUy7WktCv2LxS7BT/go-ipfs@ce22b83f24f72f18318c8649ff1bed3d3e96768e/-/blob/cmd/ipfs/daemon.go#L566
func serveHTTPApi(node *core.IpfsNode) (error, <-chan error) {
	cfg, err := node.Repo.Config()
	if err != nil {
		return fmt.Errorf("serveHTTPApi: GetConfig() failed: %s", err), nil
	}

	apiMaddr, err := ma.NewMultiaddr(cfg.Addresses.API)
	if err != nil {
		return fmt.Errorf("serveHTTPApi: invalid API address: %q (err: %s)", cfg.Addresses.API, err), nil
	}

	apiLis, err := manet.Listen(apiMaddr)
	if err != nil {
		return fmt.Errorf("serveHTTPApi: manet.Listen(%s) failed: %s", apiMaddr, err), nil
	}
	// we might have listened to /tcp/0 - lets see what we are listing on
	apiMaddr = apiLis.Multiaddr()
	fmt.Printf("API server listening on %s\n", apiMaddr)

	gatewayOpt := corehttp.GatewayOption(false, corehttp.WebUIPaths...)
	req, err := cmds.NewEmptyRequest()
	if err != nil {
		return err, nil
	}
	context := req.InvocContext()
	context.ConstructNode = func() (*core.IpfsNode, error) {
		return node, nil
	}

	var opts = []corehttp.ServeOption{
		corehttp.MetricsCollectionOption("api"),
		corehttp.CommandsOption(*context),
		corehttp.WebUIOption,
		gatewayOpt,
		corehttp.VersionOption(),
		defaultMux("/debug/vars"),
		defaultMux("/debug/pprof/"),
		corehttp.MetricsScrapingOption("/debug/metrics/prometheus"),
		corehttp.LogOption(),
	}

	if len(cfg.Gateway.RootRedirect) > 0 {
		opts = append(opts, corehttp.RedirectOption("", cfg.Gateway.RootRedirect))
	}

	if err := node.Repo.SetAPIAddr(apiMaddr); err != nil {
		return fmt.Errorf("serveHTTPApi: SetAPIAddr() failed: %s", err), nil
	}

	errc := make(chan error)
	go func() {
		errc <- corehttp.Serve(node, apiLis.NetListener(), opts...)
		close(errc)
	}()
	return nil, errc
}

func defaultMux(path string) corehttp.ServeOption {
	return func(node *core.IpfsNode, _ net.Listener, mux *http.ServeMux) (*http.ServeMux, error) {
		mux.Handle(path, http.DefaultServeMux)
		return mux, nil
	}
}
