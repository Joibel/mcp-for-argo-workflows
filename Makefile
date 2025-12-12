# Project variables
BINARY_NAME := mcp-for-argo-workflows
MODULE := github.com/Joibel/mcp-for-argo-workflows
DIST_DIR := dist

# Build variables
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
LDFLAGS := -s -w \
	-X $(MODULE)/internal/version.Version=$(VERSION) \
	-X $(MODULE)/internal/version.Commit=$(COMMIT) \
	-X $(MODULE)/internal/version.BuildTime=$(BUILD_TIME)

# Go variables
GO := go
GOFMT := gofmt
GOIMPORTS := goimports
GOLANGCI_LINT := golangci-lint

# Source files for dependency tracking (exclude dist/ and vendor/)
GO_FILES := $(shell find . -name '*.go' -type f -not -path './dist/*' -not -path './vendor/*')
GO_MOD := go.mod go.sum

# Platform-specific binary paths
DIST_DARWIN_AMD64 := $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64
DIST_DARWIN_ARM64 := $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64
DIST_LINUX_AMD64 := $(DIST_DIR)/$(BINARY_NAME)-linux-amd64
DIST_LINUX_ARM64 := $(DIST_DIR)/$(BINARY_NAME)-linux-arm64
DIST_WINDOWS_AMD64 := $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe
DIST_CHECKSUMS := $(DIST_DIR)/checksums.txt

# All distribution binaries
DIST_BINARIES := $(DIST_DARWIN_AMD64) $(DIST_DARWIN_ARM64) $(DIST_LINUX_AMD64) $(DIST_LINUX_ARM64) $(DIST_WINDOWS_AMD64)

.PHONY: all test lint lint-fix fmt vet clean tools help build-all dist-clean

# Default target
all: fmt vet lint test $(DIST_LINUX_AMD64)

## test: Run tests with race detection and coverage
test:
	@echo "Running tests..."
	$(GO) test -race -coverprofile=coverage.out -covermode=atomic ./internal/...
	@echo "Coverage report: coverage.out"

## lint: Run golangci-lint
lint:
	@echo "Running linter..."
	$(GOLANGCI_LINT) run ./...

## lint-fix: Run golangci-lint with auto-fix
lint-fix:
	@echo "Running linter with auto-fix..."
	$(GOLANGCI_LINT) run --fix ./...

## fmt: Run gofmt and goimports
fmt:
	@echo "Formatting code..."
	$(GOFMT) -s -w .
	@if command -v $(GOIMPORTS) >/dev/null 2>&1; then \
		$(GOIMPORTS) -w -local $(MODULE) .; \
	else \
		echo "goimports not installed, skipping import formatting"; \
	fi

## vet: Run go vet
vet:
	@echo "Running go vet..."
	$(GO) vet ./...

## clean: Remove build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(DIST_DIR)/
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

## tools: Install development dependencies
tools:
	@echo "Installing development tools..."
	$(GO) install golang.org/x/tools/cmd/goimports@latest
	@echo "Note: Install golangci-lint from https://golangci-lint.run/welcome/install/"

## help: Show this help message
help:
	@echo "Available targets:"
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/## /  /'

# =============================================================================
# Cross-compilation targets (real file targets with dependencies)
# =============================================================================

## build-all: Build binaries for all platforms
build-all: $(DIST_BINARIES) $(DIST_CHECKSUMS)
	@echo "All platform builds complete. Binaries in $(DIST_DIR)/"

# macOS Intel
$(DIST_DARWIN_AMD64): $(GO_FILES) $(GO_MOD)
	@echo "Building for darwin/amd64..."
	@mkdir -p $(DIST_DIR)
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GO) build -ldflags="$(LDFLAGS)" -o $@ ./cmd/$(BINARY_NAME)

# macOS Apple Silicon
$(DIST_DARWIN_ARM64): $(GO_FILES) $(GO_MOD)
	@echo "Building for darwin/arm64..."
	@mkdir -p $(DIST_DIR)
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GO) build -ldflags="$(LDFLAGS)" -o $@ ./cmd/$(BINARY_NAME)

# Linux x86_64
$(DIST_LINUX_AMD64): $(GO_FILES) $(GO_MOD)
	@echo "Building for linux/amd64..."
	@mkdir -p $(DIST_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -ldflags="$(LDFLAGS)" -o $@ ./cmd/$(BINARY_NAME)

# Linux ARM64
$(DIST_LINUX_ARM64): $(GO_FILES) $(GO_MOD)
	@echo "Building for linux/arm64..."
	@mkdir -p $(DIST_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GO) build -ldflags="$(LDFLAGS)" -o $@ ./cmd/$(BINARY_NAME)

# Windows x86_64
$(DIST_WINDOWS_AMD64): $(GO_FILES) $(GO_MOD)
	@echo "Building for windows/amd64..."
	@mkdir -p $(DIST_DIR)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO) build -ldflags="$(LDFLAGS)" -o $@ ./cmd/$(BINARY_NAME)

## checksums: Generate SHA256 checksums for all binaries
$(DIST_CHECKSUMS): $(DIST_BINARIES)
	@echo "Generating checksums..."
	@cd $(DIST_DIR) && shasum -a 256 $(BINARY_NAME)-* > checksums.txt
	@echo "Checksums written to $(DIST_CHECKSUMS)"

## dist-clean: Remove distribution artifacts
dist-clean:
	@echo "Cleaning dist directory..."
	@rm -rf $(DIST_DIR)/
