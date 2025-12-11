// Package main is the entry point for the MCP server for Argo Workflows.
package main

import (
	"log/slog"
	"os"

	"github.com/Joibel/mcp-for-argo-workflows/internal/server"
	"github.com/Joibel/mcp-for-argo-workflows/internal/version"
)

const serverName = "mcp-for-argo-workflows"

func main() {
	// Configure structured logging to stderr
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	if err := run(); err != nil {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}
}

//nolint:unparam // Will return errors in future implementations (PIP-11+)
func run() error {
	// Create the MCP server with name and version
	srv := server.NewServer(serverName, version.Version)

	// TODO: Register tools (will be implemented in future issues)
	// For now, just log that the server was created
	_ = srv
	slog.Info("MCP server created", "name", serverName, "version", version.Version)

	// TODO: Start transport (stdio/HTTP) - will be implemented in PIP-11
	slog.Info("transport setup not yet implemented", "issue", "PIP-11")

	return nil
}
