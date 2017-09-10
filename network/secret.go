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
	PublicKey *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
}

func New() (*Network, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	return &Network{
		PrivateKey: privateKey,
		PublicKey: &privateKey.PublicKey,
	}, nil
}

func (ns Network) Bytes() []byte {
	publicKey, err := ExportPublicKey(ns.PublicKey)
	if err != nil {
		panic(err)
	}
	return publicKey
}

func (ns Network) Marshal() string {
	return hex.EncodeToString(ns.Bytes())
}

func (ns Network) MarshalFromPrivateKey() string {
	return hex.EncodeToString(ExportPrivateKey(ns.PrivateKey))
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
		PublicKey: publicKey,
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
		PrivateKey: privateKey,
		PublicKey: publicKey,
	}
	return secret, nil
}

func (ns Network) InfoHash() string {
	hashBytes := sha1.Sum(ns.Bytes())
	return hex.EncodeToString(hashBytes[:])
}

func (ns Network) ExportPublicKey() ([]byte, error) {
	return ExportPublicKey(ns.PublicKey)
}

func (ns Network) ExportPrivateKey() ([]byte, error) {
	if ns.PrivateKey == nil {
		return nil, errors.New("No private key set.")
	}
	return ExportPrivateKey(ns.PrivateKey), nil
}