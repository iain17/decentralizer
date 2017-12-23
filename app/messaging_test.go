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
	app1 := fakeNew(nodes[0], false)
	assert.NotNil(t, app1)
	app2 := fakeNew(nodes[1], false)
	assert.NotNil(t, app2)

	ready := make(chan bool)
	go func() {
		msg := <- app2.DirectMessage
		assert.Equal(t, []byte("hello"), msg.Message)
		assert.Equal(t, app1.i.Identity.Pretty(), msg.PId)
		ready <- true
	}()

	err := app1.SendMessage(app2.i.Identity.Pretty(), []byte("hello"))
	assert.NoError(t, err)

	<-ready
}