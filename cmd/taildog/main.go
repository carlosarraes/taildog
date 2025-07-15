package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/alecthomas/kong"
	"github.com/carlosarraes/taildog/internal/client"
	"github.com/carlosarraes/taildog/internal/config"
	"github.com/carlosarraes/taildog/internal/signals"
	"github.com/carlosarraes/taildog/internal/tailing"
)

var version = "0.0.1"

type CLI struct {
	Query   string `kong:"arg,optional,help='Datadog query (e.g. service:my-app)'"`
	APIKey  string `kong:"env='DD_API_KEY',required,help='Datadog API key'"`
	AppKey  string `kong:"env='DD_APPLICATION_KEY',required,help='Datadog Application key'"`
	Site    string `kong:"env='DD_SITE',help='Datadog site (default: datadoghq.com)'"`
	Version bool   `kong:"help='Show version information'"`
}

func main() {
	var cli CLI
	ctx := kong.Parse(&cli,
		kong.Name("taildog"),
		kong.Description("A modern CLI tool to tail Datadog logs in real-time"),
		kong.UsageOnError(),
	)

	if cli.Version {
		fmt.Printf("taildog version %s\n", version)
		return
	}

	cfg, err := config.NewConfig(cli.APIKey, cli.AppKey, cli.Site, cli.Query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Configuration error: %v\n", err)
		os.Exit(1)
	}

	ddClient := client.NewClient(cfg.ToClientConfig())

	authCtx, authCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer authCancel()

	fmt.Println("Testing authentication...")
	if err := config.TestAuthentication(authCtx, ddClient); err != nil {
		fmt.Fprintf(os.Stderr, "Authentication failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ Authentication successful!")

	query := cfg.Query
	if query == "" {
		query = "service:*"
	}

	fmt.Printf("Starting to tail logs with query: %s\n", query)
	fmt.Println("Press Ctrl+C to stop...")
	fmt.Println()

	tailCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalHandler := signals.NewHandler(cancel)
	signalHandler.SetupGracefulShutdown()

	tailer := tailing.NewTailer(ddClient, 5*time.Second)
	if err := tailer.Start(tailCtx, query); err != nil {
		fmt.Fprintf(os.Stderr, "Tailing error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n✅ Log tailing stopped")
	ctx.Exit(0)
}
