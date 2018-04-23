package ipfs

import (
	"context"
	"gx/ipfs/QmUvjLCSYy7t4msRzrxfsfj99wboPhTUy7WktCv2LxS7BT/go-ipfs/core"
	"gx/ipfs/QmUvjLCSYy7t4msRzrxfsfj99wboPhTUy7WktCv2LxS7BT/go-ipfs/repo"
	"gx/ipfs/QmUvjLCSYy7t4msRzrxfsfj99wboPhTUy7WktCv2LxS7BT/go-ipfs/repo/config"
	"gx/ipfs/QmUvjLCSYy7t4msRzrxfsfj99wboPhTUy7WktCv2LxS7BT/go-ipfs/repo/fsrepo"
	bitswap "gx/ipfs/QmUvjLCSYy7t4msRzrxfsfj99wboPhTUy7WktCv2LxS7BT/go-ipfs/exchange/bitswap/network"
	dht "gx/ipfs/QmVSep2WwKcXxMonPASsAJ3nZVjfVMKgMcaSigxKnUWpJv/go-libp2p-kad-dht"
	utilmain "github.com/iain17/decentralizer/app/ipfs/util"
	"os"
	"strings"
	//logging "gx/ipfs/QmSpJByNKFX1sCsHBEp3R73FL4NF6FnQTEGyNAXHm2GS52/go-log"
	"github.com/iain17/logger"
)

func init() {
	//logging.SetDebugLogging()
}

func patchSystem() error {
	if err := utilmain.ManageFdLimit(); err != nil {
		logger.Errorf("setting file descriptor limit: %s", err)
	}
	return nil
}

func OpenIPFSRepo(ctx context.Context, path string, limited bool, swarmkey []byte) (*core.IpfsNode, error) {
	r, err := getIPFSRepo(path, limited, swarmkey)
	if err != nil {
		return nil, err
	}
	buildCfg := &core.BuildCfg{
		Repo:      r,
		Online:    true,
		Permanent: !limited,
		ExtraOpts: map[string]bool{
			"mplex":  true,
			"pubsub": false,
		},
	}
	if limited {
		buildCfg.Routing = core.DHTClientOption
	} else {
		buildCfg.Routing = core.DHTOption
	}

	err = patchSystem()
	if err != nil {
		return nil, err
	}

	node, err := core.NewNode(ctx, buildCfg)
	if err != nil {
		return nil, err
	}

	cfg, err := node.Repo.Config()
	if err != nil {
		return nil, err
	}

	//IPFS UPDATER
	//if !cfg.Experimental.FilestoreEnabled {
	//	logger.Info("Enabling experimental file store")
	//	cfg.Experimental.FilestoreEnabled = true
	//	r.SetConfig(cfg)
	//	node.Close()
	//	return OpenIPFSRepo(ctx, path, limited, swarmkey)
	//}

	//Start gateway etc..
	go func() {
		err, gwErrc := serveHTTPGateway(node)
		if err != nil {
			logger.Error(err)
			return
		}
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
		logger.Infof("IPFS API running on %s", cfg.Addresses.API)
		for err := range gwErrc {
			logger.Warning(err)
		}
	}()

	return node, nil
}

func getIPFSRepo(path string, limited bool, swarmkey []byte) (repo.Repo, error) {
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

	err = changeConfig(r, limited)
	return makePrivateRepo(r, swarmkey), err
}

type privateRepo struct{
	repo.Repo
	swarmkey []byte
}

func makePrivateRepo(r repo.Repo, swarmkey []byte) repo.Repo {
	return &privateRepo{Repo: r, swarmkey:swarmkey}
}

func (r *privateRepo) SwarmKey() ([]byte, error) {
	return r.swarmkey, nil
}

func changeConfig(r repo.Repo, limited bool) error {
	rc, err := r.Config()
	if err != nil {
		return err
	}
	rc.Bootstrap = []string{}

	//apiPort := fmt.Sprintf("500%d", 12)
	//gatewayPort := fmt.Sprintf("808%d", 12)
	swarmPort := "4123"

	bitswap.ProtocolBitswap = "/decentralizer/bitswap/testnet/1.1.0"
	dht.ProtocolDHT = "/decentralizer/kad/testnet/1.0.0"

	//rc.Addresses.API = strings.Replace(rc.Addresses.API, "5001", apiPort, -1)
	//rc.Addresses.Gateway = strings.Replace(rc.Addresses.Gateway, "8080", gatewayPort, -1)
	for i, addr := range rc.Addresses.Swarm {
		rc.Addresses.Swarm[i] = strings.Replace(addr, "4001", swarmPort, -1)
	}
	//rc.Swarm.DisableNatPortMap = true
	rc.Swarm.EnableRelayHop = !limited
	err = r.SetConfig(rc)
	return err
}
