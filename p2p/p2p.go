package p2p

import (
	core "github.com/ipfs/go-ipfs/core"
	fsrepo "github.com/ipfs/go-ipfs/repo/fsrepo"
	"os"
	"golang.org/x/net/context"
	"github.com/ipfs/go-ipfs/repo/config"
	"github.com/ipfs/go-ipfs/path"
	logger "github.com/Sirupsen/logrus"
	"fmt"
)

type p2p struct {
	Ctx context.Context
	Node *core.IpfsNode
	apps map[string]*p2pApp
}

type P2p interface {
	NewApp(appName string) *p2pApp
	GetApp(appName string) *p2pApp
}

func New(ctx context.Context, repoRoot string) (P2p, error) {
	node, err := ipfsSetup(ctx, repoRoot)
	if err != nil {
		return nil, err
	}
	logger.Infof("I am peer: %s\n", node.Identity.Pretty())
	return &p2p{
		Ctx: ctx,
		Node: node,
		apps: make(map[string]*p2pApp),
	}, nil
}

func ipfsSetup(ctx context.Context, repoRoot string) (*core.IpfsNode, error) {

	if err := checkWriteable(repoRoot); err != nil {
		return nil, err
	}
	if !fsrepo.IsInitialized(repoRoot) {
		conf, err := config.Init(os.Stdout, 2048)
		if err != nil {
			panic(err)
		}
		fsrepo.Init(repoRoot, conf)
	}

	r, err := fsrepo.Open(repoRoot)
	if err != nil {
		return nil, err
	}

	cfg := &core.BuildCfg{
		Repo:   r,
		Online: true,
	}

	return core.NewNode(ctx, cfg)
}

func checkWriteable(dir string) error {
	_, err := os.Stat(dir)
	if err == nil {
		// dir exists, make sure we can write to it
		testfile := path.Join([]string{dir, "test"})
		fi, err := os.Create(testfile)
		if err != nil {
			if os.IsPermission(err) {
				return fmt.Errorf("%s is not writeable by the current user", dir)
			}
			return fmt.Errorf("unexpected error while checking writeablility of repo root: %s", err)
		}
		fi.Close()
		return os.Remove(testfile)
	}

	if os.IsNotExist(err) {
		// dir doesnt exist, check that we can create it
		return os.Mkdir(dir, 0775)
	}

	if os.IsPermission(err) {
		return fmt.Errorf("cannot write to %s, incorrect permissions", err)
	}

	return err
}

func (s *p2p) NewApp(appName string) *p2pApp {
	if s.GetApp(appName) == nil {
		s.apps[appName] = newP2pApp(s, appName)
	}
	return s.apps[appName]
}

func (s *p2p) GetApp(appName string) *p2pApp {
	return s.apps[appName]
}