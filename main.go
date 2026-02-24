package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/kellyhum/conduit/commandline"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cmd := commandline.Entry()

	err := cmd.Run(context.Background(), os.Args)

	if err != nil {
		logger.Error("Error running program. Aborting...", "error", err)
		os.Exit(1)
	}
}
