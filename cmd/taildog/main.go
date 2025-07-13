package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/alecthomas/kong"
	"github.com/carlosarraes/taildog/internal/client"
)

var version = "0.0.1"

type CLI struct {
	Query   string `kong:"arg,optional,help='Datadog query (e.g. service:my-app)'"`
	APIKey  string `kong:"env='DD_API_KEY',required,help='Datadog API key'"`
	AppKey  string `kong:"env='DD_APPLICATION_KEY',required,help='Datadog Application key'"`
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

	cfg := &client.Config{
		APIKey:     cli.APIKey,
		AppKey:     cli.AppKey,
		Site:       os.Getenv("DD_SITE"),
		Timeout:    30 * time.Second,
		MaxRetries: 3,
	}

	ddClient := client.NewClient(cfg)

	testCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := cli.Query
	if query == "" {
		query = "service:*"
	}

	fmt.Printf("Testing Datadog client with query: %s\n", query)
	response, err := ddClient.FetchLogs(testCtx, query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "API call failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("API call successful! Response type: %T\n", response)
	if response.Data != nil {
		fmt.Printf("Found %d log entries\n", len(response.Data))
		if len(response.Data) > 0 {
			firstLog := response.Data[0]
			fmt.Printf("First log type: %T\n", firstLog)
			if firstLog.Attributes != nil {
				if firstLog.Attributes.Message != nil {
					fmt.Printf("First log has message: %s\n", *firstLog.Attributes.Message)
				}
			}
		}
	}

	fmt.Println("✅ Client creation and API call test successful!")
	ctx.Exit(0)
}
