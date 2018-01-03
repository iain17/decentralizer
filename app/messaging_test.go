package app

import (
	"testing"
	"context"
	"github.com/iain17/decentralizer/app/ipfs"
	"github.com/stretchr/testify/assert"
	"sync"
)

func TestDecentralizer_SendMessage(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	nodes := ipfs.FakeNewIPFSNodes(ctx,2)
	app1 := fakeNew(ctx, nodes[0], false)
	assert.NotNil(t, app1)
	app2 := fakeNew(ctx, nodes[1], false)
	assert.NotNil(t, app2)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		msg := <- app2.GetMessagingChan(0)
		assert.Equal(t, []byte("hello"), msg.Message)
		assert.Equal(t, app1.i.Identity.Pretty(), msg.PId)
		wg.Done()
	}()

	err := app1.SendMessage(0, app2.i.Identity.Pretty(), []byte("hello"))
	assert.NoError(t, err)

	wg.Wait()
}