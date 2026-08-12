package commandline

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"

	"github.com/urfave/cli/v3"

	"github.com/kellyhum/conduit/config"
	"github.com/kellyhum/conduit/cryptography"
)

func Setup() *cli.Command {
	return &cli.Command{
		Name:      "setup",
		Usage:     "Set up your conduit client for the first time",
		ArgsUsage: "<username>",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() == 0 {
				return fmt.Errorf("ERROR: no username provided. Please run `conduit setup <username>` to set up your client")
			}

			// CHECK EXISTING
			if config.PrevUser() {
				fmt.Println("Existing configuration found. Attempting to sync with server...")
				user, err := config.GetUser()
				if err != nil {
					return err
				}

				err = registerUser(user.Username, user.Id, user.PublicKey, user.EncryptionPublicKey)
				if err != nil {
					fmt.Printf("ERROR: Failed to sync with server: %v\n", err)
				} else {
					fmt.Println("Successfully synced with server!")
				}

				return nil
			}

			// LOCAL SETUP
			username := cmd.Args().First()
			fmt.Printf("Setting up your conduit client with username: %s...\n", username)

			// dir initializtion
			if err := config.InitializeLibraryDir(); err != nil {
				return fmt.Errorf("failed to initialize directory: %w", err)
			}

			// keygen
			// signing
			pub, priv, err := cryptography.GenerateKeyPair()
			if err != nil {
				return fmt.Errorf("cryptography error: %w", err)
			}

			// data encryption
			encKey, err := cryptography.GenerateEncryptionKey()
			if err != nil {
				return fmt.Errorf("cryptography error: %w", err)
			}

			// keys -> hex b/c saving in json format
			pubString := hex.EncodeToString(pub)
			privString := hex.EncodeToString(priv)
			dataEncPubKey := hex.EncodeToString(encKey.PublicKey().Bytes())
			dataEncPrivKey := hex.EncodeToString(encKey.Bytes())

			// generate random id for user
			id := cryptography.HashedKey(pub)

			// actually make the json
			err = config.CreateUser(username, id, pubString, privString, dataEncPubKey, dataEncPrivKey)
			if err != nil {
				return fmt.Errorf("failed to create user: %w", err)
			}

			fmt.Println("Client setup complete!")

			// SERVER SETUP
			err = registerUser(username, id, pubString, dataEncPubKey)
			if err != nil {
				fmt.Printf("ERROR: Failed to register user with server: %v\n", err)
			} else {
				fmt.Println("Successfully registered user with server!")
			}

			fmt.Println("Setup fully completed - run `./conduit` to view all available commands")
			return nil
		},
	}
}

func registerUser(username string, id string, publicKey string, encryptionPublicKey string) error {
	formData := url.Values{}
	formData.Set("username", username)
	formData.Set("user_id", id)
	formData.Set("public_key", publicKey)
	formData.Set("encryption_public_key", encryptionPublicKey)

	resp, err := http.PostForm(config.BaseURL+"/saveuser", formData)
	if err != nil {
		return fmt.Errorf("Failed to save user to server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Server returned error status %d", resp.StatusCode)
	}

	return nil
}
