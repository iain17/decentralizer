package sessionstore

import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
	"github.com/iain17/decentralizer/pb"
)

func TestStore_FindByPeerId(t *testing.T) {
	store, err := New(10, 10 * time.Second)
	assert.NoError(t, err)
	store.Insert(&pb.SessionInfo {
		DnId: 1,
		PId: "test1",
	})
	store.Insert(&pb.SessionInfo {
		DnId: 2,
		PId: "test2",
	})
	sessions, err := store.FindByPeerId("test1")
	assert.NoError(t, err)
	assert.Equal(t, sessions[0].PId, "test1")
}

func TestStore_FindByDetails(t *testing.T) {
	store, err := New(10, 10 * time.Second)
	assert.NoError(t, err)
	store.Insert(&pb.SessionInfo {
		DnId: 1,
		PId: "test1",
		Details: map[string]string{
			"cool": "yes",
		},
	})
	store.Insert(&pb.SessionInfo {
		DnId: 2,
		PId: "test2",
		Details: map[string]string{
			"cool": "no",
		},
	})
	sessions, err := store.FindByDetails("cool", "no")
	assert.NoError(t, err)
	assert.Equal(t, sessions[0].PId, "test2")
	assert.Equal(t, len(sessions), 1)

}