package ipfs

import (
	"os"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/repo/config"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	"context"
	"fmt"
	"strings"
	"github.com/ipfs/go-ipfs/repo"
	//logging "gx/ipfs/QmSpJByNKFX1sCsHBEp3R73FL4NF6FnQTEGyNAXHm2GS52/go-log"
)

func init() {
	//logging.SetDebugLogging()
}

func OpenIPFSRepo(path string, portIdx int) (*core.IpfsNode, error) {
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

	resetRepoConfigPorts(r, portIdx)

	cfg := &core.BuildCfg{
		Repo:   r,
		Online: true,
		Permament: true,
		ExtraOpts: map[string]bool{
			"pubsub": true,
		},
	}

	node, err := core.NewNode(context.Background(), cfg)
	if err != nil {
		return nil, err
	}
	return node, nil
}

func resetRepoConfigPorts(r repo.Repo, nodeIdx int) {
	if nodeIdx < 0 || nodeIdx > 9 {
		return
	}

	apiPort := fmt.Sprintf("500%d", nodeIdx)
	gatewayPort := fmt.Sprintf("808%d", nodeIdx)
	swarmPort := fmt.Sprintf("400%d", nodeIdx)

	rc, err := r.Config()
	if err != nil {
		panic(err)
	}

	rc.Addresses.API = strings.Replace(rc.Addresses.API, "5001", apiPort, -1)
	rc.Addresses.Gateway = strings.Replace(rc.Addresses.Gateway, "8080", gatewayPort, -1)
	for i, addr := range rc.Addresses.Swarm {
		rc.Addresses.Swarm[i] = strings.Replace(addr, "4001", swarmPort, -1)
	}
	err = r.SetConfig(rc)
	if err != nil {
		panic(err)
	}
}
