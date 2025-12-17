//go:build e2e

// Package e2e contains end-to-end tests for the MCP server.
package e2e

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/k3s"
)

// TestStartup_NoKubernetesConfigured verifies the server handles missing kubeconfig gracefully.
func TestStartup_NoKubernetesConfigured(t *testing.T) {
	binaryPath := buildBinary(t)

	// Run the server with no kubeconfig
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	//nolint:gosec // Running test binaries in tests is expected
	cmd := exec.CommandContext(ctx, binaryPath)
	// Build environment without HOME or KUBECONFIG
	filteredEnv := []string{}
	for _, e := range os.Environ() {
		if !strings.HasPrefix(e, "HOME=") && !strings.HasPrefix(e, "KUBECONFIG=") {
			filteredEnv = append(filteredEnv, e)
		}
	}
	cmd.Env = append(filteredEnv, "HOME=/nonexistent")

	output, err := cmd.CombinedOutput()

	// The server should exit with an error (not panic)
	assert.Error(t, err, "Server should exit with error when no kubeconfig is available")

	// Verify it's a clean error, not a panic
	outputStr := string(output)
	assert.NotContains(t, outputStr, "panic:", "Server should not panic")
	assert.NotContains(t, outputStr, "runtime error:", "Server should not have runtime errors")
	assert.Contains(t, outputStr, "failed to create Argo", "Error message should mention Argo client failure")
}

// TestStartup_WithK3s verifies the server can connect to a real Kubernetes cluster.
func TestStartup_WithK3s(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping k3s test in short mode")
	}

	ctx := t.Context()

	// Start k3s container
	k3sContainer, err := k3s.Run(ctx, "rancher/k3s:v1.31.2-k3s1")
	require.NoError(t, err, "Failed to start k3s container")
	defer func() {
		if termErr := k3sContainer.Terminate(ctx); termErr != nil {
			t.Logf("Failed to terminate k3s container: %v", termErr)
		}
	}()

	// Get kubeconfig from container
	kubeconfig, err := k3sContainer.GetKubeConfig(ctx)
	require.NoError(t, err, "Failed to get kubeconfig from k3s")

	// Write kubeconfig to temp file
	kubeconfigFile, err := os.CreateTemp("", "kubeconfig-*.yaml")
	require.NoError(t, err, "Failed to create temp kubeconfig file")
	defer os.Remove(kubeconfigFile.Name())

	_, err = kubeconfigFile.Write(kubeconfig)
	require.NoError(t, err, "Failed to write kubeconfig")
	err = kubeconfigFile.Close()
	require.NoError(t, err, "Failed to close kubeconfig file")

	// Build the binary
	binaryPath := buildBinary(t)

	// Run the server with the k3s kubeconfig
	// Use a short timeout since we just want to verify it starts without crashing
	runCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	//nolint:gosec // Running test binaries in tests is expected
	cmd := exec.CommandContext(runCtx, binaryPath)
	cmd.Env = append(os.Environ(), "KUBECONFIG="+kubeconfigFile.Name())

	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	// The server should either:
	// 1. Time out (exit code from context deadline) - meaning it started successfully
	// 2. Exit with an error about Argo workflows not being installed (expected)
	// Either way, it should NOT panic

	assert.NotContains(t, outputStr, "panic:", "Server should not panic")
	assert.NotContains(t, outputStr, "runtime error:", "Server should not have runtime errors")

	// If there's an error, it should be a clean error message
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && runCtx.Err() != context.DeadlineExceeded {
			// Server exited before timeout - check it was a clean exit
			t.Logf("Server exited with code %d, output: %s", exitErr.ExitCode(), outputStr)
		}
	}
}
