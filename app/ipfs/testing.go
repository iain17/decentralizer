package ipfs

import (
	"context"
	"gx/ipfs/QmUvjLCSYy7t4msRzrxfsfj99wboPhTUy7WktCv2LxS7BT/go-ipfs/core"
	"gx/ipfs/QmUvjLCSYy7t4msRzrxfsfj99wboPhTUy7WktCv2LxS7BT/go-ipfs/core/mock"
	namesys "gx/ipfs/QmUvjLCSYy7t4msRzrxfsfj99wboPhTUy7WktCv2LxS7BT/go-ipfs/namesys"
	pstore "gx/ipfs/QmXauCuJzmzapetmC6W4TuDJLL1yFFrVzSHoWv8YdbmnxH/go-libp2p-peerstore"
	"gx/ipfs/QmNh1kGFFdsPu79KNSaL4NUKUPb4Eiz4KHdMtFY6664RDp/go-libp2p/p2p/net/mock"
	"time"
	"github.com/iain17/logger"
)

func FakeNewIPFSNodes(ctx context.Context, numPeers int) []*core.IpfsNode {
	// create network
	mn := mocknet.New(ctx)
	return FakeNewIPFSNodesNetworked(mn, ctx, numPeers, nil)
}

func connectNodes(ctx context.Context, nodes []*core.IpfsNode) {
	for i, n1 := range nodes {
		ii := i + 1
		if ii > len(nodes)-1 {
			continue
		}
		logger.Debugf("Connecting node %d with %d", i, ii)
		n2 := nodes[ii]
		p2 := n2.PeerHost.Peerstore().PeerInfo(n2.PeerHost.ID())
		if err := n1.PeerHost.Connect(ctx, p2); err != nil {
			panic(err)
		}
	}
}

func FakeNewIPFSNodesNetworked(mn mocknet.Mocknet, ctx context.Context, numPeers int, existing []*core.IpfsNode) []*core.IpfsNode {

	var nodes []*core.IpfsNode
	for i := 0; i < numPeers; i++ {
		n, err := core.NewNode(ctx, &core.BuildCfg{
			Online:    true,
			Permanent: false,
			Host:      coremock.MockHostOption(mn),
			ExtraOpts: map[string]bool{
				"mplex":  true,
				"pubsub": true,
			},
		})
		if err != nil {
			panic(err)
		}
		n.Namesys = namesys.NewNameSystem(n.Routing, n.Repo.Datastore(), 0)
		nodes = append(nodes, n)
	}

	mn.LinkAll()

	var bsinf core.BootstrapConfig

	// connect them
	connectNodes(ctx, nodes)
	if existing != nil {
		connectNodes(ctx, existing)
		bsinf = core.BootstrapConfigWithPeers(
			[]pstore.PeerInfo{
				existing[0].Peerstore.PeerInfo(existing[0].Identity),
			},
		)
	} else {
		bsinf = core.BootstrapConfigWithPeers(
			[]pstore.PeerInfo{
				nodes[0].Peerstore.PeerInfo(nodes[0].Identity),
			},
		)
	}

	for _, n := range nodes[1:] {
		if err := n.Bootstrap(bsinf); err != nil {
			panic(err)
		}
	}
	time.Sleep(1 * time.Second)

	return nodes
}
