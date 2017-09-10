//This file takes care of signing data and verifying data.
package network

//Check: https://github.com/brainattica/Golang-RSA-sample/blob/master/rsa_sample.go

//func (ns *Network) Encrypt(data []byte) ([]byte, error) {
//	hash := sha256.New()
//	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, raulPublicKey, message, label)
//
//}
//
//
//func (ns *Network) Decrypt(data []byte) ([]byte, error) {
//	hash := sha256.New()
//	ciphertext, err := rsa.DecryptOAEP(hash, rand.Reader, raulPublicKey, message, label)
//}