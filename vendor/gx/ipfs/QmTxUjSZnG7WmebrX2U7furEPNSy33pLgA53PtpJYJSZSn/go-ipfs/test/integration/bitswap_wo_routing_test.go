package integrationtest

import (
	"bytes"
	"context"
	"testing"

	"gx/ipfs/QmTxUjSZnG7WmebrX2U7furEPNSy33pLgA53PtpJYJSZSn/go-ipfs/core"
	"gx/ipfs/QmTxUjSZnG7WmebrX2U7furEPNSy33pLgA53PtpJYJSZSn/go-ipfs/core/mock"
	"gx/ipfs/QmYsEQydGrsxNZfAiskvQ76N2xE9hDQtSAkRSynwMiUK3c/go-block-format"

	mocknet "gx/ipfs/Qma23bpHwQrQyvKeBemaeJh7sAoRHggPkgnge1B9489ff5/go-libp2p/p2p/net/mock"
	cid "gx/ipfs/QmeSrf6pzut73u6zLQkRFQ3ygt3k6XFT2kjdYP8Tnkwwyg/go-cid"
)

func TestBitswapWithoutRouting(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	const numPeers = 4

	// create network
	mn := mocknet.New(ctx)

	var nodes []*core.IpfsNode
	for i := 0; i < numPeers; i++ {
		n, err := core.NewNode(ctx, &core.BuildCfg{
			Online:  true,
			Host:    coremock.MockHostOption(mn),
			Routing: core.NilRouterOption, // no routing
		})
		if err != nil {
			t.Fatal(err)
		}
		defer n.Close()
		nodes = append(nodes, n)
	}

	mn.LinkAll()

	// connect them
	for _, n1 := range nodes {
		for _, n2 := range nodes {
			if n1 == n2 {
				continue
			}

			log.Debug("connecting to other hosts")
			p2 := n2.PeerHost.Peerstore().PeerInfo(n2.PeerHost.ID())
			if err := n1.PeerHost.Connect(ctx, p2); err != nil {
				t.Fatal(err)
			}
		}
	}

	// add blocks to each before
	log.Debug("adding block.")
	block0 := blocks.NewBlock([]byte("block0"))
	block1 := blocks.NewBlock([]byte("block1"))

	// put 1 before
	if err := nodes[0].Blockstore.Put(block0); err != nil {
		t.Fatal(err)
	}

	//  get it out.
	for i, n := range nodes {
		// skip first because block not in its exchange. will hang.
		if i == 0 {
			continue
		}

		log.Debugf("%d %s get block.", i, n.Identity)
		b, err := n.Blocks.GetBlock(ctx, cid.NewCidV0(block0.Multihash()))
		if err != nil {
			t.Error(err)
		} else if !bytes.Equal(b.RawData(), block0.RawData()) {
			t.Error("byte comparison fail")
		} else {
			log.Debug("got block: %s", b.Cid())
		}
	}

	// put 1 after
	if err := nodes[1].Blockstore.Put(block1); err != nil {
		t.Fatal(err)
	}

	//  get it out.
	for _, n := range nodes {
		b, err := n.Blocks.GetBlock(ctx, cid.NewCidV0(block1.Multihash()))
		if err != nil {
			t.Error(err)
		} else if !bytes.Equal(b.RawData(), block1.RawData()) {
			t.Error("byte comparison fail")
		} else {
			log.Debug("got block: %s", b.Cid())
		}
	}
}