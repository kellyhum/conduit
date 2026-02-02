package cryptography

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
)

func GenerateKeyPair() (ed25519.PublicKey, ed25519.PrivateKey, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	return publicKey, privateKey, nil
}

func HashedKey(key ed25519.PublicKey) string {
	hashBytes := sha256.Sum256([]byte(key))
	return string(hashBytes[:])
}
