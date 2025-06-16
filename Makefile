.PHONY: build test clean lint install dev run cross-build release help

# Build configuration
BUILD_DIR := build
BINARY_NAME := axon
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X main.version=$(VERSION) -s -w"
GOFLAGS := -trimpath

# Go build configuration
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

# Default target
default: build

## Build the binary
build:
	@echo "Building $(BINARY_NAME) v$(VERSION) for $(GOOS)/$(GOARCH)..."
	@mkdir -p $(BUILD_DIR)
	go build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

## Run all tests
test:
	@echo "Running tests..."
	go test -v ./...

## Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## Run tests with race detection
test-race:
	@echo "Running tests with race detection..."
	go test -race -v ./...

## Run benchmarks
bench:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

## Lint the code
lint:
	@echo "Running linter..."
	golangci-lint run

## Format the code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	goimports -w .

## Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	rm -f *.prof

## Install the binary to GOPATH/bin
install:
	@echo "Installing $(BINARY_NAME)..."
	go install $(GOFLAGS) $(LDFLAGS) .

## Run in development mode
dev:
	@echo "Running in development mode..."
	go run $(LDFLAGS) . --debug

## Run the built binary
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME)

## Cross-compile for multiple platforms
cross-build:
	@echo "Cross-compiling for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	# Linux AMD64
	GOOS=linux GOARCH=amd64 go build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .
	# Linux ARM64
	GOOS=linux GOARCH=arm64 go build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 .
	# macOS AMD64
	GOOS=darwin GOARCH=amd64 go build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .
	# macOS ARM64 (Apple Silicon)
	GOOS=darwin GOARCH=arm64 go build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .
	# Windows AMD64
	GOOS=windows GOARCH=amd64 go build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .
	# FreeBSD AMD64
	GOOS=freebsd GOARCH=amd64 go build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-freebsd-amd64 .
	@echo "Cross-compilation complete. Binaries in $(BUILD_DIR)/"

## Create checksums for release
checksums: cross-build
	@echo "Generating checksums..."
	cd $(BUILD_DIR) && sha256sum * > checksums.txt
	@echo "Checksums generated: $(BUILD_DIR)/checksums.txt"

## Prepare release package
release: clean cross-build checksums
	@echo "Preparing release package..."
	@mkdir -p $(BUILD_DIR)/release
	cp README.md $(BUILD_DIR)/release/
	cp DOCUMENTATION.md $(BUILD_DIR)/release/
	cp LICENSE $(BUILD_DIR)/release/ 2>/dev/null || echo "No LICENSE file found"
	@echo "Release package prepared in $(BUILD_DIR)/release/"

## Check dependencies
deps:
	@echo "Checking dependencies..."
	go mod verify
	go mod tidy

## Update dependencies
update-deps:
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy

## Generate documentation
docs:
	@echo "Generating documentation..."
	godoc -http=:6060 &
	@echo "Documentation server started at http://localhost:6060"

## Profile memory usage
profile-mem:
	@echo "Profiling memory usage..."
	go test -memprofile=mem.prof -bench=. ./...
	go tool pprof mem.prof

## Profile CPU usage
profile-cpu:
	@echo "Profiling CPU usage..."
	go test -cpuprofile=cpu.prof -bench=. ./...
	go tool pprof cpu.prof

## Check for vulnerabilities
security:
	@echo "Checking for security vulnerabilities..."
	govulncheck ./...

## Setup development environment
setup-dev:
	@echo "Setting up development environment..."
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
	@echo "Development tools installed"

## Display help
help:
	@echo "Available targets:"
	@echo "  build         - Build the binary"
	@echo "  test          - Run all tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  test-race     - Run tests with race detection"
	@echo "  bench         - Run benchmarks"
	@echo "  lint          - Run linter"
	@echo "  fmt           - Format code"
	@echo "  clean         - Clean build artifacts"
	@echo "  install       - Install binary to GOPATH/bin"
	@echo "  dev           - Run in development mode"
	@echo "  run           - Build and run binary"
	@echo "  cross-build   - Cross-compile for multiple platforms"
	@echo "  checksums     - Generate checksums for binaries"
	@echo "  release       - Prepare complete release package"
	@echo "  deps          - Check dependencies"
	@echo "  update-deps   - Update dependencies"
	@echo "  docs          - Start documentation server"
	@echo "  profile-mem   - Profile memory usage"
	@echo "  profile-cpu   - Profile CPU usage"
	@echo "  security      - Check for vulnerabilities"
	@echo "  setup-dev     - Setup development environment"
	@echo "  help          - Show this help message"

