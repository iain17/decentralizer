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
	var validatorFunc = func(key string, val []byte) error{
		return nil
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	nodes := FakeNewIPFSNodes(ctx, 2)
	app1, err := NewBitSwap(nodes[0])
	assert.NoError(t, err)
	app2, err := NewBitSwap(nodes[1])
	assert.NoError(t, err)

	app1.RegisterValidator(KEYTYPE, validatorFunc, true)
	app2.RegisterValidator(KEYTYPE, validatorFunc, true)

	fmt.Printf("app1 = %s\n", nodes[0].Identity.Pretty())
	fmt.Printf("app2 = %s\n", nodes[1].Identity.Pretty())

	//Execute
	app1.PutValue(KEYTYPE, KEY, []byte{1})

	values, err := app2.GetValues(ctx, KEYTYPE, KEY, 99)
	assert.NoError(t, err)
	assert.Len(t, values, 1)
	assert.Equal(t, values[0].Val[0], byte(1))

	app2.PutValue(KEYTYPE, KEY, []byte{2})

	values, err = app2.GetValues(ctx, KEYTYPE, KEY, 99)
	assert.NoError(t, err)
	assert.True(t, len(values) != 0)
	assert.Equal(t, values[0].Val[0], byte(2))

	values, err = app1.GetValues(ctx, KEYTYPE, KEY, 99)
	assert.NoError(t, err)
	assert.True(t, len(values) != 0)
	assert.Equal(t, values[0].Val[0], byte(2))

	app2.PutValue(KEYTYPE, KEY, []byte{3})
	values, err = app2.GetValues(ctx, KEYTYPE, KEY, 99)
	assert.NoError(t, err)
	assert.True(t, len(values) != 0)
	assert.Equal(t, values[0].Val[0], byte(3))

	values, err = app1.GetValues(ctx, KEYTYPE, KEY, 99)
	assert.NoError(t, err)
	assert.True(t, len(values) != 0)
	assert.Equal(t, values[0].Val[0], byte(3))
}