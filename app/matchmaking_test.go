package app

import (
	"testing"
	"github.com/iain17/decentralizer/app/ipfs"
	"context"
	"github.com/getlantern/testify/assert"
	"time"
)

func TestDecentralizer_GetSessionsByDetails(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	const num = 10
	nodes := ipfs.FakeNewIPFSNodes(ctx, num)
	var apps []*Decentralizer
	for i := 0; i < num; i++ {
		app := fakeNew(ctx, nodes[i], false)
		assert.NotNil(t, app)
		apps = append(apps, app)
	}
	app1 := apps[5]//Somewhere in near middle
	app2 := apps[8]//At the end

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

	time.Sleep(1 * time.Second)

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