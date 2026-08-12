package cryptography

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

// SETUP()
func GenerateKeyPair() (ed25519.PublicKey, ed25519.PrivateKey, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	return publicKey, privateKey, nil
}

func HashedKey(key ed25519.PublicKey) string {
	hashBytes := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hashBytes[:])
}

// TRANSFER()
func SignData(data []byte, privateKey []byte) []byte {
	return ed25519.Sign(privateKey, data)
}

// LISTEN()
func VerifySignature(publicKey []byte, data []byte, signature []byte) bool {
	return ed25519.Verify(publicKey, data, signature)
}
