//With this file you can create a unique new network or Unmarshal an existing one from a public key
package network

import (
	"crypto/sha1"
	"encoding/hex"
	"crypto/rsa"
	"crypto/rand"
	"errors"
)

type Network struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
	publicKeyBytes []byte
}

func New() (*Network, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	return &Network{
		privateKey: privateKey,
		publicKey:  &privateKey.PublicKey,
	}, nil
}

func (ns Network) Bytes() []byte {
	publicKey, err := ExportPublicKey(ns.publicKey)
	if err != nil {
		panic(err)
	}
	return publicKey
}

func (ns Network) Marshal() string {
	return hex.EncodeToString(ns.Bytes())
}

func (ns Network) MarshalFromPrivateKey() string {
	return hex.EncodeToString(ExportPrivateKey(ns.privateKey))
}

func Unmarshal(v string) (*Network, error) {
	data, err := hex.DecodeString(v)
	if err != nil {
		return nil, err
	}
	publicKey, err := ParsePublicKey(data)
	if err != nil {
		return nil, err
	}
	secret := &Network{
		publicKey: publicKey,
	}
	return secret, nil
}

func UnmarshalFromPrivateKey(v string) (*Network, error) {
	data, err := hex.DecodeString(v)
	if err != nil {
		return nil, err
	}
	privateKey, err := ParsePrivateKey(data)
	if err != nil {
		return nil, err
	}
	publicKey := &privateKey.PublicKey
	secret := &Network{
		privateKey: privateKey,
		publicKey:  publicKey,
	}
	return secret, nil
}

func (ns Network) InfoHash() [20]byte {
	return sha1.Sum(ns.Bytes())
}

func (ns Network) ExportPublicKey() []byte {
	if ns.publicKeyBytes == nil {
		var err error
		ns.publicKeyBytes, err = ExportPublicKey(ns.publicKey)
		if err != nil {
			panic(err)
		}
	}
	return ns.publicKeyBytes
}

func (ns Network) ExportPrivateKey() ([]byte, error) {
	if ns.privateKey == nil {
		return nil, errors.New("No private key set.")
	}
	return ExportPrivateKey(ns.privateKey), nil
}