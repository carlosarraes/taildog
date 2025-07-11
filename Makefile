# Taildog - Datadog Log Tailing CLI Tool
.PHONY: help build install clean test lint fmt vet run dev release tag push check deps update-deps

# Default target
help: ## Show this help message
	@echo "Taildog - Datadog Log Tailing CLI Tool"
	@echo ""
	@echo "Available commands:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build commands
build: ## Build the binary and copy to ~/.local/bin
	go build -ldflags="-w -s" -o taildog ./cmd/taildog
	@mkdir -p ~/.local/bin
	@cp taildog ~/.local/bin/
	@echo "Copied to ~/.local/bin/taildog"

build-dev: ## Build without optimizations (faster, for development)
	go build -o taildog ./cmd/taildog

install: build ## Build and install to GOPATH/bin
	go install ./cmd/taildog

# Cross-platform builds
build-all: ## Build for all platforms
	GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o dist/taildog-linux-x86_64 ./cmd/taildog
	GOOS=linux GOARCH=arm64 go build -ldflags="-w -s" -o dist/taildog-linux-aarch64 ./cmd/taildog
	GOOS=darwin GOARCH=amd64 go build -ldflags="-w -s" -o dist/taildog-darwin-x86_64 ./cmd/taildog
	GOOS=darwin GOARCH=arm64 go build -ldflags="-w -s" -o dist/taildog-darwin-aarch64 ./cmd/taildog

# Development commands
run: build-dev ## Build and run with example query
	./taildog --help

dev: build-dev ## Quick development cycle (build + show help)
	./taildog --help

test: ## Run tests
	go test -v ./...

test-race: ## Run tests with race detection
	go test -race -v ./...

test-cover: ## Run tests with coverage
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

fmt: ## Format code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

check: fmt vet test ## Run all checks (format, vet, test)

# Dependency management
deps: ## Download dependencies
	go mod download

tidy: ## Tidy dependencies
	go mod tidy

update-deps: ## Update dependencies
	go get -u ./...
	go mod tidy

# Release commands
version: ## Show current version
	@grep 'var version' cmd/taildog/main.go | cut -d'"' -f2

release: clean check build ## Prepare for release (clean, check, build)
	@echo "Ready for release. Current version: $$(make version)"
	@echo "To create release:"
	@echo "  1. Update version in cmd/taildog/main.go"
	@echo "  2. Run: make tag"
	@echo "  3. Run: make push"

tag: ## Create and push git tag for current version
	$(eval VERSION := $(shell make version))
	git tag -a v$(VERSION) -m "Release v$(VERSION)"
	@echo "Created tag v$(VERSION)"
	@echo "Run 'make push' to push tag and trigger release"

push: ## Push current branch and tags
	git push origin
	git push origin --tags

# Utility commands
clean: ## Clean build artifacts
	rm -f taildog
	rm -rf dist/
	rm -f coverage.out coverage.html

demo: build-dev ## Run demo commands
	@echo "=== Taildog Demo ==="
	@echo "1. Version info:"
	./taildog --version
	@echo ""
	@echo "2. Help output:"
	./taildog --help
	@echo ""
	@echo "3. Example query (requires DD_API_KEY and DD_APPLICATION_KEY):"
	@echo "   ./taildog \"service:example\" --dry-run"

# Development helpers
watch: ## Watch for changes and rebuild (requires entr)
	@if command -v entr >/dev/null 2>&1; then \
		find . -name "*.go" | entr -r make dev; \
	else \
		echo "entr not installed. Install with your package manager"; \
		echo "  macOS: brew install entr"; \
		echo "  Ubuntu: apt install entr"; \
	fi

size: build ## Show binary size
	@ls -lh taildog | awk '{print "Binary size: " $$5}'

# Create dist directory
dist:
	mkdir -p dist