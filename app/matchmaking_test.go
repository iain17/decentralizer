package app

import (
	"testing"
	"github.com/iain17/decentralizer/app/ipfs"
	"context"
	"github.com/getlantern/testify/assert"
	"time"
	"github.com/iain17/logger"
	"github.com/iain17/decentralizer/pb"
	"gx/ipfs/QmNh1kGFFdsPu79KNSaL4NUKUPb4Eiz4KHdMtFY6664RDp/go-libp2p/p2p/net/mock"
	"github.com/iain17/decentralizer/stime"
)

func TestDecentralizer_GetSessionsByDetailsSimple(t *testing.T) {
	EXPIRE_TIME_SESSION = 3 * time.Hour
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	nodes := ipfs.FakeNewIPFSNodes(ctx, 5)
	app1 := fakeNew(ctx, nodes[0], false)
	assert.NotNil(t, app1)
	app2 := fakeNew(ctx, nodes[1], false)
	assert.NotNil(t, app2)

	sessId, err := app2.UpsertSession(1337, "App 2 session :D", 304, map[string]string{
		"cool": "no",
	})
	assert.NoError(t, err)
	assert.True(t, sessId > 0)

	time.Sleep(500 * time.Millisecond)

	//App 1 gets only non cool sessions
	sessions, err := app1.GetSessionsByDetails(1337, "cool", "no")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(sessions))
	if len(sessions) != 0 {
		assert.Equal(t, sessions[0].Name, "App 2 session :D")
	}

	sessId, err = app2.UpsertSession(1337, "App 2 session :D", 304, map[string]string{
		"cool": "yes",
	})
	assert.NoError(t, err)
	assert.True(t, sessId > 0)
}

func TestDecentralizer_GetSessionsByDetailsTrio(t *testing.T) {
	EXPIRE_TIME_SESSION = 3 * time.Hour
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	nodes := ipfs.FakeNewIPFSNodes(ctx, 3)
	app1 := fakeNew(ctx, nodes[0], false)
	assert.NotNil(t, app1)
	app2 := fakeNew(ctx, nodes[1], false)
	assert.NotNil(t, app2)
	app3 := fakeNew(ctx, nodes[2], false)
	assert.NotNil(t, app3)

	_, err := app2.UpsertSession(1337, "App 2 session :D", 304, map[string]string{"0": "0"})
	assert.NoError(t, err)

	app1Search := app1.getSessionSearch(1337)
	app1Store, err := app1Search.fetch()
	assert.NoError(t, err)
	app2Search := app2.getSessionSearch(1337)
	app2Store, err := app2Search.fetch()
	assert.NoError(t, err)
	app3Search := app3.getSessionSearch(1337)
	app3Store, err := app3Search.fetch()
	assert.NoError(t, err)
	assert.Equal(t, 1, app1Store.Len())
	assert.Equal(t, 1, app2Store.Len())
	assert.Equal(t, 1, app3Store.Len())

	_, err = app3.UpsertSession(1337, "App 3 session :D", 308, map[string]string{"0": "0"})
	assert.NoError(t, err)

	app1Store, err = app1Search.fetch()
	assert.NoError(t, err)
	app2Store, err = app2Search.fetch()
	assert.NoError(t, err)
	app3Store, err = app3Search.fetch()
	assert.NoError(t, err)
	assert.Equal(t, 2, app1Store.Len())
	assert.Equal(t, 2, app2Store.Len())
	assert.Equal(t, 2, app3Store.Len())

	_, err = app1.UpsertSession(1337, "App 1 session :D", 111, map[string]string{"0": "0"})
	assert.NoError(t, err)

	app1Store, err = app1Search.fetch()
	assert.NoError(t, err)
	app2Store, err = app2Search.fetch()
	assert.NoError(t, err)
	app3Store, err = app3Search.fetch()
	assert.NoError(t, err)
	assert.Equal(t, 3, app1Store.Len())
	assert.Equal(t, 3, app2Store.Len())
	assert.Equal(t, 3, app3Store.Len())
}

