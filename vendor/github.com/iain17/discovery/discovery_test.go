package discovery

import (
	"testing"
	"context"
	"github.com/iain17/discovery/network"
	"time"
	"github.com/stretchr/testify/assert"
)


func TestNetTableService_GetPeers(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	n, err := network.New()
	if err != nil {
		panic(err)
	}
	app1, err := New(ctx, n, 10, false)
	if err != nil {
		panic(err)
	}

	app2, err := New(ctx, n, 10, false)
	if err != nil {
		panic(err)
	}

	app1Peers := app1.WaitForPeers(1, 10 * time.Second)
	app2Peers := app2.WaitForPeers(1, 10 * time.Second)

	assert.Len(t, app1Peers, 1)
	assert.Len(t, app2Peers, 1)
}