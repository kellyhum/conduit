package config

import (
	"os"
	"path/filepath"
)

func SavePrivateKey(key []byte) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	privPath := filepath.Join(homeDir, "conduit", "private.txt")
	return os.WriteFile(privPath, key, 0600)
}

func LoadPrivateKey() ([]byte, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	privPath := filepath.Join(homeDir, "conduit", "private.txt")
	return os.ReadFile(privPath)
}
