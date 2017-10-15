package network

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var master *Network//He started this lovely network. Has both the private and public key
var slave *Network//Joined this lovely network. So he just has the public key

func init() {
	var err error
	master, err = New()
	if err != nil {
		panic(err)
	}
	slave, err = Unmarshal(master.Marshal())
	if err != nil {
		panic(err)
	}
}

func TestNetwork_Encrypt(t *testing.T) {
	message := []byte("the code must be like a piece of music")
	ciphertext, err := slave.Encrypt(message)//Yes a slave can encrypt data.
	assert.NoError(t, err)
	assert.NotNil(t, ciphertext)
}

func TestNetwork_Sign_Verify(t *testing.T) {
	message := []byte("This message was signed by the creator.")
	ciphertext, err := master.Encrypt(message)
	assert.NoError(t, err)
	assert.NotNil(t, ciphertext)
	//Only a master can sign it.
	_, signature, err := slave.Sign(ciphertext)
	assert.Error(t, err)
	//So let us try that again. This time with the master
	hash, signature, err := master.Sign(ciphertext)
	assert.NoError(t, err)
	//Now a slave can verify if that was sent by a master
	err = slave.Verify(hash, signature)
	assert.NoError(t, err)
}