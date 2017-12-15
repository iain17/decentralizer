package sessionstore

import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
	"github.com/iain17/decentralizer/pb"
)

func TestStore_FindByPeerId(t *testing.T) {
	store, err := New(10, 1 * time.Hour)
	assert.NoError(t, err)
	_, err = store.Insert(&pb.Session {
		DnId: 1,
		PId: "test1",
		Details: map[string]string{
			"hey": "ho",
		},
	})
	assert.NoError(t, err)
	_, err = store.Insert(&pb.Session {
		DnId: 2,
		PId: "test2",
		Details: map[string]string{
			"hey": "ho",
		},
	})
	assert.NoError(t, err)
	sessions, err := store.FindByPeerId("test1")
	assert.NoError(t, err)
	assert.True(t, len(sessions) == 1)
	assert.Equal(t, sessions[0].PId, "test1")
}

func TestStore_Expire(t *testing.T) {
	store, err := New(10, 1 * time.Second)
	assert.NoError(t, err)
	_, err = store.Insert(&pb.Session {
		DnId: 1,
		PId: "test1",
		Details: map[string]string{
			"hey": "ho",
		},
	})
	assert.NoError(t, err)

	sessions, err := store.FindAll()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(sessions))
	time.Sleep(2 * time.Second)
	sessions, err = store.FindAll()
	assert.NoError(t, err)
	assert.Equal(t, 0, len(sessions))
}

func TestStore_Limit(t *testing.T) {
	store, err := New(1, 2 * time.Hour)
	assert.NoError(t, err)
	_, err = store.Insert(&pb.Session {
		DnId: 1,
		PId: "test1",
		Details: map[string]string{
			"hey": "ho",
		},
	})
	assert.NoError(t, err)
	_, err = store.Insert(&pb.Session {
		DnId: 2,
		PId: "test2",
		Details: map[string]string{
			"hey": "ho",
		},
	})
	assert.NoError(t, err)
	time.Sleep(1 * time.Second)
	sessions, err := store.FindAll()
	assert.NoError(t, err)
	assert.True(t, len(sessions) == 1)
}

func TestStore_FindByDetails(t *testing.T) {
	store, err := New(10, 1 * time.Hour)
	assert.NoError(t, err)
	_, err = store.Insert(&pb.Session {
		DnId: 1,
		PId: "test1",
		Details: map[string]string{
			"cool": "yes",
		},
	})
	assert.NoError(t, err)
	_, err = store.Insert(&pb.Session {
		DnId: 2,
		PId: "test2",
		Details: map[string]string{
			"cool": "no",
		},
	})
	assert.NoError(t, err)
	sessions, err := store.FindByDetails("cool", "no")
	assert.NoError(t, err)
	assert.Equal(t, sessions[0].PId, "test2")
	assert.Equal(t, len(sessions), 1)
}