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

// waitForServer waits until the server is reachable or times out.
func waitForServer(t *testing.T, addr string) {
	t.Helper()
	client := &http.Client{Timeout: 100 * time.Millisecond}
	require.Eventually(t, func() bool {
		resp, err := client.Get(fmt.Sprintf("http://%s/", addr))
		if err != nil {
			return false
		}
		_ = resp.Body.Close()
		return true
	}, 500*time.Millisecond, 10*time.Millisecond, "server should become reachable")
}

// getAvailableAddr returns an available address for testing.
func getAvailableAddr(t *testing.T) string {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	addr := listener.Addr().String()
	require.NoError(t, listener.Close())
	return addr
}

// TestRunHTTP_BasicStartup tests that the HTTP server starts and listens successfully.
func TestRunHTTP_BasicStartup(t *testing.T) {
	srv := NewServer("test-server", "1.0.0")

	// Pick an addr we can actually probe
	addr := getAvailableAddr(t)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Start server in goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- srv.RunHTTP(ctx, addr)
	}()

	// Wait for server to become reachable
	waitForServer(t, addr)

	// Server should still be running
	select {
	case err := <-errChan:
		require.FailNowf(t, "server exited prematurely", "%v", err)
	default:
		// Server is still running, this is expected
	}

	// Cancel context to trigger shutdown
	cancel()

	// Server should shut down gracefully
	err := <-errChan
	assert.NoError(t, err, "server should shut down without error")
}

// TestRunHTTP_GracefulShutdown tests that cancelling the context triggers graceful shutdown.
func TestRunHTTP_GracefulShutdown(t *testing.T) {
	srv := NewServer("test-server", "1.0.0")
	addr := getAvailableAddr(t)

	ctx, cancel := context.WithCancel(context.Background())

	// Start server
	errChan := make(chan error, 1)
	go func() {
		errChan <- srv.RunHTTP(ctx, addr)
	}()

	// Wait for server to become reachable before testing shutdown
	waitForServer(t, addr)

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
		setup   func(t *testing.T) (string, func())
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
			name: "port in use",
			setup: func(t *testing.T) (string, func()) {
				t.Helper()
				listener, err := net.Listen("tcp", "127.0.0.1:0")
				require.NoError(t, err)
				return listener.Addr().String(), func() { _ = listener.Close() }
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr := tt.addr
			if tt.setup != nil {
				var cleanup func()
				addr, cleanup = tt.setup(t)
				defer cleanup()
			}

			srv := NewServer("test-server", "1.0.0")

			// Use a short timeout since we expect failure
			ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
			defer cancel()

			runErr := srv.RunHTTP(ctx, addr)

			if tt.wantErr {
				assert.Error(t, runErr, "expected error for address: %s", addr)
			} else {
				assert.NoError(t, runErr)
			}
		})
	}
}

// TestRunHTTP_HTTPServerConfiguration tests that the HTTP server has correct configuration.
func TestRunHTTP_HTTPServerConfiguration(t *testing.T) {
	srv := NewServer("test-server", "1.0.0")
	addr := getAvailableAddr(t)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Start server
	errChan := make(chan error, 1)
	go func() {
		errChan <- srv.RunHTTP(ctx, addr)
	}()

	// Verify server is listening and responds
	client := &http.Client{Timeout: 100 * time.Millisecond}
	require.Eventually(t, func() bool {
		resp, err := client.Get(fmt.Sprintf("http://%s/", addr))
		if err != nil {
			return false
		}
		_ = resp.Body.Close()
		return resp.StatusCode != 0
	}, 500*time.Millisecond, 10*time.Millisecond, "server should respond to HTTP requests")

	// Cancel to trigger shutdown
	cancel()
	runErr := <-errChan
	assert.NoError(t, runErr)
}

// TestRunHTTP_MultipleShutdowns tests that multiple cancellations don't cause issues.
func TestRunHTTP_MultipleShutdowns(t *testing.T) {
	srv := NewServer("test-server", "1.0.0")
	addr := getAvailableAddr(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errChan := make(chan error, 1)
	go func() {
		errChan <- srv.RunHTTP(ctx, addr)
	}()

	// Wait for server to become reachable
	waitForServer(t, addr)

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
	addr := getAvailableAddr(t)

	// Create an already-cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Server should start and then immediately shut down
	runErr := srv.RunHTTP(ctx, addr)

	// This should either return no error (clean shutdown) or context.Canceled
	// Both are acceptable behaviors
	assert.True(t, runErr == nil || errors.Is(runErr, context.Canceled),
		"expected nil or context.Canceled, got: %v", runErr)
}

// TestRunHTTP_ConcurrentRequests tests that the server can handle multiple concurrent requests.
func TestRunHTTP_ConcurrentRequests(t *testing.T) {
	srv := NewServer("test-server", "1.0.0")
	addr := getAvailableAddr(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	errChan := make(chan error, 1)
	go func() {
		errChan <- srv.RunHTTP(ctx, addr)
	}()

	// Wait for server to become reachable
	waitForServer(t, addr)

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

	// Cancel to trigger shutdown
	cancel()
	runErr := <-errChan
	assert.NoError(t, runErr)
}

// TestRunHTTP_ServerStartsNormally is a smoke test verifying the HTTP server
// starts and stops without errors. This indirectly confirms that server configuration
// (including ReadHeaderTimeout) doesn't prevent normal operation.
func TestRunHTTP_ServerStartsNormally(t *testing.T) {
	srv := NewServer("test-server", "1.0.0")
	addr := getAvailableAddr(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errChan := make(chan error, 1)
	go func() {
		errChan <- srv.RunHTTP(ctx, addr)
	}()

	// Wait for server to become reachable
	waitForServer(t, addr)

	// Server started normally, trigger shutdown
	cancel()
	runErr := <-errChan
	assert.NoError(t, runErr, "server should start and stop normally")
}

// TestRunHTTP_QuickSuccession tests starting servers on the same port in quick succession.
func TestRunHTTP_QuickSuccession(t *testing.T) {
	addr := getAvailableAddr(t)

	// First server
	srv1 := NewServer("test-server-1", "1.0.0")
	ctx1, cancel1 := context.WithCancel(context.Background())

	errChan1 := make(chan error, 1)
	go func() {
		errChan1 <- srv1.RunHTTP(ctx1, addr)
	}()

	// Wait for first server to become reachable
	waitForServer(t, addr)

	// Shut down first server
	cancel1()
	runErr1 := <-errChan1
	assert.NoError(t, runErr1)

	// Wait for port to be released (TCP TIME_WAIT state)
	// Use a simple polling approach to check when the port is available
	require.Eventually(t, func() bool {
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			return false
		}
		_ = listener.Close()
		return true
	}, 5*time.Second, 100*time.Millisecond, "port should be released after first server shutdown")

	// Second server on same port - now that port is available
	srv2 := NewServer("test-server-2", "1.0.0")
	ctx2, cancel2 := context.WithCancel(context.Background())
	defer cancel2()

	errChan2 := make(chan error, 1)
	go func() {
		errChan2 <- srv2.RunHTTP(ctx2, addr)
	}()

	// Wait for second server to become reachable
	waitForServer(t, addr)

	// Verify second server is running by making a request
	client := &http.Client{Timeout: 100 * time.Millisecond}
	resp, err := client.Get(fmt.Sprintf("http://%s/", addr))
	require.NoError(t, err, "second server should respond to requests")
	_ = resp.Body.Close()

	// Shut down second server
	cancel2()
	runErr2 := <-errChan2
	assert.NoError(t, runErr2)
}
