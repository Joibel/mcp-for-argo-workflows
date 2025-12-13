// Package server provides the MCP server implementation for Argo Workflows.
package server

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/Joibel/mcp-for-argo-workflows/internal/argo"
	"github.com/Joibel/mcp-for-argo-workflows/internal/tools"
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

// RegisterTools registers all Argo Workflows MCP tools with the server.
func (s *Server) RegisterTools(client *argo.Client) {
	// Register submit_workflow tool
	mcp.AddTool(s.mcp, tools.SubmitWorkflowTool(), tools.SubmitWorkflowHandler(client))

	// Register list_workflows tool
	mcp.AddTool(s.mcp, tools.ListWorkflowsTool(), tools.ListWorkflowsHandler(client))

	// Register get_workflow tool
	mcp.AddTool(s.mcp, tools.GetWorkflowTool(), tools.GetWorkflowHandler(client))

	// Register delete_workflow tool
	mcp.AddTool(s.mcp, tools.DeleteWorkflowTool(), tools.DeleteWorkflowHandler(client))

	// Register watch_workflow tool
	mcp.AddTool(s.mcp, tools.WatchWorkflowTool(), tools.WatchWorkflowHandler(client))

	// Register logs_workflow tool
	mcp.AddTool(s.mcp, tools.LogsWorkflowTool(), tools.LogsWorkflowHandler(client))

	// Register wait_workflow tool
	mcp.AddTool(s.mcp, tools.WaitWorkflowTool(), tools.WaitWorkflowHandler(client))

	// Register lint_workflow tool
	mcp.AddTool(s.mcp, tools.LintWorkflowTool(), tools.LintWorkflowHandler(client))

	// Register retry_workflow tool
	mcp.AddTool(s.mcp, tools.RetryWorkflowTool(), tools.RetryWorkflowHandler(client))

	// Register resubmit_workflow tool
	mcp.AddTool(s.mcp, tools.ResubmitWorkflowTool(), tools.ResubmitWorkflowHandler(client))

	// Register suspend_workflow tool
	mcp.AddTool(s.mcp, tools.SuspendWorkflowTool(), tools.SuspendWorkflowHandler(client))

	// Register resume_workflow tool
	mcp.AddTool(s.mcp, tools.ResumeWorkflowTool(), tools.ResumeWorkflowHandler(client))

	// Register stop_workflow tool
	mcp.AddTool(s.mcp, tools.StopWorkflowTool(), tools.StopWorkflowHandler(client))

	// Register terminate_workflow tool
	mcp.AddTool(s.mcp, tools.TerminateWorkflowTool(), tools.TerminateWorkflowHandler(client))

	// Register list_workflow_templates tool
	mcp.AddTool(s.mcp, tools.ListWorkflowTemplatesTool(), tools.ListWorkflowTemplatesHandler(client))
}

// GetMCPServer returns the underlying MCP server instance.
// This is useful for transport setup and starting the server.
func (s *Server) GetMCPServer() *mcp.Server {
	return s.mcp
}
