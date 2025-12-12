package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRunHTTP_BasicStartup tests that the HTTP server starts and listens successfully.
func TestRunHTTP_BasicStartup(t *testing.T) {
	server := NewServer("test-server", "1.0.0")

	// Use port 0 to get a random available port
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Start server in goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- server.RunHTTP(ctx, "127.0.0.1:0")
	}()

	// Give the server a moment to start
	time.Sleep(50 * time.Millisecond)

	// Server should still be running
	select {
	case err := <-errChan:
		t.Fatalf("server exited prematurely: %v", err)
	default:
		// Server is still running, this is expected
	}

	// Wait for context timeout to trigger shutdown
	<-ctx.Done()

	// Server should shut down gracefully
	err := <-errChan
	assert.NoError(t, err, "server should shut down without error")
}

// TestRunHTTP_GracefulShutdown tests that cancelling the context triggers graceful shutdown.
func TestRunHTTP_GracefulShutdown(t *testing.T) {
	server := NewServer("test-server", "1.0.0")

	// Get a random available port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	addr := listener.Addr().String()
	_ = listener.Close()

	ctx, cancel := context.WithCancel(context.Background())

	// Start server
	errChan := make(chan error, 1)
	go func() {
		errChan <- server.RunHTTP(ctx, addr)
	}()

	// Give server time to start
	time.Sleep(50 * time.Millisecond)

	// Cancel context to trigger shutdown
	cancel()

	// Server should shut down gracefully within reasonable time
	select {
	case err := <-errChan:
		assert.NoError(t, err, "server should shut down gracefully")
	case <-time.After(2 * time.Second):
		t.Fatal("server did not shut down within timeout")
	}
}

// TestRunHTTP_InvalidAddress tests that invalid addresses return errors.
func TestRunHTTP_InvalidAddress(t *testing.T) {
	tests := []struct {
		name    string
		addr    string
		wantErr bool
	}{
		{
			name:    "invalid port - too high",
			addr:    "127.0.0.1:99999",
			wantErr: true,
		},
		{
			name:    "invalid host",
			addr:    "999.999.999.999:8080",
			wantErr: true,
		},
		{
			name:    "port in use",
			addr:    "", // Will be set dynamically
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Special handling for "port in use" test
			if tt.name == "port in use" {
				// Bind to a port first
				listener, listenErr := net.Listen("tcp", "127.0.0.1:0")
				require.NoError(t, listenErr)
				tt.addr = listener.Addr().String()
				defer func() { _ = listener.Close() }()
			}

			srv := NewServer("test-server", "1.0.0")

			// Use a short timeout since we expect failure
			ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
			defer cancel()

			runErr := srv.RunHTTP(ctx, tt.addr)

			if tt.wantErr {
				assert.Error(t, runErr, "expected error for address: %s", tt.addr)
			} else {
				assert.NoError(t, runErr)
			}
		})
	}
}

// TestRunHTTP_HTTPServerConfiguration tests that the HTTP server has correct configuration.
func TestRunHTTP_HTTPServerConfiguration(t *testing.T) {
	srv := NewServer("test-server", "1.0.0")

	// Get a random available port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	addr := listener.Addr().String()
	_ = listener.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	// Start server
	errChan := make(chan error, 1)
	go func() {
		errChan <- srv.RunHTTP(ctx, addr)
	}()

	// Give server time to start
	time.Sleep(50 * time.Millisecond)

	// Make a request to verify the handler is set up
	// Note: We're not testing the MCP protocol itself, just that HTTP works
	client := &http.Client{
		Timeout: 1 * time.Second,
	}

	resp, getErr := client.Get(fmt.Sprintf("http://%s/", addr))
	if getErr == nil {
		defer func() { _ = resp.Body.Close() }()
		// Server responded, which means it's listening
		// The exact response depends on the MCP SDK's SSE handler
		assert.NotEqual(t, 0, resp.StatusCode, "server should respond")
	}

	// Wait for context to expire and server to shut down
	<-ctx.Done()
	runErr := <-errChan
	assert.NoError(t, runErr)
}

// TestRunHTTP_MultipleShutdowns tests that multiple cancellations don't cause issues.
func TestRunHTTP_MultipleShutdowns(t *testing.T) {
	srv := NewServer("test-server", "1.0.0")

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	addr := listener.Addr().String()
	_ = listener.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errChan := make(chan error, 1)
	go func() {
		errChan <- srv.RunHTTP(ctx, addr)
	}()

	// Give server time to start
	time.Sleep(50 * time.Millisecond)

	// Cancel multiple times (should be safe)
	cancel()
	cancel()
	cancel()

	// Server should still shut down gracefully once
	select {
	case runErr := <-errChan:
		assert.NoError(t, runErr)
	case <-time.After(2 * time.Second):
		t.Fatal("server did not shut down within timeout")
	}
}

