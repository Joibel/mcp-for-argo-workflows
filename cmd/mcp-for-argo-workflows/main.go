// Package main is the entry point for the MCP server for Argo Workflows.
package main

import (
	"context"
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

func run() error {
	// Create the MCP server with name and version
	srv := server.NewServer(serverName, version.Version)

	// TODO: Register tools (will be implemented in future issues)
	slog.Info("MCP server created", "name", serverName, "version", version.Version)

	// Start the server with stdio transport
	// This will block until interrupted (SIGINT/SIGTERM)
	return srv.RunStdio(context.Background())
}
