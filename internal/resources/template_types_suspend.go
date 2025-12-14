// Package resources implements MCP resources for Argo Workflows template type documentation.
package resources

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// TemplateTypesSuspendResource returns the MCP resource definition for suspend template documentation.
func TemplateTypesSuspendResource() *mcp.Resource {
	return &mcp.Resource{
		URI:         "argo://docs/template-types/suspend",
		Name:        "template-types-suspend",
		Title:       "Suspend Template Type",
		Description: "Documentation for the Suspend Template Type",
		MIMEType:    "text/markdown",
	}
}

// TemplateTypesSuspendHandler returns a handler function for the suspend template type resource.
func TemplateTypesSuspendHandler() mcp.ResourceHandler {
	return func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		if req.Params.URI != "argo://docs/template-types/suspend" {
			return nil, mcp.ResourceNotFoundError(req.Params.URI)
		}

		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:      "argo://docs/template-types/suspend",
					MIMEType: "text/markdown",
					Text:     suspendMarkdown,
				},
			},
		}, nil
	}
}

const suspendMarkdown = `# Suspend Template Type

Pause workflow execution for approval gates or timed delays.

## Key Fields

- **duration** - How long to suspend (e.g., "10s", "5m", "1h")
  - If omitted, suspends indefinitely until manually resumed

## See Full Documentation

This is a summary. For complete examples and best practices, refer to the Argo Workflows documentation.
`
