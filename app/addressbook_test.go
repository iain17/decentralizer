package app

import (
	"testing"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/iain17/decentralizer/app/ipfs"
)

func TestDecentralizer_FindByPeerId(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	nodes := ipfs.FakeNewIPFSNodes(ctx,2)
	app1 := fakeNew(nodes[0])
	assert.NotNil(t, app1)
	app2 := fakeNew(nodes[1])
	assert.NotNil(t, app2)

	err := app1.UpsertPeer("self", map[string]string{
		"quote": "these violent delights have violent ends",
	})
	assert.NoError(t, err)

	err = app2.UpsertPeer("self", map[string]string{
		"quote": "these violent delights have violent beginnings",
	})
	assert.NoError(t, err)

	peer, err := app2.FindByPeerId(app1.i.Identity.Pretty())
	assert.NoError(t, err)
	assert.NotNil(t, peer)
	assert.Equal(t, peer.Details["quote"], "these violent delights have violent ends")
}