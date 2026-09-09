package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Role1776/mcp-x/app/internal/app"
	"github.com/Role1776/mcp-x/app/internal/config"
	"github.com/Role1776/mcp-x/app/internal/pkg/logger"
	"github.com/Role1776/mcp-x/app/internal/pkg/mcpserver"
)

func main() {
	var (
		envPath     string
		showVersion bool
	)

	flag.StringVar(&envPath, "env", "", "path to .env file (optional)")
	flag.BoolVar(&showVersion, "version", false, "print the version and exit")
	flag.Parse()

	// Answerable without credentials on purpose: every other way of asking the
	// binary what it is requires a working X account.
	if showVersion {
		fmt.Fprintln(os.Stdout, mcpserver.Version)

		return
	}

	cfg, err := config.Load(envPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logger := logger.SetupLogger(cfg.Logger)
	if err := app.New(cfg, logger).Run(); err != nil {
		log.Fatalf("failed to run app: %v", err)
	}
}
