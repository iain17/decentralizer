package ipfs

import (
	"context"
	"gx/ipfs/QmNUKMfTHQQpEwE8bUdv5qmKC3ymdW7zw82LFS8D6MQXmu/go-ipfs/core"
	"gx/ipfs/QmNUKMfTHQQpEwE8bUdv5qmKC3ymdW7zw82LFS8D6MQXmu/go-ipfs/core/mock"
	namesys "gx/ipfs/QmNUKMfTHQQpEwE8bUdv5qmKC3ymdW7zw82LFS8D6MQXmu/go-ipfs/namesys"
	pstore "gx/ipfs/QmPgDWmTmuzvP7QE5zwo1TmjbJme9pmZHNujB2453jkCTr/go-libp2p-peerstore"
	"gx/ipfs/QmRQ76P5dgvxTujhfPsCRAG83rC15jgb1G9bKLuomuC6dQ/go-libp2p/p2p/net/mock"
	"time"
)

func FakeNewIPFSNodes(ctx context.Context, numPeers int) []*core.IpfsNode {
	// create network
	mn := mocknet.New(ctx)

	var nodes []*core.IpfsNode
	for i := 0; i < numPeers; i++ {
		n, err := core.NewNode(ctx, &core.BuildCfg{
			Online:    true,
			Permament: true,
			Host:      coremock.MockHostOption(mn),
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
	time.Sleep(1 * time.Second)

	return nodes
}
