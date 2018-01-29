package peerstore

import (
	"github.com/iain17/decentralizer/pb"
	"github.com/stretchr/testify/assert"
	libp2pPeer "gx/ipfs/QmXYjuNuxVzXKJCfWasQk1RqkhVLDM9jtUKhqc2WPQmFSB/go-libp2p-peer"
	"testing"
	"time"
	"context"
)

func TestStore_FindByPeerId(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	self, err := libp2pPeer.IDB58Decode("QmTq1jNgbHarKgYkZfLJAmtUewyWYniTupQf7ZYsSFQ381")
	assert.NoError(t, err)
	store, err := New(ctx,10, 1*time.Hour, self, "")
	assert.NoError(t, err)
	err = store.Upsert(&pb.Peer{
		DnId: 1,
		PId:  "QmXjGNgJsaktGo7k6h8sBg9zt7fgNC8Rq9JaREVg4QKwuZ",
		Details: map[string]string{
			"hey": "ho",
		},
	})
	assert.NoError(t, err)
	err = store.Upsert(&pb.Peer{
		DnId: 2,
		PId:  "QmQc35EXfiSzJntQUuX8ExpEo6z7PQRriDVnSrtXsNVAdR",
		Details: map[string]string{
			"hey": "ho",
		},
	})
	assert.NoError(t, err)
	peer, err := store.FindByPeerId("QmXjGNgJsaktGo7k6h8sBg9zt7fgNC8Rq9JaREVg4QKwuZ")
	assert.NoError(t, err)
	assert.NotNil(t, peer)
	assert.Equal(t, peer.PId, "QmXjGNgJsaktGo7k6h8sBg9zt7fgNC8Rq9JaREVg4QKwuZ")
}

func TestStore_Expire(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	self, err := libp2pPeer.IDB58Decode("QmTq1jNgbHarKgYkZfLJAmtUewyWYniTupQf7ZYsSFQ381")
	assert.NoError(t, err)
	store, err := New(ctx, 10, 1*time.Second, self, "")
	assert.NoError(t, err)
	err = store.Upsert(&pb.Peer{
		DnId: 1,
		PId:  "QmXjGNgJsaktGo7k6h8sBg9zt7fgNC8Rq9JaREVg4QKwuZ",
		Details: map[string]string{
			"hey": "ho",
		},
	})
	assert.NoError(t, err)

	peers, err := store.FindAll()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(peers))
	time.Sleep(2 * time.Second)
	peers, err = store.FindAll()
	assert.NoError(t, err)
	assert.Equal(t, 0, len(peers))
}

func TestStore_Limit(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	self, err := libp2pPeer.IDB58Decode("QmTq1jNgbHarKgYkZfLJAmtUewyWYniTupQf7ZYsSFQ381")
	assert.NoError(t, err)
	store, err := New(ctx,1, 2*time.Hour, self, "")
	assert.NoError(t, err)
	err = store.Upsert(&pb.Peer{
		DnId: 1,
		PId:  "QmXjGNgJsaktGo7k6h8sBg9zt7fgNC8Rq9JaREVg4QKwuZ",
		Details: map[string]string{
			"hey": "ho",
		},
	})
	assert.NoError(t, err)
	err = store.Upsert(&pb.Peer{
		DnId: 2,
		PId:  "QmQc35EXfiSzJntQUuX8ExpEo6z7PQRriDVnSrtXsNVAdR",
		Details: map[string]string{
			"hey": "ho",
		},
	})
	assert.NoError(t, err)
	time.Sleep(2 * time.Second)
	sessions, err := store.FindAll()
	assert.NoError(t, err)
	assert.True(t, len(sessions) == 1)
}

//Self should never be deleted.
func TestStore_Self(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	self, err := libp2pPeer.IDB58Decode("QmXjGNgJsaktGo7k6h8sBg9zt7fgNC8Rq9JaREVg4QKwuZ")
	assert.NoError(t, err)
	store, err := New(ctx,1, 2*time.Hour, self, "")
	assert.NoError(t, err)
	err = store.Upsert(&pb.Peer{
		DnId: 1,
		PId:  "QmXjGNgJsaktGo7k6h8sBg9zt7fgNC8Rq9JaREVg4QKwuZ",
		Details: map[string]string{
			"username": "mr.ford",
		},
	})
	assert.NoError(t, err)
	err = store.Upsert(&pb.Peer{
		DnId: 2,
		PId:  "QmQc35EXfiSzJntQUuX8ExpEo6z7PQRriDVnSrtXsNVAdR",
		Details: map[string]string{
			"username": "dolores",
		},
	})
	assert.NoError(t, err)
	err = store.Upsert(&pb.Peer{
		DnId: 1,
		PId:  "QmTq1jNgbHarKgYkZfLJAmtUewyWYniTupQf7ZYsSFQ381",
		Details: map[string]string{
			"username": "daddy",
		},
	})
	assert.NoError(t, err)

	time.Sleep(1 * time.Second)
	sessions, err := store.FindAll()
	assert.NoError(t, err)
	assert.True(t, len(sessions) == 2)

	selfPeer, err := store.FindByPeerId("QmXjGNgJsaktGo7k6h8sBg9zt7fgNC8Rq9JaREVg4QKwuZ")
	assert.Equal(t, "mr.ford", selfPeer.Details["username"])
}

func TestStore_FindByDetails(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	self, err := libp2pPeer.IDB58Decode("QmTq1jNgbHarKgYkZfLJAmtUewyWYniTupQf7ZYsSFQ381")
	assert.NoError(t, err)
	store, err := New(ctx,10, 1*time.Hour, self, "")
	assert.NoError(t, err)
	err = store.Upsert(&pb.Peer{
		DnId: 1,
		PId:  "QmXjGNgJsaktGo7k6h8sBg9zt7fgNC8Rq9JaREVg4QKwuZ",
		Details: map[string]string{
			"cool": "yes",
		},
	})
	assert.NoError(t, err)
	err = store.Upsert(&pb.Peer{
		DnId: 2,
		PId:  "QmQc35EXfiSzJntQUuX8ExpEo6z7PQRriDVnSrtXsNVAdR",
		Details: map[string]string{
			"cool": "no",
		},
	})
	assert.NoError(t, err)
	sessions, err := store.FindByDetails("cool", "no")
	assert.NoError(t, err)
	assert.Equal(t, sessions[0].PId, "QmQc35EXfiSzJntQUuX8ExpEo6z7PQRriDVnSrtXsNVAdR")
	assert.Equal(t, len(sessions), 1)
}
