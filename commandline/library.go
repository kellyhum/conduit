package commandline

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"

	"github.com/kellyhum/conduit/config"
)

func Library() *cli.Command {
	return &cli.Command{
		Name:  "library",
		Usage: "Show your user profile and saved/received files",
		Action: func(context.Context, *cli.Command) error {
			userProfile, err := config.GetUser()
			if err != nil {
				fmt.Println("ERROR (LIBRARY): no configuration file found. Please run `./conduit setup <username>` to begin")
				return nil
			}

			fmt.Println("PROFILE ------------")
			fmt.Printf("Username: %s\n", userProfile.Username)
			fmt.Printf("ID: %s\n", userProfile.Id)
			fmt.Printf("Public Key: %s\n", userProfile.PublicKey)
			fmt.Println("--------------------")
			fmt.Println()
			fmt.Println("LIBRARY ------------")
			files, err := GetFiles()
			if err != nil {
				fmt.Println("ERROR: failed to retrieve files from library")
				return nil
			}

			if len(files) > 0 {
				fmt.Println("Files in your library:")
				for _, file := range files {
					fmt.Printf("- %s\n", file)
				}
			} else {
				fmt.Println("Your library is empty. Use `./conduit upload <filename>` to add files to your library!")
			}

			return nil
		},
	}
}

func GetFiles() ([]string, error) {
	dirEntries, err := os.ReadDir(config.LibraryDir())
	if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range dirEntries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}
	return files, nil
}
