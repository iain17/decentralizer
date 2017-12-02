package ipfs

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestStringToCid(t *testing.T) {
	expected := "QmbUq44GnfDE5QGVrBKZFBnoHTB9KRXsiaDp2bKKa1WabU"
	actual := StringToCid("cool")
	assert.Equal(t, expected, actual.String())
}
