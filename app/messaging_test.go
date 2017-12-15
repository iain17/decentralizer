package app

import (
	"testing"
	"context"
	"github.com/iain17/decentralizer/app/ipfs"
	"github.com/stretchr/testify/assert"
)

func TestDecentralizer_SendMessage(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	nodes := ipfs.FakeNewIPFSNodes(ctx,2)
	app1 := fakeNew(nodes[0])
	assert.NotNil(t, app1)
	app2 := fakeNew(nodes[1])
	assert.NotNil(t, app2)

	go func() {
		msg := <- app2.directMessage
		assert.Equal(t, msg.Message, []byte("hello"))
	}()

	err := app1.SendMessage(app2.i.Identity.Pretty(), []byte("hello"))
	assert.NoError(t, err)
}