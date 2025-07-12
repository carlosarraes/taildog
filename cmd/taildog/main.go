package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/kong"
)

var version = "0.0.1"

type CLI struct {
	Query   string `kong:"arg,optional,help='Datadog query (e.g. service:my-app)'"`
	APIKey  string `kong:"env='DD_API_KEY',help='Datadog API key'"`
	AppKey  string `kong:"env='DD_APPLICATION_KEY',help='Datadog Application key'"`
	Version bool   `kong:"help='Show version information'"`
}

func validateAuth() error {
	apiKey := os.Getenv("DD_API_KEY")
	if strings.TrimSpace(apiKey) == "" {
		return fmt.Errorf("Missing required environment variable: DD_API_KEY")
	}
	appKey := os.Getenv("DD_APPLICATION_KEY")
	if strings.TrimSpace(appKey) == "" {
		return fmt.Errorf("Missing required environment variable: DD_APPLICATION_KEY")
	}
	return nil
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

	if err := validateAuth(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("CLI parsed successfully")
	ctx.Exit(0)
}
