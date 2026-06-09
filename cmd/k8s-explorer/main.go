package main

import (
	"log/slog"
	"os"

	"github.com/maz/k8s-mini-explorer/internal/cli"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, nil)))

	if err := cli.NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
