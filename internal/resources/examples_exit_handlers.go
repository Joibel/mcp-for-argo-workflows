// Package resources implements MCP resources for Argo Workflows examples.
package resources

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	examplesExitHandlersURI  = "argo://examples/exit-handlers"
	examplesExitHandlersName = "examples-exit-handlers"
)

// ExamplesExitHandlersResource returns the MCP resource definition for exit-handlers example.
func ExamplesExitHandlersResource() *mcp.Resource {
	return &mcp.Resource{
		URI:         examplesExitHandlersURI,
		Name:        examplesExitHandlersName,
		Title:       "Exit Handlers Example",
		Description: "OnExit handlers for cleanup and status-specific actions",
		MIMEType:    "text/markdown",
	}
}

// ExamplesExitHandlersHandler returns a handler function for the exit-handlers example resource.
func ExamplesExitHandlersHandler() mcp.ResourceHandler {
	return func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		if req.Params.URI != examplesExitHandlersURI {
			return nil, mcp.ResourceNotFoundError(req.Params.URI)
		}

		content, err := readDoc("examples_exit_handlers.md")
		if err != nil {
			return nil, err
		}

		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:      examplesExitHandlersURI,
					MIMEType: "text/markdown",
					Text:     content,
				},
			},
		}, nil
	}
}
