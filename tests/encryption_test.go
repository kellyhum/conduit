package tests

import (
	"bytes"
	"testing"

	"github.com/kellyhum/conduit/cryptography"
)

func TestSharedKeyMatches(t *testing.T) {
	alice, _ := cryptography.GenerateEncryptionKey()
	bob, _ := cryptography.GenerateEncryptionKey()

	alicesKey, _ := cryptography.GetSharedAESKey(alice.Bytes(), bob.PublicKey().Bytes())
	bobsKey, _ := cryptography.GetSharedAESKey(bob.Bytes(), alice.PublicKey().Bytes())

	if !bytes.Equal(alicesKey, bobsKey) {
		t.Errorf("shared key bytes dont match")
	}

	if len(bobsKey) != 32 {
		t.Errorf("bob aes key is not 32 bytes")
	}

	if len(alicesKey) != 32 {
		t.Errorf("alice aes key is not 32 bytes")
	}
}

func TestEncryptDecryptMsg(t *testing.T) {
	alice, _ := cryptography.GenerateEncryptionKey()
	bob, _ := cryptography.GenerateEncryptionKey()
	key, _ := cryptography.GetSharedAESKey(alice.Bytes(), bob.PublicKey().Bytes())

	msg := []byte("hello from alice")

	ciphertext, err := cryptography.EncryptData(msg, key)
	if err != nil {
		t.Fatalf("encryption failed: %v", err)
	}

	if bytes.Contains(ciphertext, msg) {
		t.Errorf("ciphertext still contains the plaintext msg")
	}

	plaintext, err := cryptography.DecryptData(ciphertext, key)
	if err != nil {
		t.Fatalf("decryption failed: %v", err)
	}

	if !bytes.Equal(plaintext, msg) {
		t.Errorf("decrypted msg doesnt match the original")
	}
}

func TestTamperedCiphertextFails(t *testing.T) {
	alice, _ := cryptography.GenerateEncryptionKey()
	bob, _ := cryptography.GenerateEncryptionKey()
	key, _ := cryptography.GetSharedAESKey(alice.Bytes(), bob.PublicKey().Bytes())

	ciphertext, _ := cryptography.EncryptData([]byte("hello from alice"), key)
	ciphertext[len(ciphertext)-1] ^= 0xFF // flip the last byte

	_, err := cryptography.DecryptData(ciphertext, key)
	if err == nil {
		t.Errorf("tampered ciphertext decrypted without an error")
	}
}

func TestDecryptTooShort(t *testing.T) {
	alice, _ := cryptography.GenerateEncryptionKey()
	bob, _ := cryptography.GenerateEncryptionKey()
	key, _ := cryptography.GetSharedAESKey(alice.Bytes(), bob.PublicKey().Bytes())

	_, err := cryptography.DecryptData([]byte{1, 2, 3}, key)
	if err == nil {
		t.Errorf("3 byte input decrypted without an error")
	}
}
