package ipfs

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"context"
	"fmt"
)

func TestStringToCid(t *testing.T) {
	expected := "QmbUq44GnfDE5QGVrBKZFBnoHTB9KRXsiaDp2bKKa1WabU"
	actual := StringToCid("cool")
	assert.Equal(t, expected, actual.String())
}

func TestDHTSimple(t *testing.T) {
	//Setup
	const KEYTYPE = "test"
	const KEY = "hey"
	var validatorFunc = func(key string, value []byte) error {
		return nil
	}
	var selectFunc = func(key string, values [][]byte) (int, error) {
		return 0, nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	nodes := FakeNewIPFSNodes(ctx, 2)
	app1, err := NewBitSwap(nodes[0])
	assert.NoError(t, err)
	app2, err := NewBitSwap(nodes[1])
	assert.NoError(t, err)

	app1.RegisterValidator(KEYTYPE, validatorFunc, selectFunc, false)

	app2.RegisterValidator(KEYTYPE, validatorFunc, selectFunc, false)

	fmt.Printf("app1 = %s\n", nodes[0].Identity.Pretty())
	fmt.Printf("app2 = %s\n", nodes[1].Identity.Pretty())

	//Because we can't have randomness in our tests. But live its fine.
	app1.slot = 0
	app2.slot = 1

	//Execute
	app1.PutShardedValues(KEYTYPE, KEY, []byte{1})

	values, err := app2.GetShardedValues(ctx, KEYTYPE, KEY)
	assert.NoError(t, err)
	assert.True(t, len(values) == 1)
	if len(values) != 0 {
		assert.Equal(t, values[0], []byte{1})
	}

	app2.PutShardedValues(KEYTYPE, KEY, []byte{2})

	values, err = app2.GetShardedValues(ctx, KEYTYPE, KEY)
	assert.NoError(t, err)
	assert.True(t, len(values) == 2)
	if len(values) != 0 {
		assert.True(t,   values[0][0] == byte(1) || values[1][0] == byte(1))
		assert.True(t,   values[0][0] == byte(2) || values[1][0] == byte(2))
	}

	values, err = app1.GetShardedValues(ctx, KEYTYPE, KEY)
	assert.NoError(t, err)
	assert.True(t, len(values) == 2)
	if len(values) != 0 {
		assert.True(t,   values[0][0] == byte(1) || values[1][0] == byte(1))
		assert.True(t,   values[0][0] == byte(2) || values[1][0] == byte(2))
	}
}