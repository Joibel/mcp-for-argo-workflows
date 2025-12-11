package argo

import (
	"os"
	"strconv"
)

// Config holds the configuration for connecting to Argo Workflows.
type Config struct {
	// ArgoServer is the Argo Server host:port (e.g., "localhost:2746").
	// When empty, the client will use direct Kubernetes API access.
	ArgoServer string

	// ArgoToken is the bearer token for authentication with Argo Server.
	ArgoToken string

	// Namespace is the default namespace for operations.
	Namespace string

	// Kubeconfig is the path to the kubeconfig file.
	// Used for direct Kubernetes API access when ArgoServer is empty.
	Kubeconfig string

	// Secure indicates whether to use TLS when connecting to Argo Server.
	// Only applies when ArgoServer is set.
	Secure bool

	// InsecureSkipVerify skips TLS certificate verification.
	// Only applies when ArgoServer is set and Secure is true.
	InsecureSkipVerify bool
}

// NewConfigFromEnv creates a Config from environment variables.
func NewConfigFromEnv() *Config {
	config := &Config{
		ArgoServer: os.Getenv("ARGO_SERVER"),
		ArgoToken:  os.Getenv("ARGO_TOKEN"),
		Namespace:  os.Getenv("ARGO_NAMESPACE"),
		Kubeconfig: os.Getenv("KUBECONFIG"),
		Secure:     true, // Default to secure
	}

	// Parse ARGO_SECURE if set
	if secureStr := os.Getenv("ARGO_SECURE"); secureStr != "" {
		if secure, err := strconv.ParseBool(secureStr); err == nil {
			config.Secure = secure
		}
	}

	// Parse ARGO_INSECURE_SKIP_VERIFY if set
	if skipVerifyStr := os.Getenv("ARGO_INSECURE_SKIP_VERIFY"); skipVerifyStr != "" {
		if skipVerify, err := strconv.ParseBool(skipVerifyStr); err == nil {
			config.InsecureSkipVerify = skipVerify
		}
	}

	// Default namespace if not set
	if config.Namespace == "" {
		config.Namespace = "default"
	}

	return config
}
