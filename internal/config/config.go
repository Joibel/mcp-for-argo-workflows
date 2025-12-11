// Package config handles configuration parsing and validation.
package config

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/spf13/pflag"

	"github.com/Joibel/mcp-for-argo-workflows/internal/argo"
)

// Config holds the combined configuration for the MCP server.
type Config struct {
	// Server settings
	Transport string // "stdio" or "http"
	HTTPAddr  string // HTTP listen address (e.g., ":8080")

	// Argo connection settings
	ArgoServer string // Argo Server host:port (empty = direct K8s)
	ArgoToken  string // Bearer token for Argo Server auth
	Namespace  string // Default namespace for operations

	// Kubernetes settings (when not using Argo Server)
	Kubeconfig string // Path to kubeconfig file
	Context    string // Kubernetes context to use

	// TLS settings (grouped together for alignment)
	Secure             bool // Use TLS when connecting to Argo Server
	InsecureSkipVerify bool // Skip TLS certificate verification
}

// DefaultConfig returns a Config with default values.
func DefaultConfig() *Config {
	return &Config{
		Transport: "stdio",
		HTTPAddr:  ":8080",
		Namespace: "default",
		Secure:    true,
	}
}

// NewFromFlags creates a Config from CLI flags and environment variables.
// Precedence: CLI flags > Environment variables > Default values.
func NewFromFlags() (*Config, error) {
	cfg := DefaultConfig()

	// Define CLI flags
	pflag.StringVar(&cfg.Transport, "transport", cfg.Transport, "MCP transport mode: stdio or http")
	pflag.StringVar(&cfg.HTTPAddr, "http-addr", cfg.HTTPAddr, "HTTP listen address")
	pflag.StringVar(&cfg.ArgoServer, "argo-server", cfg.ArgoServer, "Argo Server host:port (empty = direct K8s)")
	pflag.StringVar(&cfg.ArgoToken, "argo-token", cfg.ArgoToken, "Bearer token for Argo Server auth")
	pflag.StringVar(&cfg.Namespace, "namespace", cfg.Namespace, "Default namespace for operations")
	pflag.BoolVar(&cfg.Secure, "argo-secure", cfg.Secure, "Use TLS when connecting to Argo Server")
	pflag.BoolVar(&cfg.InsecureSkipVerify, "argo-insecure-skip-verify", cfg.InsecureSkipVerify, "Skip TLS certificate verification")
	pflag.StringVar(&cfg.Kubeconfig, "kubeconfig", cfg.Kubeconfig, "Path to kubeconfig file")
	pflag.StringVar(&cfg.Context, "context", cfg.Context, "Kubernetes context to use")

	// Parse CLI flags
	pflag.Parse()

	// Apply environment variables for values not set via CLI flags
	applyEnvOverrides(cfg)

	return cfg, nil
}

// applyEnvOverrides applies environment variable values for unset flags.
func applyEnvOverrides(cfg *Config) {
	// Only override if the flag was not explicitly set
	if !pflag.CommandLine.Changed("transport") {
		if v := os.Getenv("MCP_TRANSPORT"); v != "" {
			cfg.Transport = v
		}
	}

	if !pflag.CommandLine.Changed("http-addr") {
		if v := os.Getenv("MCP_HTTP_ADDR"); v != "" {
			cfg.HTTPAddr = v
		}
	}

	if !pflag.CommandLine.Changed("argo-server") {
		if v := os.Getenv("ARGO_SERVER"); v != "" {
			cfg.ArgoServer = v
		}
	}

	if !pflag.CommandLine.Changed("argo-token") {
		if v := os.Getenv("ARGO_TOKEN"); v != "" {
			cfg.ArgoToken = v
		}
	}

	if !pflag.CommandLine.Changed("namespace") {
		if v := os.Getenv("ARGO_NAMESPACE"); v != "" {
			cfg.Namespace = v
		}
	}

	if !pflag.CommandLine.Changed("argo-secure") {
		if v := os.Getenv("ARGO_SECURE"); v != "" {
			b, err := strconv.ParseBool(v)
			if err != nil {
				slog.Warn("invalid ARGO_SECURE value, using default",
					"value", v, "default", cfg.Secure)
			} else {
				cfg.Secure = b
			}
		}
	}

	if !pflag.CommandLine.Changed("argo-insecure-skip-verify") {
		if v := os.Getenv("ARGO_INSECURE_SKIP_VERIFY"); v != "" {
			b, err := strconv.ParseBool(v)
			if err != nil {
				slog.Warn("invalid ARGO_INSECURE_SKIP_VERIFY value, using default",
					"value", v, "default", cfg.InsecureSkipVerify)
			} else {
				cfg.InsecureSkipVerify = b
			}
		}
	}

	if !pflag.CommandLine.Changed("kubeconfig") {
		if v := os.Getenv("KUBECONFIG"); v != "" {
			cfg.Kubeconfig = v
		}
	}

	// Note: There's no standard env var for Kubernetes context,
	// so --context is CLI-only
}

// ToArgoConfig converts the Config to an argo.Config for creating the Argo client.
func (c *Config) ToArgoConfig() *argo.Config {
	return &argo.Config{
		ArgoServer:         c.ArgoServer,
		ArgoToken:          c.ArgoToken,
		Namespace:          c.Namespace,
		Kubeconfig:         c.Kubeconfig,
		Secure:             c.Secure,
		InsecureSkipVerify: c.InsecureSkipVerify,
	}
}

// IsHTTPTransport returns true if the HTTP transport mode is configured.
func (c *Config) IsHTTPTransport() bool {
	return c.Transport == "http"
}
