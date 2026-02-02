package main

import (
	"log/slog"
	"os"

	"github.com/kellyhum/conduit/cli"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	p := tea.NewProgram(cli.InitialModel())
	_, err := p.Run()

	if err != nil {
		logger.Error("Error running program. Aborting...", "error", err)
		os.Exit(1)
	}
}
