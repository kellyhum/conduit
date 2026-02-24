package commandline

import (
	"context"
	"encoding/hex"
	"fmt"

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

			// SETUP
			username := cmd.Args().First()
			fmt.Printf("Setting up your conduit client with username: %s...\n", username)

			// dir initializtion
			if err := config.InitializeDir(); err != nil {
				return fmt.Errorf("failed to initialize directory: %w", err)
			}

			// keygen
			pub, priv, err := cryptography.GenerateKeyPair()
			if err != nil {
				return fmt.Errorf("cryptography error: %w", err)
			}

			// keys -> hex b/c saving in json format
			pubString := hex.EncodeToString(pub)
			privString := hex.EncodeToString(priv)

			// generate random id for user
			id := cryptography.HashedKey(pub)

			// actually make the json
			err = config.CreateUser(username, id, pubString, privString)
			if err != nil {
				return fmt.Errorf("failed to create user: %w", err)
			}

			fmt.Println("Client setup complete! Run `./conduit --help` to view all available commands")

			return nil
		},
	}
}
