package cryptography

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/hkdf"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
)

func GenerateEncryptionKey() (*ecdh.PrivateKey, error) {
	return ecdh.X25519().GenerateKey(rand.Reader)
}

func EncryptData(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, data, nil), nil
}

func DecryptData(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(data) < gcm.NonceSize() {
		return nil, fmt.Errorf("ciphertext too short, currently at %d bytes but requires min %d bytes", len(data), gcm.NonceSize())
	}

	nonce := data[:gcm.NonceSize()]
	ciphertext := data[gcm.NonceSize():]

	// decrypt + proves data wasn't tampered with
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func GetSharedAESKey(privKey []byte, pubKey []byte) ([]byte, error) {
	// rebuild keys
	curve := ecdh.X25519()
	privateKey, err := curve.NewPrivateKey(privKey)
	if err != nil {
		return nil, fmt.Errorf("private key error: %w", err)
	}
	publicKey, err := curve.NewPublicKey(pubKey)
	if err != nil {
		return nil, fmt.Errorf("public key error: %w", err)
	}

	// compute shared key
	sharedKey, err := privateKey.ECDH(publicKey)
	if err != nil {
		return nil, err
	}

	// key derivation func (hkdf), resehape key to 32 bytes
	aesKey, err := hkdf.Key(sha256.New, sharedKey, nil, "conduit-file-encryption", 32)
	if err != nil {
		return nil, err
	}

	return aesKey, nil
}
