package commandline

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"

	"github.com/kellyhum/conduit/config"
	"github.com/kellyhum/conduit/cryptography"
	"github.com/urfave/cli/v3"
)

func Listen() *cli.Command {
	return &cli.Command{
		Name:  "listen",
		Usage: "Listen for incoming connections",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			ln, err := net.Listen("tcp", ":9000")
			if err != nil {
				return fmt.Errorf("Failed to start listener: %w", err)
			}
			defer ln.Close()

			for {
				conn, err := ln.Accept()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Accept error: %v\n", err)
					continue
				}
				go handleIncoming(conn)
			}
		},
	}
}

func handleIncoming(conn net.Conn) {
	defer conn.Close()

	header := make([]byte, 1024)
	_, err := io.ReadFull(conn, header)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read header: %v\n", err)
		return
	}

	fileName := string(bytes.Trim(header[0:512], "\x00"))
	signature := header[512:576]
	senderPubKey := header[576:608]

	fileData, err := io.ReadAll(conn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read file data: %v\n", err)
		return
	}

	if !cryptography.VerifySignature(senderPubKey, fileData, signature) {
		fmt.Println("Verification failed")
		return
	}

	if err := config.InitializeLibraryDir(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create the library directory: %v\n", err)
		return
	}

	destPath := filepath.Join(config.LibraryDir(), fileName)

	err = os.WriteFile(destPath, fileData, 0600)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write file: %v\n", err)
		return
	}

	fmt.Println("Received: " + fileName)
}
