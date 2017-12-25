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

// Funnily enough both a private key user and a public key user can encrypt data
// But only the private key user can sign and decrypt
// The private key itself cannot encrypt. It can only encrypt because it also holds the public key
func TestNetwork_Encrypt(t *testing.T) {
	message := []byte("the code must be like a piece of music")
	ciphertext, err := slave.Encrypt(message)//Yes a slave can encrypt data.
	assert.NoError(t, err)
	assert.NotNil(t, ciphertext)

	ciphertext2, err := master.Encrypt(message)//Yes a slave can encrypt data.
	assert.NoError(t, err)
	assert.NotNil(t, ciphertext)
	assert.Equal(t, ciphertext, ciphertext2)
}

// In this test case the message has been written by the creator.
// By simply signing the data we can ensure no other slave is able to change it. because each slave will check if the hash of the message is compatible with the signature.
// Only the creator (master) can create a signature.
func TestNetwork_Sign_Verify(t *testing.T) {
	message := []byte("the code must be like a piece of music. ~creator")
	signature, err := master.Sign(message)
	assert.NoError(t, err)
	assert.NotNil(t, signature)
	//Only a master can sign it.
	_, err = slave.Sign(message)
	assert.Error(t, err)

	//Slave verifies that this message was written by the master
	err = slave.Verify(message, signature)
	assert.NoError(t, err)
	//If the message was changed for whatever reason. It won't check out
	err = slave.Verify([]byte("the code must be like a piece of music. ~evil slave"), signature)
	assert.Error(t, err)
}