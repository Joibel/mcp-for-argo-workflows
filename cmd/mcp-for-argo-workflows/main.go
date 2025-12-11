// Package main is the entry point for the MCP server for Argo Workflows.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/Joibel/mcp-for-argo-workflows/internal/config"
	"github.com/Joibel/mcp-for-argo-workflows/internal/server"
	"github.com/Joibel/mcp-for-argo-workflows/internal/version"
)

const serverName = "mcp-for-argo-workflows"

func main() {
	// Configure structured logging to stderr
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	if err := run(context.Background()); err != nil {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	// Parse configuration from CLI flags and environment variables
	cfg, err := config.NewFromFlags()
	if err != nil {
		return fmt.Errorf("failed to parse configuration: %w", err)
	}

	// Create the MCP server with name and version
	srv := server.NewServer(serverName, version.Version)

	// TODO: Register tools (will be implemented in future issues)
	slog.Info("MCP server created",
		"name", serverName,
		"version", version.Version,
		"transport", cfg.Transport,
		"namespace", cfg.Namespace,
	)

	// Start the server with the configured transport
	if cfg.IsHTTPTransport() {
		slog.Info("starting HTTP transport", "addr", cfg.HTTPAddr)
		return srv.RunHTTP(ctx, cfg.HTTPAddr)
	}

	// Default to stdio transport
	return srv.RunStdio(ctx)
}
