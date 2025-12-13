package server

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RunHTTP runs the MCP server with HTTP/SSE transport.
// It handles graceful shutdown on SIGINT and SIGTERM signals.
func (s *Server) RunHTTP(ctx context.Context, addr string) error {
	// Create a context that cancels on SIGINT or SIGTERM
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	slog.Info("starting MCP server", "transport", "http", "addr", addr)

	// Create an SSE handler that returns our MCP server for each new session
	handler := mcp.NewSSEHandler(func(_ *http.Request) *mcp.Server {
		return s.mcp
	}, nil)

	// Create HTTP server with timeouts to prevent Slowloris attacks
	httpServer := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
	}

	// Start HTTP server in a goroutine
	errChan := make(chan error, 1)
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
		close(errChan)
	}()

	// Wait for shutdown signal or error
	select {
	case <-ctx.Done():
		slog.Info("shutting down HTTP server")
		//nolint:contextcheck // Use fresh context for graceful shutdown after cancellation
		if err := httpServer.Shutdown(context.Background()); err != nil {
			return err
		}
	case err := <-errChan:
		return err
	}

	slog.Info("MCP server shutdown gracefully", "transport", "http")
	return nil
}
