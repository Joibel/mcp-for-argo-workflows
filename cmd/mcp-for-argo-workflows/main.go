// Package main is the entry point for the MCP server for Argo Workflows.
package main

import (
	"fmt"
	"os"

	"github.com/Joibel/mcp-for-argo-workflows/internal/server"
	"github.com/Joibel/mcp-for-argo-workflows/internal/version"
)

const serverName = "mcp-for-argo-workflows"

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
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
	fmt.Fprintf(os.Stderr, "MCP server created: %s version %s\n", serverName, version.Version)

	// TODO: Start transport (stdio/HTTP) - will be implemented in PIP-11
	fmt.Fprintf(os.Stderr, "Transport setup not yet implemented (PIP-11)\n")

	return nil
}
