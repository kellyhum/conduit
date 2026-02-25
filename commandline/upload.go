package commandline

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v3"

	"github.com/kellyhum/conduit/config"
)

func Upload() *cli.Command {
	return &cli.Command{
		Name:  "upload",
		Usage: "Upload a file to both library and server",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() == 0 {
				return fmt.Errorf("file path is required")
			}

			filePath := cmd.Args().First()

			userProfile, err := config.GetUser()
			if err != nil {
				fmt.Println("ERROR (UPLOAD): no configuration file found. Please run `./conduit setup <username>` to begin")
				return nil
			}

			fmt.Printf("Uploading file: %s\n", filePath)
			fmt.Printf("User: %s (ID: %s)\n", userProfile.Username, userProfile.Id)

			// upload to local library
			sourceFile, err := os.Open(filePath)
			if err != nil {
				fmt.Printf("ERROR: failed to open file %s\n", filePath)
				return nil
			}
			defer sourceFile.Close()

			home, _ := os.UserHomeDir()
			destinationPath := filepath.Join(home, ".conduit", "library", filepath.Base(filePath))

			libraryDir := filepath.Join(home, ".conduit", "library")
			err = os.MkdirAll(libraryDir, 0755)
			if err != nil {
				return fmt.Errorf("failed to create library directory: %w", err)
			}

			destinationFile, err := os.Create(destinationPath)
			if err != nil {
				fmt.Printf("ERROR: failed to create destination file %s\n", destinationPath)
				return nil
			}
			defer destinationFile.Close()

			copied, err := io.Copy(destinationFile, sourceFile)
			if err != nil {
				fmt.Printf("ERROR: failed to copy file %s to %s\n", filePath, destinationPath)
				return nil
			}
			fmt.Printf("Successfully uploaded %d bytes to library at %s\n", copied, destinationPath)

			// upload to sql
			formData := url.Values{}
			formData.Set("filename", filepath.Base(filePath))
			formData.Set("owner_username", userProfile.Username)

			resp, err := http.PostForm(config.BaseURL+"/upload", formData)
			if err != nil {
				fmt.Printf("ERROR: Saved locally, but failed to upload file metadata to server: %v\n", err)
				return nil
			} else {

				defer resp.Body.Close()
				if resp.StatusCode != http.StatusOK {
					bodyBytes, _ := io.ReadAll(resp.Body)
					fmt.Printf("ERROR: Server responded with status %d: %s\n", resp.StatusCode, string(bodyBytes))
				} else {
					fmt.Printf("Successfully uploaded file metadata to server: %s\n", config.BaseURL)
				}
			}

			return nil
		},
	}
}
