# taildog

A modern CLI tool to tail Datadog logs in real-time, inspired by `tail -f` but for your cloud infrastructure.

## Features

- **Real-time log tailing** from Datadog with live updates
- **Multiple query support** - run multiple instances for different services/environments
- **File output** with optional rotation and buffering
- **Flexible formatting** - JSON for searchability, text for readability
- **Smart polling** with cursor-based pagination to avoid duplicates
- **Colorized output** for better terminal visibility
- **Graceful shutdown** with Ctrl+C handling

## Installation

### Quick Install (Linux/macOS)
```bash
curl -sSf https://raw.githubusercontent.com/carlosarraes/taildog/main/install.sh | sh
```

### Manual Download
Download the binary for your platform from the [releases page](https://github.com/carlosarraes/taildog/releases).

## Quick Start

```bash
# Set your Datadog credentials
export DD_API_KEY="your-api-key"
export DD_APPLICATION_KEY="your-app-key"

# Tail all logs to stdout
taildog

# Tail specific service logs to stdout
taildog "service:my-app"

# Tail and write to file
taildog "service:my-app AND @http.status_code:>=400" --output errors.log

# Multiple services
taildog "service:app1 OR service:app2"
```

## Usage Examples

### Basic Tailing

```bash
# Tail ALL logs (no filter)
taildog

# Tail all logs from a service
taildog "service:my-api"

# Tail with environment filter
taildog "env:prod service:api"

# Tail specific environment and source
taildog "env:homolog source:django"

# Tail errors only
taildog "service:my-app AND status:error"
```

### Output to File

```bash
# Write logs to file
taildog "service:my-app" --output app.log

# JSON format for searchability
taildog "service:my-app" --output app.json --format json

# With file rotation
taildog "service:my-app" --output app.log --rotate-size 100MB
```

### Advanced Usage

```bash
# Custom polling interval
taildog "service:my-app" --interval 10s

# Fetch last hour, then follow
taildog "service:my-app" --time-range 1h

# Buffer writes for performance
taildog "service:my-app" --output app.log --buffer 100

# One-time fetch (no follow)
taildog "service:my-app" --follow=false
```

### Multiple Instances

```bash
# Terminal 1: Production API logs
taildog "env:prod service:api" --output prod-api.log

# Terminal 2: Development Aspose service
taildog "env:dev service:aspose" --output dev-aspose.log

# Terminal 3: All production errors
taildog "env:prod AND status:error" --output prod-errors.log
```

## Configuration

### Environment Variables

```bash
DD_API_KEY           # Required: Datadog API key
DD_APPLICATION_KEY   # Required: Datadog Application key
DD_SITE              # Optional: Datadog site (default: datadoghq.com)
```

### Command Line Options

```
Usage: taildog <query> [flags]

Arguments:
  <query>    Datadog query (e.g. "service:my-app")

Flags:
  -o, --output STRING       Output file path
  -f, --follow              Follow logs like tail -f (default: true)
  -i, --interval DURATION   Polling interval (default: 5s)
  -t, --time-range STRING   Initial time range to fetch (default: 15m)
  -F, --format STRING       Output format: text|json (default: text)
  -R, --rotate-size STRING  Rotate file at size (e.g. 100MB)
  -b, --buffer INT          Buffer N logs before writing (default: 1)
  --max-retries INT         Maximum retry attempts (default: 5)
  -h, --help               Show help
```

## Output Formats

### Text Format (Default)

```
2024-01-01T10:00:00Z my-service INFO Starting application
2024-01-01T10:00:01Z my-service ERROR Database connection failed
```

### JSON Format

```json
{"timestamp":"2024-01-01T10:00:00Z","service":"my-service","level":"INFO","message":"Starting application"}
{"timestamp":"2024-01-01T10:00:01Z","service":"my-service","level":"ERROR","message":"Database connection failed"}
```

## Query Syntax

taildog uses Datadog's search syntax. Common patterns:

```bash
# Service filtering
taildog "service:my-app"

# Environment and service
taildog "env:prod service:api"

# Error logs only
taildog "status:error"

# HTTP status codes
taildog "@http.status_code:>=400"

# Time-based (for initial fetch)
taildog "service:my-app" --time-range 1h

# Complex queries
taildog "service:my-app AND (status:error OR @http.status_code:>=400)"
```

## Use Cases

- **Development**: Monitor your service logs during development
- **Production Monitoring**: Keep an eye on production errors and performance
- **Debugging**: Tail specific error patterns or trace IDs
- **Log Analysis**: Export logs to files for further analysis
- **Multi-environment**: Run multiple instances for different environments

## Development

### Implementation Priority

**Phase 1 (MVP)**: Basic real-time log tailing

- Environment variable authentication (DD_API_KEY, DD_APP_KEY)
- Simple query parsing and API calls
- Raw log output to stdout
- Basic error handling with configurable retries
- Example: `taildog "env:homolog source:django"`

**Phase 2**: Enhanced functionality

- TOML configuration support (`~/.config/taildog/config.toml`)
- Multiple output formats (text/JSON)
- File output with rotation and buffering
- Advanced query support and filtering

**Phase 3**: Advanced features

- Performance optimizations
- Multiple output destinations
- Metrics and monitoring
- Plugin system

### Project Structure

Following standard Go project layout:

```
cmd/taildog/main.go          # CLI entry point
internal/client/             # Datadog API client
internal/config/             # Configuration management
internal/output/             # Output formatting and writing
pkg/types/                   # Shared types and interfaces
```

### Code Guidelines

- Maximum 400 lines per source file (tests excluded)
- Modular design with clear separation of concerns
- TDD approach with comprehensive test coverage
- Follow Go best practices and conventions

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

MIT License - see LICENSE file for details

