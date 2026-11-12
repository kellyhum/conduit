package commandline

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v3"

	"github.com/kellyhum/conduit/config"
	"github.com/kellyhum/conduit/cryptography"
)

type UserInfo struct {
	PublicKey           string `json:"public_key"`
	EncryptionPublicKey string `json:"encryption_public_key"`
	IPAddress           string `json:"ip_address"`
	Error               string `json:"error"`
}

func Transfer() *cli.Command {
	return &cli.Command{
		Name:      "transfer",
		Usage:     "Transfer files between users",
		ArgsUsage: "<filename> <username>",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() < 2 {
				return fmt.Errorf("usage: conduit transfer <filename> <target_username>")
			}

			// SETUP
			// get current user info
			fileToTransfer := cmd.Args().Get(0)

			userProfile, err := config.GetUser()
			if err != nil {
				return fmt.Errorf("No user found - run `conduit setup <username>` first. Error message: %w", err)
			}

			// ed25519 signing keys
			signPubKey, err := hex.DecodeString(userProfile.PublicKey)
			if err != nil {
				return fmt.Errorf("Error decoding ed25519 signing public key: %w", err)
			}

			signPrivKey, err := hex.DecodeString(userProfile.PrivateKey)
			if err != nil {
				return fmt.Errorf("Error decoding ed25519 signing private key: %w", err)
			}

			// x25519 encryption keys
			encPrivKey, err := hex.DecodeString(userProfile.EncryptionPrivateKey)
			if err != nil {
				return fmt.Errorf("Error decoding x25519 encryption private key: %w", err)
			}

			// file
			localPath := filepath.Join(config.LibraryDir(), fileToTransfer)
			fileData, err := os.ReadFile(localPath)
			if err != nil {
				return fmt.Errorf("Could not read file from library: %w", err)
			}

			// get target user info and data
			targetUser := cmd.Args().Get(1)

			fmt.Printf("Finding user %s...\n", targetUser)
			resp, err := http.Get(config.BaseURL + "/getuser/" + url.PathEscape(targetUser))
			if err != nil {
				return fmt.Errorf("Failed to get user info: %w", err)
			}
			defer resp.Body.Close()

			var targetUserInfo UserInfo
			err = json.NewDecoder(resp.Body).Decode(&targetUserInfo)
			if err != nil {
				return fmt.Errorf("Failed to decode user info: %w", err)
			}

			if targetUserInfo.Error != "" {
				return fmt.Errorf("User %s not found", targetUser)
			}

			targetEncPubKey, err := hex.DecodeString(targetUserInfo.EncryptionPublicKey)
			if err != nil {
				return fmt.Errorf("Error decoding %s's x25519 encryption public key: %w", targetUser, err)
			}

			// ENCRYPT (x25519 + hkdf -> aes-gcm)
			aesKey, err := cryptography.GetSharedAESKey(encPrivKey, targetEncPubKey)
			if err != nil {
				return fmt.Errorf("Failed to get AES key: %w", err)
			}

			ct, err := cryptography.EncryptData(fileData, aesKey)
			if err != nil {
				return fmt.Errorf("Failed to encrypt file data with AES key: %w", err)
			}

			// SIGN (ed25519) + build header
			signature := cryptography.SignData(ct, signPrivKey)
			header := make([]byte, 1024)
			copy(header[0:512], []byte(fileToTransfer))
			copy(header[512:576], signature)
			copy(header[576:608], signPubKey)
			copy(header[608:672], []byte(userProfile.Username)) // b/c u need a username to fetch their pubkey/etc

			// SEND
			// connect to target user ip
			fmt.Println("Connecting to target user at " + targetUserInfo.IPAddress + "...")
			conn, err := net.Dial("tcp", targetUserInfo.IPAddress+":9000")
			if err != nil {
				return fmt.Errorf("Failed to connect to target user: %w", err)
			}
			defer conn.Close()

			// file transfer
			_, err = conn.Write(header)
			if err != nil {
				return fmt.Errorf("failed to send header: %w", err)
			}

			fmt.Printf("Sending %s (%d bytes)...\n", fileToTransfer, len(ct))
			_, err = conn.Write(ct)
			if err != nil {
				return fmt.Errorf("failed to send file data: %w", err)
			}

			fmt.Println("File transfer successfully completed!")
			return nil
		},
	}
}