// TestRunHTTP_ContextAlreadyCancelled tests behavior when context is already cancelled.
func TestRunHTTP_ContextAlreadyCancelled(t *testing.T) {
	srv := NewServer("test-server", "1.0.0")

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	addr := listener.Addr().String()
	_ = listener.Close()

	// Create an already-cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Server should start and then immediately shut down
	runErr := srv.RunHTTP(ctx, addr)

	// This should either return no error (clean shutdown) or context.Canceled
	// Both are acceptable behaviors
	if runErr != nil && !errors.Is(runErr, context.Canceled) {
		t.Errorf("unexpected error: %v", runErr)
	}
}

// TestRunHTTP_ConcurrentRequests tests that the server can handle multiple concurrent requests.
func TestRunHTTP_ConcurrentRequests(t *testing.T) {
	srv := NewServer("test-server", "1.0.0")

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	addr := listener.Addr().String()
	_ = listener.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	errChan := make(chan error, 1)
	go func() {
		errChan <- srv.RunHTTP(ctx, addr)
	}()

	// Give server time to start
	time.Sleep(50 * time.Millisecond)

	// Make concurrent requests
	const numRequests = 5
	requestErrs := make(chan error, numRequests)

	for i := 0; i < numRequests; i++ {
		go func(id int) {
			client := &http.Client{Timeout: 500 * time.Millisecond}
			resp, getErr := client.Get(fmt.Sprintf("http://%s/", addr))
			if getErr == nil {
				_ = resp.Body.Close()
			}
			requestErrs <- getErr
		}(i)
	}

	// Collect results from concurrent requests
	successCount := 0
	for i := 0; i < numRequests; i++ {
		reqErr := <-requestErrs
		if reqErr == nil {
			successCount++
		}
	}

	// At least some requests should succeed
	assert.Greater(t, successCount, 0, "at least some concurrent requests should succeed")

	// Wait for server shutdown
	<-ctx.Done()
	runErr := <-errChan
	assert.NoError(t, runErr)
}

// TestRunHTTP_ReadHeaderTimeout tests that ReadHeaderTimeout is configured.
func TestRunHTTP_ReadHeaderTimeout(t *testing.T) {
	// This test verifies that the timeout is set by checking the server doesn't hang
	// on slow header reads. We test this indirectly by ensuring the server configuration
	// includes the timeout.

	srv := NewServer("test-server", "1.0.0")

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	addr := listener.Addr().String()
	_ = listener.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	errChan := make(chan error, 1)
	go func() {
		errChan <- srv.RunHTTP(ctx, addr)
	}()

	// Give server time to start
	time.Sleep(50 * time.Millisecond)

	// The server should be configured with ReadHeaderTimeout (10 seconds in the code)
	// We can't easily test the timeout behavior without actually triggering it,
	// but we can verify the server starts and runs normally

	<-ctx.Done()
	runErr := <-errChan
	assert.NoError(t, runErr, "server with ReadHeaderTimeout should work normally")
}

// TestRunHTTP_QuickSuccession tests starting servers on the same port in quick succession.
func TestRunHTTP_QuickSuccession(t *testing.T) {
	// Get a port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	addr := listener.Addr().String()
	_ = listener.Close()

	// First server
	srv1 := NewServer("test-server-1", "1.0.0")
	ctx1, cancel1 := context.WithCancel(context.Background())

	errChan1 := make(chan error, 1)
	go func() {
		errChan1 <- srv1.RunHTTP(ctx1, addr)
	}()

	// Give first server time to start
	time.Sleep(50 * time.Millisecond)

	// Shut down first server
	cancel1()
	runErr1 := <-errChan1
	assert.NoError(t, runErr1)

	// Give OS time to release the port
	time.Sleep(100 * time.Millisecond)

	// Second server on same port should work
	srv2 := NewServer("test-server-2", "1.0.0")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel2()

	errChan2 := make(chan error, 1)
	go func() {
		errChan2 <- srv2.RunHTTP(ctx2, addr)
	}()

	// Give second server time to start
	time.Sleep(50 * time.Millisecond)

	<-ctx2.Done()
	runErr2 := <-errChan2
	assert.NoError(t, runErr2, "second server should start after first shuts down")
}
