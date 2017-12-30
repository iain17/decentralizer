//This file takes care of signing data and verifying data.
package network

import (
	"crypto/sha256"
	"crypto/rsa"
	"crypto/rand"
	"crypto"
	"errors"
)

//Check: https://github.com/brainattica/Golang-RSA-sample/blob/master/rsa_sample.go
var opts rsa.PSSOptions
func init() {
	opts.SaltLength = rsa.PSSSaltLengthAuto
}

func (ns *Network) Encrypt(data []byte) ([]byte, error) {
	hash := sha256.New()
	return rsa.EncryptOAEP(hash, rand.Reader, ns.publicKey, data, []byte(""))
}

func (ns *Network) Decrypt(data []byte) ([]byte, error) {
	if ns.privateKey == nil {
		return nil, errors.New("no private key set")
	}
	hash := sha256.New()
	return rsa.DecryptOAEP(hash, rand.Reader, ns.privateKey, data, []byte(""))
}

func (ns *Network) Sign(data []byte) (signature []byte, err error) {
	if ns.privateKey == nil {
		return nil, errors.New("no private key set")
	}
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(data)
	hash := pssh.Sum(nil)
	signature, err = rsa.SignPSS(rand.Reader, ns.privateKey, newhash, hash, &opts)
	return
}

func (ns *Network) Verify(data []byte, signature []byte) error {
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(data)
	hash := pssh.Sum(nil)
	return rsa.VerifyPSS(ns.publicKey, crypto.SHA256, hash, signature, &opts)
}