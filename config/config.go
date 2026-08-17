package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type UserConfig struct {
	Username             string `json:"username"`
	Id                   string `json:"id"`
	PublicKey            string `json:"public_key"`
	PrivateKey           string `json:"private_key"`
	EncryptionPublicKey  string `json:"encryption_public_key"`
	EncryptionPrivateKey string `json:"encryption_private_key"`
}

const BaseURL = "http://127.0.0.1:8000"

// HELPERS
func HomeDir() string {
	// use the CONDUIT_DIR path if that's set (for testing)
	envHome := os.Getenv("CONDUIT_DIR")
	if envHome != "" {
		return envHome
	}

	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".conduit")
}

func LibraryDir() string {
	// points to ~/.conduit/library/
	return filepath.Join(HomeDir(), "library")
}

func ConfigJSON() string {
	// points to ~/.conduit/config.json
	return filepath.Join(HomeDir(), "config.json")
}

func InitializeLibraryDir() error {
	return os.MkdirAll(LibraryDir(), 0700) // 700 = read/write/execute
}

// USERS
func PrevUser() bool {
	_, err := os.Stat(ConfigJSON()) // t = fileinfo, f = error
	return err == nil
}

func CreateUser(username string, id string, publicKey string, privateKey string, encPublicKey string, encPrivateKey string) error {
	config := UserConfig{
		Username:             username,
		Id:                   id,
		PublicKey:            publicKey,
		PrivateKey:           privateKey,
		EncryptionPublicKey:  encPublicKey,
		EncryptionPrivateKey: encPrivateKey,
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(ConfigJSON(), data, 0600) // 600 = read/write
}

func GetUser() (UserConfig, error) {
	var config UserConfig

	data, err := os.ReadFile(ConfigJSON())
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
