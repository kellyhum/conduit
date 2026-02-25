package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type UserConfig struct {
	Username   string `json:"username"`
	Id         string `json:"id"`
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}

const BaseURL = "http://127.0.0.1:8000"

func GetPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".conduit", "config.json") // points to ~/.conduit/config.json
}

func InitializeDir() error {
	home, _ := os.UserHomeDir()
	path := filepath.Join(home, ".conduit")   // ~/.conduit
	library := filepath.Join(path, "library") // ~/.conduit/library

	return os.MkdirAll(library, 0700) // 700 = read/write/execute
}

func PrevUser() bool {
	_, err := os.Stat(GetPath()) // t = fileinfo, f = error
	return err == nil
}

func CreateUser(username string, id string, publicKey string, privateKey string) error {
	config := UserConfig{
		Username:   username,
		Id:         id,
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(GetPath(), data, 0600) // 600 = read/write
}

func GetUser() (UserConfig, error) {
	var config UserConfig

	data, err := os.ReadFile(GetPath())
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
