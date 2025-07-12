package main

import (
	"fmt"

	"github.com/alecthomas/kong"
)

var version = "0.0.1"

type CLI struct {
	Query   string `kong:"arg,optional,help='Datadog query (e.g. service:my-app)'"`
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

	fmt.Println("taildog foundation is working!")
	ctx.Exit(0)
}
