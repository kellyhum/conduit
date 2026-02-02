package config

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"

	"github.com/kellyhum/conduit/cryptography"
	"github.com/kellyhum/conduit/shared"
)

type UserData struct {
	Username  string
	PublicKey []byte
	Files     []string
}

const ServerURL = "" // replace with your server url

func CreateFirstTimeUser(username string) (string, error) {
	// make keypair
	pub, priv, err := cryptography.GenerateKeyPair()
	if err != nil {
		return "", err
	}

	// generate dir if doens't exist yet
	err = GenerateFolder()
	if err != nil {
		return "", err
	}

	// upload data
	err = SavePrivateKey(priv)
	if err != nil {
		return "", err
	}

	payload := map[string]string{
		"username":   username,
		"public_key": base64.StdEncoding.EncodeToString(pub),
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(ServerURL+"/register", "application/json", bytes.NewReader(payloadBytes))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return "", errors.New(shared.ServerErrorResp + resp.Status)
	}

	return cryptography.HashedKey(pub), err
}

func VerifyReturningUser(username string) (UserData, error) {
	userData, err := LoadUserData()
	if err != nil {
		return UserData{}, err
	}

	// data checks
	if userData.Username != username {
		return UserData{}, errors.New(shared.UsernameNotFound)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return UserData{}, err
	}
	privKeyPath := filepath.Join(homeDir, "conduit", "private.txt")
	_, err = os.Stat(privKeyPath)
	if err != nil {
		return UserData{}, errors.New(shared.PrivateKeyNotFound)
	}

	return userData, nil
}

func LoadUserData() (UserData, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return UserData{}, err
	}

	userFilePath := filepath.Join(homeDir, "conduit", "user_data.json")
	dataBytes, err := os.ReadFile(userFilePath)
	if err != nil {
		return UserData{}, err
	}

	var userData UserData
	err = json.Unmarshal(dataBytes, &userData)
	if err != nil {
		return UserData{}, err
	}

	return userData, nil
}

func GenerateFolder() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	dirPath := filepath.Join(homeDir, "conduit")

	err = os.MkdirAll(dirPath, 0700)
	if err != nil && !os.IsExist(err) { // doesn't exist and throws error
		return err
	}

	return nil
}

func (u UserData) GetUsername() string {
	return u.Username
}

func (u UserData) GetPublicKey() []byte {
	return u.PublicKey
}

func (u UserData) GetFiles() []string {
	return u.Files
}
