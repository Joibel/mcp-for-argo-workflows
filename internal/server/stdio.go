// Package server provides the MCP server implementation for Argo Workflows.
package server

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RunStdio runs the MCP server with stdio transport.
// It handles graceful shutdown on SIGINT and SIGTERM signals.
func (s *Server) RunStdio(ctx context.Context) error {
	// Create a context that cancels on SIGINT or SIGTERM
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	slog.Info("starting MCP server", "transport", "stdio")

	// Run the server with stdio transport
	// This blocks until the transport is closed or context is cancelled
	transport := &mcp.StdioTransport{}
	if err := s.mcp.Run(ctx, transport); err != nil {
		return err
	}

	slog.Info("MCP server shutdown gracefully", "transport", "stdio")
	return nil
}