//2 peer have published their session. Then all a third peer joins the network. He should have two sessions.
func TestDecentralizer_GetSessionsLateJoiner(t *testing.T) {
	EXPIRE_TIME_SESSION = 3 * time.Hour
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	mn := mocknet.New(ctx)
	nodes := ipfs.FakeNewIPFSNodesNetworked(mn, ctx, 5, nil)
	app1 := fakeNew(ctx, nodes[0], false)
	assert.NotNil(t, app1)
	app2 := fakeNew(ctx, nodes[1], false)
	assert.NotNil(t, app2)
	app3 := fakeNew(ctx, nodes[2], false)
	assert.NotNil(t, app3)

	_, err := app1.UpsertSession(1337, "App 1 session :D", 304, map[string]string{
		"cool": "no",
	})
	assert.NoError(t, err)

	_, err = app2.UpsertSession(1337, "App 2 session :D", 305, map[string]string{
		"cool": "maybe",
	})
	assert.NoError(t, err)

	_, err = app3.UpsertSession(1337, "App 3 session :D", 306, map[string]string{
		"cool": "yes",
	})
	assert.NoError(t, err)

	time.Sleep(500 * time.Millisecond)

	app1Search := app1.getSessionSearch(1337)
	app1Search.refresh(ctx)

	time.Sleep(1 * time.Second)

	assert.Equal(t, 3, app1Search.storage.Len())

	//Now the late joiner
	lateNodes := ipfs.FakeNewIPFSNodesNetworked(mn, ctx, 2, nodes)

	app4 := fakeNew(ctx, lateNodes[1], false)
	assert.NotNil(t, app2)

	app4Search := app4.getSessionSearch(1337)

	time.Sleep(500 * time.Millisecond)

	app4Search.refresh(ctx)

	time.Sleep(1 * time.Second)

	assert.Equal(t, 3, app4Search.storage.Len())
}

func TestDecentralizer_GetSessionsByDetailsSimple2(t *testing.T) {
	EXPIRE_TIME_SESSION = 3 * time.Hour
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

	time.Sleep(500 * time.Millisecond)

	search := app1.getSessionSearch(1337)
	search.refresh(ctx)
	store, err := search.fetch()
	assert.NoError(t, err)
	sessions, err := store.FindAll()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(sessions), "App 1 should have both sessions")
}

func TestDecentralizer_GetSessionsByDetails(t *testing.T) {
	EXPIRE_TIME_SESSION = 3 * time.Hour
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

	time.Sleep(500 * time.Millisecond)

	//App 1 gets only non cool sessions
	sessions, err := app1.GetSessionsByDetails(1337, "cool", "no")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(sessions))
	assert.Equal(t, sessions[0].Name, "App 2 session :D")
}

//One peer is trying to be evil
func TestDecentralizer_GetSessionsByDetailsEvil(t *testing.T) {
	EXPIRE_TIME_SESSION = 3 * time.Hour
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

	time.Sleep(500 * time.Millisecond)

	//App 1 gets only non cool sessions
	sessions, err := app1.GetSessionsByDetails(1337, "cool", "no")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(sessions))
	assert.Equal(t, sessions[0].Name, "App 2 session :D")
}

func TestValidateDNSessions(t *testing.T) {
	EXPIRE_TIME_SESSION = 3 * time.Hour

	//Future
	assert.Error(t, validateDNSessionsRecord(&pb.DNSessionsRecord{
		Published: uint64(stime.Now().Add(2 * time.Hour).Unix()),
	}))

	//Expired
	assert.Error(t, validateDNSessionsRecord(&pb.DNSessionsRecord{
		Published: uint64(stime.Now().AddDate(-3, 0, 0).Unix()),
	}))

	//Alright
	assert.NoError(t, validateDNSessionsRecord(&pb.DNSessionsRecord{
		Published: uint64(stime.Now().Unix()),
	}))

	EXPIRE_TIME_SESSION = 1 * time.Second
	//Just in time
	assert.NoError(t, validateDNSessionsRecord(&pb.DNSessionsRecord{
		Published: uint64(stime.Now().Unix()),
	}))
}

func TestDecentralizer_GetSessionsByDetailsExpire(t *testing.T) {
	EXPIRE_TIME_SESSION = 2 * time.Second
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

	time.Sleep(500 * time.Millisecond)

	search := app1.getSessionSearch(1337)
	search.refresh(ctx)
	store, err := search.fetch()
	assert.NoError(t, err)
	assert.NoError(t, err)
	assert.Equal(t, 1, store.Len(), "Because it hasn't YET expired")

	time.Sleep(EXPIRE_TIME_SESSION * 2)

	searchCtx, cancel := context.WithCancel(app1.i.Context())
	search.refresh(searchCtx)
	assert.NoError(t, err)
	assert.Equal(t, 0, store.Len(), "Because it expired")
	cancel()

	_, err = app2.UpsertSession(1337, "App 2 session :D", 304, map[string]string{
		"cool": "no",
	})
	assert.NoError(t, err)
	time.Sleep(500 * time.Millisecond)

	searchCtx, cancel = context.WithCancel(app1.i.Context())
	search.refresh(searchCtx)
	time.Sleep(500 * time.Millisecond)
	assert.NoError(t, err)
	assert.Equal(t, 1, store.Len(), "Because app2 has republished again")
	cancel()
}