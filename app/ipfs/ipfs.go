package ipfs

import (
	"context"
	"fmt"
	"gx/ipfs/QmNUKMfTHQQpEwE8bUdv5qmKC3ymdW7zw82LFS8D6MQXmu/go-ipfs/core"
	"gx/ipfs/QmNUKMfTHQQpEwE8bUdv5qmKC3ymdW7zw82LFS8D6MQXmu/go-ipfs/repo"
	"gx/ipfs/QmNUKMfTHQQpEwE8bUdv5qmKC3ymdW7zw82LFS8D6MQXmu/go-ipfs/repo/config"
	"gx/ipfs/QmNUKMfTHQQpEwE8bUdv5qmKC3ymdW7zw82LFS8D6MQXmu/go-ipfs/repo/fsrepo"
	"os"
	"strings"
	//logging "gx/ipfs/QmSpJByNKFX1sCsHBEp3R73FL4NF6FnQTEGyNAXHm2GS52/go-log"
	"github.com/iain17/logger"
)

func init() {
	//logging.SetDebugLogging()
}

func OpenIPFSRepo(ctx context.Context, path string, portIdx int) (*core.IpfsNode, error) {
	r, err := getIPFSRepo(path, portIdx)
	if err != nil {
		return nil, err
	}
	cfg := &core.BuildCfg{
		Repo:      r,
		Online:    true,
		Permament: true,
		ExtraOpts: map[string]bool{
		//"pubsub": true,
		},
	}

	node, err := core.NewNode(ctx, cfg)
	if err != nil {
		return nil, err
	}

	//Start gateway etc..
	go func() {
		err, gwErrc := serveHTTPGateway(node)
		if err != nil {
			logger.Error(err)
			return
		}
		cfg, _ := node.Repo.Config()
		logger.Infof("IPFS Gateway running on %s", cfg.Addresses.Gateway)
		for err := range gwErrc {
			logger.Warning(err)
		}
	}()
	go func() {
		err, gwErrc := serveHTTPApi(node)
		if err != nil {
			logger.Error(err)
			return
		}
		cfg, _ := node.Repo.Config()
		logger.Infof("IPFS API running on %s", cfg.Addresses.API)
		for err := range gwErrc {
			logger.Warning(err)
		}
	}()

	return node, nil
}

func getIPFSRepo(path string, portIdx int) (repo.Repo, error) {
	r, err := fsrepo.Open(path)
	if _, ok := err.(fsrepo.NoRepoError); ok {
		var conf *config.Config
		conf, err = config.Init(os.Stdout, 2048)
		if err != nil {
			return nil, err
		}
		err = fsrepo.Init(path, conf)
		if err != nil {
			return nil, err
		}
		r, err = fsrepo.Open(path)
	}
	if err != nil {
		return nil, err
	}

	//err = resetRepoConfigPorts(r, portIdx)
	return r, err
}

func resetRepoConfigPorts(r repo.Repo, nodeIdx int) error {
	rc, err := r.Config()
	if err != nil {
		return err
	}
	if nodeIdx > 0 && nodeIdx < 9 {
		apiPort := fmt.Sprintf("500%d", nodeIdx)
		gatewayPort := fmt.Sprintf("808%d", nodeIdx)
		swarmPort := fmt.Sprintf("400%d", nodeIdx)

		rc.Addresses.API = strings.Replace(rc.Addresses.API, "5001", apiPort, -1)
		rc.Addresses.Gateway = strings.Replace(rc.Addresses.Gateway, "8080", gatewayPort, -1)
		for i, addr := range rc.Addresses.Swarm {
			rc.Addresses.Swarm[i] = strings.Replace(addr, "4001", swarmPort, -1)
		}
	}
	err = r.SetConfig(rc)
	return err
}
