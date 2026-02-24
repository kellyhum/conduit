package commandline

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

func Transfer() *cli.Command {
	return &cli.Command{
		Name:  "transfer",
		Usage: "Transfer files between users",
		Action: func(context.Context, *cli.Command) error {
			fmt.Println("TRANSFER ----------------------------------")

			return nil
		},
	}
}
