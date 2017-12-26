package app

import (
	"testing"
	"github.com/iain17/decentralizer/app/ipfs"
	"context"
	"github.com/getlantern/testify/assert"
)

func TestDecentralizer_GetSessionsByDetails(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	nodes := ipfs.FakeNewIPFSNodes(ctx,2)
	app1 := fakeNew(nodes[0], false)
	assert.NotNil(t, app1)
	app2 := fakeNew(nodes[1], false)
	assert.NotNil(t, app2)

	sessId, err := app1.UpsertSession(1337, "App 1 session :D", 303, map[string]string{
		"cool": "yes",
	})
	assert.NoError(t, err)
	assert.True(t, sessId > 0)


	sessId, err = app2.UpsertSession(1337, "App 2 session :D", 304, map[string]string{
		"cool": "no",
	})
	assert.NoError(t, err)
	assert.True(t, sessId > 0)

	//App 2 gets all sessions.
	sessions, err := app2.GetSessions(1337)//All
	assert.NoError(t, err)
	assert.Equal(t, 2, len(sessions))


	//App 1 gets only non cool sessions
	sessions, err = app2.GetSessionsByDetails(1337, "cool", "no")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(sessions))
	assert.Equal(t, sessions[0].Name, "App 2 session :D")
}