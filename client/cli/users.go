package cli

import (
	"fmt"
	"log/slog"
	"os"
)

func UserSetup() error {
	// if currdirectory doesn't exist
	// return error

	// check if private.key exists + is readable
	// if key doesn't exist in current directory
	// make key
	// if keygeneration error
	// log
	// return error
	// else (it does exist in current directory, so is a prev user)
	// redirect to loaduserdata()

	_, err := os.Getwd()
	if err != nil {
		slog.Error("Failed to get current working directory", "error", err)
		return fmt.Errorf("Failed to get current working directory: %w", err)
	}

	return nil
}
