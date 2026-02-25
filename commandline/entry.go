package commandline

import (
	"github.com/urfave/cli/v3"
)

func Entry() *cli.Command {
	return &cli.Command{
		Name:  "conduit",
		Usage: "CLI File Transfer Tool",
		Commands: []*cli.Command{
			Setup(),
			Library(),
			Transfer(),
			Listen(),
			Upload(),
		},
	}
}
