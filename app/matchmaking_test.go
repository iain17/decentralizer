package app

import (
	"testing"
	"github.com/iain17/decentralizer/app/ipfs"
	"context"
	"github.com/getlantern/testify/assert"
	"time"
	"github.com/iain17/logger"
)

func TestDecentralizer_GetSessionsByDetailsSimple(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	nodes := ipfs.FakeNewIPFSNodes(ctx, 2)
	app1 := fakeNew(ctx, nodes[0], false)
	assert.NotNil(t, app1)
	app2 := fakeNew(ctx, nodes[1], false)
	assert.NotNil(t, app2)

	sessId, err := app2.UpsertSession(1337, "App 2 session :D", 304, map[string]string{
		"cool": "no",
	})
	assert.NoError(t, err)
	assert.True(t, sessId > 0)

	time.Sleep(1 * time.Second)

	//App 1 gets only non cool sessions
	sessions, err := app1.GetSessionsByDetails(1337, "cool", "no")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(sessions))
	assert.Equal(t, sessions[0].Name, "App 2 session :D")
}

func TestDecentralizer_GetSessionsByDetailsSimple2(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	nodes := ipfs.FakeNewIPFSNodes(ctx, 2)
	app1 := fakeNew(ctx, nodes[0], false)
	assert.NotNil(t, app1)
	app2 := fakeNew(ctx, nodes[1], false)
	assert.NotNil(t, app2)
	logger.Infof("app1 = %s", app1.i.Identity.Pretty())
	logger.Infof("app2 = %s", app2.i.Identity.Pretty())

	sessId, err := app2.UpsertSession(1337, "App 2 session :D", 304, map[string]string{
		"cool": "no",
	})
	assert.NoError(t, err)
	assert.True(t, sessId > 0)

	//Now app1 will also publish to DHT.
	sessId, err = app1.UpsertSession(1337, "App 1 session :D", 308, map[string]string{
		"cool": "no",
	})
	assert.NoError(t, err)
	assert.True(t, sessId > 0)

	time.Sleep(1 * time.Second)

	search := app1.getSessionSearch(1337)
	search.refresh()
	store := search.fetch()
	sessions, err := store.FindAll()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(sessions), "App 1 should have both sessions")
}

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

	sessId, err := app2.UpsertSession(1337, "App 2 session :D", 304, map[string]string{
		"cool": "no",
	})
	assert.NoError(t, err)
	assert.True(t, sessId > 0)

	time.Sleep(1 * time.Second)

	//App 1 gets only non cool sessions
	sessions, err := app1.GetSessionsByDetails(1337, "cool", "no")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(sessions))
	assert.Equal(t, sessions[0].Name, "App 2 session :D")
}

//One peer is trying to be evil
func TestDecentralizer_GetSessionsByDetailsEvil(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	nodes := ipfs.FakeNewIPFSNodes(ctx, 3)
	app1 := fakeNew(ctx, nodes[0], false)
	assert.NotNil(t, app1)
	app2 := fakeNew(ctx, nodes[1], false)
	assert.NotNil(t, app2)
	app3 := fakeNew(ctx, nodes[2], false)
	assert.NotNil(t, app2)

	sessId, err := app2.UpsertSession(1337, "App 2 session :D", 304, map[string]string{
		"cool": "no",
	})
	assert.NoError(t, err)
	assert.True(t, sessId > 0)

	go func() {
		for {
			select {
				case <- ctx.Done():
					return
				default:
					app3.b.PutValue(DHT_SESSIONS_KEY_TYPE, app3.getMatchmakingKey(1337), []byte{0, 1, 2})
					time.Sleep(300 * time.Millisecond)
			}
		}
	}()

	time.Sleep(1 * time.Second)

	//App 1 gets only non cool sessions
	sessions, err := app1.GetSessionsByDetails(1337, "cool", "no")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(sessions))
	assert.Equal(t, sessions[0].Name, "App 2 session :D")
}