package ipfs

import (
	"context"
	utilmain "github.com/iain17/decentralizer/app/ipfs/util"
	//dht "gx/ipfs/QmTktQYCKzQjhxF6dk5xJPRuhHn3JBiKGvMLoiDy1mYmxC/go-libp2p-kad-dht"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/core"
	bitswap "gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/exchange/bitswap/network"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/repo"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/repo/config"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/repo/fsrepo"
	"os"
	"strings"
	//logging "gx/ipfs/QmcVVHfdyv15GVPk7NrxdWjh2hLVccXnoD8j2tyQShiXJb/go-log"
	"github.com/iain17/logger"
	//"gx/ipfs/QmZNkThpqfVXs9GNbexPrfBbXSLNYeKrE7jwFM2oqHbyqN/go-libp2p-protocol"
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

	//cfg, err := node.Repo.Config()
	//if err != nil {
	//	return nil, err
	//}

	//IPFS UPDATER
	//if !cfg.Experimental.FilestoreEnabled {
	//	logger.Info("Enabling experimental file store")
	//	cfg.Experimental.FilestoreEnabled = true
	//	r.SetConfig(cfg)
	//	node.Close()
	//	return OpenIPFSRepo(ctx, path, limited, swarmkey)
	//}

	//Start gateway etc..
	//go func() {
	//	err, gwErrc := serveHTTPGateway(node)
	//	if err != nil {
	//		logger.Error(err)
	//		return
	//	}
	//	logger.Infof("IPFS Gateway running on %s", cfg.Addresses.Gateway)
	//	for err := range gwErrc {
	//		logger.Warning(err)
	//	}
	//}()
	//go func() {
	//	err, gwErrc := serveHTTPApi(node)
	//	if err != nil {
	//		logger.Error(err)
	//		return
	//	}
	//	logger.Infof("IPFS API running on %s", cfg.Addresses.API)
	//	for err := range gwErrc {
	//		logger.Warning(err)
	//	}
	//}()

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

type privateRepo struct {
	repo.Repo
	swarmkey []byte
}

func makePrivateRepo(r repo.Repo, swarmkey []byte) repo.Repo {
	return &privateRepo{Repo: r, swarmkey: swarmkey}
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

	bitswap.ProtocolBitswapOne = "/decentralizer/bitswap/1.0.0"
	bitswap.ProtocolBitswapNoVers = "/decentralizer/bitswap"
	bitswap.ProtocolBitswap = "/decentralizer/bitswap/1.1.0"

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
