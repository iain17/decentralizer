package network

import ("testing"
	"github.com/stretchr/testify/assert")

func TestNew(t *testing.T) {
	network, err := New()
	assert.NoError(t, err)
	assert.NotEmpty(t, network.Marshal())
	assert.NotNil(t, network.PrivateKey)
	assert.NotNil(t, network.PublicKey)
}

func TestNetwork_InfoHash(t *testing.T) {
	network, err := New()
	assert.NoError(t, err)
	assert.NotEmpty(t, network.InfoHash())
	assert.Equal(t, len(network.InfoHash()), 40)
}

func TestNetwork_Marshal(t *testing.T) {
	network, err := New()
	assert.NoError(t, err)
	networkMarshaled := network.Marshal()

	unmarshalednetwork, err := Unmarshal(networkMarshaled)
	assert.NoError(t, err)

	exportPublicKey1, err := network.ExportPublicKey()
	assert.NoError(t, err)

	exportPublicKey2, err := unmarshalednetwork.ExportPublicKey()
	assert.NoError(t, err)

	assert.Equal(t, exportPublicKey1, exportPublicKey2)
}

func TestNetwork_MarshalFromPrivateKey(t *testing.T) {
	network, err := New()
	assert.NoError(t, err)
	networkMarshaled := network.MarshalFromPrivateKey()

	unmarshalednetwork, err := UnmarshalFromPrivateKey(networkMarshaled)
	assert.NoError(t, err)

	exportPrivateKey1, err := network.ExportPrivateKey()
	assert.NoError(t, err)

	exportPrivateKey2, err := unmarshalednetwork.ExportPrivateKey()
	assert.NoError(t, err)

	exportPublicKey1, err := network.ExportPublicKey()
	assert.NoError(t, err)

	exportPublicKey2, err := unmarshalednetwork.ExportPublicKey()
	assert.NoError(t, err)

	assert.Equal(t, exportPrivateKey1, exportPrivateKey2)
	assert.Equal(t, exportPublicKey1, exportPublicKey2)
}