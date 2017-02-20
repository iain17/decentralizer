package p2p

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGetkey(t *testing.T) {
	hash := "iain"
	key := getKey(hash)

	assert.Equal(t, "iain-2017-02-17T22:45:25", key)
}
