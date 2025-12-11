// Package server provides the MCP server implementation for Argo Workflows.
package server

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Server wraps the MCP server and provides methods for managing tools and resources.
type Server struct {
	mcp *mcp.Server
}

// NewServer creates and initializes a new MCP server instance.
// It configures the server with the given name and version.
func NewServer(name, version string) *Server {
	implementation := &mcp.Implementation{
		Name:    name,
		Version: version,
	}

	// Create the MCP server with basic options
	// Tools capability is enabled by default when tools are added
	mcpServer := mcp.NewServer(implementation, nil)

	return &Server{
		mcp: mcpServer,
	}
}

// GetMCPServer returns the underlying MCP server instance.
// This is useful for transport setup and starting the server.
func (s *Server) GetMCPServer() *mcp.Server {
	return s.mcp
}
