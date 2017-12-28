package app

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"context"
	"github.com/iain17/decentralizer/app/ipfs"
	"time"
)

//One user saves a file. The other gets it by its hash.
func TestDecentralizer_SaveGetFile(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	nodes := ipfs.FakeNewIPFSNodes(ctx,2)
	app1 := fakeNew(nodes[0], false)
	assert.NotNil(t, app1)
	app2 := fakeNew(nodes[1], false)
	assert.NotNil(t, app2)

	message := []byte("Hey ho this is cool.")

	cid, err := app2.SavePeerFile("test.txt", message)
	assert.NoError(t, err)
	assert.NotNil(t, cid)

	data, err := app1.GetFile(cid)
	assert.NoError(t, err)
	assert.Equal(t, message, data)
}

//One user saves a file. The other gets it by its name and the peer id that saved it.
func TestDecentralizer_SaveGetUserFile(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	nodes := ipfs.FakeNewIPFSNodes(ctx,2)
	app1 := fakeNew(nodes[0], false)
	assert.NotNil(t, app1)
	app2 := fakeNew(nodes[1], false)
	assert.NotNil(t, app2)

	message := []byte("Hey ho this is cool.")
	updatedMessage := []byte("Nah not that cool.")
	filename := "test.txt"

	_, err := app1.SavePeerFile(filename, message)
	assert.NoError(t, err)

	_, err = app1.SavePeerFile(filename, updatedMessage)
	assert.NoError(t, err)

	time.Sleep(2 * time.Second)

	file, err := app2.GetPeerFile(app1.i.Identity.Pretty(), filename)
	assert.NoError(t, err)
	assert.Equal(t, updatedMessage, file)

	_, err = app2.GetPeerFile(app1.i.Identity.Pretty(), "random shit")
	assert.Error(t, err)
}

