package argo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient_NilConfig(t *testing.T) {
	client, err := NewClient(context.Background(), nil)
	require.Error(t, err)
	assert.Nil(t, client)
	assert.Contains(t, err.Error(), "config cannot be nil")
}

func TestNewClient_InvalidKubeconfig(t *testing.T) {
	// Test that an invalid kubeconfig path returns an error instead of panicking
	config := &Config{
		Kubeconfig: "/nonexistent/path/to/kubeconfig",
		Namespace:  "default",
	}

	client, err := NewClient(context.Background(), config)
	require.Error(t, err)
	assert.Nil(t, client)
	assert.Contains(t, err.Error(), "failed to create Argo API client")
}

func TestClient_DefaultNamespace(t *testing.T) {
	// Test the DefaultNamespace method by creating a client struct directly
	// (bypassing NewClient which requires a real connection)
	client := &Client{
		config: &Config{
			Namespace: "test-namespace",
		},
	}

	assert.Equal(t, "test-namespace", client.DefaultNamespace())
}

func TestClient_IsArgoServerMode(t *testing.T) {
	tests := []struct {
		name       string
		argoServer string
		expected   bool
	}{
		{
			name:       "with argo server",
			argoServer: "localhost:2746",
			expected:   true,
		},
		{
			name:       "without argo server",
			argoServer: "",
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{
				config: &Config{
					ArgoServer: tt.argoServer,
				},
			}
			assert.Equal(t, tt.expected, client.IsArgoServerMode())
		})
	}
}
