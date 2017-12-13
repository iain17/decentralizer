package ipfs

import (
	"github.com/ipfs/go-ipfs/core"
	"context"
	mocknet "gx/ipfs/QmXZyBQMkqSYigxhJResC6fLWDGFhbphK67eZoqMDUvBmK/go-libp2p/p2p/net/mock"
	"github.com/ipfs/go-ipfs/core/mock"
	namesys "github.com/ipfs/go-ipfs/namesys"
	pstore "gx/ipfs/QmPgDWmTmuzvP7QE5zwo1TmjbJme9pmZHNujB2453jkCTr/go-libp2p-peerstore"
)

func FakeNewIPFSNodes(ctx context.Context, numPeers int) []*core.IpfsNode {
	// create network
	mn := mocknet.New(ctx)

	var nodes []*core.IpfsNode
	for i := 0; i < numPeers; i++ {
		n, err := core.NewNode(ctx, &core.BuildCfg{
			Online:  true,
			Permament: true,
			Host:    coremock.MockHostOption(mn),
		})
		if err != nil {
			panic(err)
		}
		n.Namesys = namesys.NewNameSystem(n.Routing, n.Repo.Datastore(), 0)
		nodes = append(nodes, n)
	}

	mn.LinkAll()

	// connect them
	for _, n1 := range nodes {
		for _, n2 := range nodes {
			if n1 == n2 {
				continue
			}
			p2 := n2.PeerHost.Peerstore().PeerInfo(n2.PeerHost.ID())
			if err := n1.PeerHost.Connect(ctx, p2); err != nil {
				panic(err)
			}
		}
	}

	bsinf := core.BootstrapConfigWithPeers(
		[]pstore.PeerInfo{
			nodes[0].Peerstore.PeerInfo(nodes[0].Identity),
		},
	)

	for _, n := range nodes[1:] {
		if err := n.Bootstrap(bsinf); err != nil {
			panic(err)
		}
	}

	return nodes
}