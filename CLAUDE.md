# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

taildog is a Go CLI application for real-time log tailing from Datadog's Logs API v2. It's built with Kong for CLI parsing and follows standard Go project structure.

## Development Commands

### Building and Testing

- `make build` - Build binary and install to ~/.local/bin
- `make build-dev` - Quick build without optimizations for development
- `make test` - Run all tests
- `make test-cover` - Run tests with coverage report
- `make fmt` - Format all Go code
- `make vet` - Run go vet
- `make check` - Run fmt, vet, and test together

### Development Workflow

- `make dev` - Quick development cycle (build + show help)
- `make run` - Build and run with example query
- `make demo` - Run demo commands showing version and help

### Release Management

- `make version` - Show current version from cmd/taildog/main.go
- `make tag` - Create git tag for current version
- `make release` - Full release preparation (clean, check, build)
- `make build-all` - Cross-compile for all platforms (Linux/macOS, x86_64/arm64)

## Architecture

### Project Structure

```
cmd/taildog/main.go          # CLI entry point with Kong
internal/client/             # Datadog API client (empty, to be implemented)
internal/config/             # Configuration management (empty, to be implemented)
internal/output/             # Output formatting and writing (empty, to be implemented)
pkg/types/types.go           # Shared types and interfaces (currently empty)
```

### Key Design Decisions

- Uses Kong for CLI argument parsing
- Supports optional query argument (for getting all logs vs filtered logs)
- Standard Go project layout with cmd/ and internal/ separation
- Modular design planned for client, config, and output handling

## Technology Stack

- **Go 1.24+** with standard library
- **Kong** for CLI parsing (already integrated)
- **Datadog Logs API v2** for log retrieval
- **Standard net/http** for API calls
- Future: fatih/color for terminal output, TOML for config

## Code Guidelines

- Maximum 400 lines per source file (excluding tests)
- Modular design with clear separation of concerns
- Follow Go best practices and conventions
- TDD approach with comprehensive test coverage planned

## Authentication

Uses environment variables:

- `DD_API_KEY` (required)
- `DD_APPLICATION_KEY` (required)
- `DD_SITE` (optional, defaults to datadoghq.com)

## Version Management

Version is stored as a variable in cmd/taildog/main.go. Use `make version` to extract current version, `make tag` to create git tags.

