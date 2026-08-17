package commandline

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
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
			// get current user info
			fileToTransfer := cmd.Args().Get(0)

			userProfile, err := config.GetUser()
			if err != nil {
				return fmt.Errorf("No user found - run `conduit setup <username>` first")
			}

			pubKeyBytes, err := hex.DecodeString(userProfile.PublicKey)
			if err != nil {
				return fmt.Errorf("Error decoding public key: %w", err)
			}

			privKeyBytes, err := hex.DecodeString(userProfile.PrivateKey)
			if err != nil {
				return fmt.Errorf("Error decoding private key: %w", err)
			}

			localPath := filepath.Join(config.LibraryDir(), fileToTransfer)
			fileData, err := os.ReadFile(localPath)
			if err != nil {
				return fmt.Errorf("Could not read file from library: %w", err)
			}

			signature := cryptography.SignData(fileData, privKeyBytes) // sign file data w/ priv key for authenticity/integrity
			header := make([]byte, 1024)
			copy(header[0:512], []byte(fileToTransfer))
			copy(header[512:576], signature)
			copy(header[576:608], pubKeyBytes)

			// get target user info
			targetUser := cmd.Args().Get(1)

			fmt.Printf("Finding user %s...\n", targetUser)
			resp, err := http.Get(config.BaseURL + "/getuser/" + targetUser)
			if err != nil {
				return fmt.Errorf("Failed to get user info: %w", err)
			}
			defer resp.Body.Close()

			// get target user data
			var info UserInfo
			err = json.NewDecoder(resp.Body).Decode(&info)
			if err != nil {
				return fmt.Errorf("Failed to decode user info: %w", err)
			}

			// connect to target user ip
			fmt.Println("Connecting to target user at " + info.IPAddress + "...")
			conn, err := net.Dial("tcp", info.IPAddress+":9000")
			if err != nil {
				return fmt.Errorf("Failed to connect to target user: %w", err)
			}
			defer conn.Close()

			// file transfer
			_, err = conn.Write(header)
			if err != nil {
				return fmt.Errorf("failed to send header: %w", err)
			}

			fmt.Printf("Sending %s (%d bytes)...\n", fileToTransfer, len(fileData))
			_, err = conn.Write(fileData)
			if err != nil {
				return fmt.Errorf("failed to send file data: %w", err)
			}

			fmt.Println("File transfer successfully completed!")
			return nil
		},
	}
}
