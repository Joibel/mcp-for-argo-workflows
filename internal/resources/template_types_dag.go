// Package resources implements MCP resources for Argo Workflows template type documentation.
package resources

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// TemplateTypesDAGResource returns the MCP resource definition for dag template documentation.
func TemplateTypesDAGResource() *mcp.Resource {
	return &mcp.Resource{
		URI:         "argo://docs/template-types/dag",
		Name:        "template-types-dag",
		Title:       "DAG Template Type",
		Description: "Documentation for the DAG Template Type",
		MIMEType:    "text/markdown",
	}
}

// TemplateTypesDAGHandler returns a handler function for the dag template type resource.
func TemplateTypesDAGHandler() mcp.ResourceHandler {
	return func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		if req.Params.URI != "argo://docs/template-types/dag" {
			return nil, mcp.ResourceNotFoundError(req.Params.URI)
		}

		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:      "argo://docs/template-types/dag",
					MIMEType: "text/markdown",
					Text:     dagMarkdown,
				},
			},
		}, nil
	}
}

const dagMarkdown = `# DAG Template Type

Define tasks with explicit dependencies for maximum parallelism.

## Key Fields

- **tasks** (required) - List of tasks to execute
- Each task has:
  - **name** - Task name
  - **template** - Template to execute
  - **dependencies** - Other tasks to wait for

## See Full Documentation

This is a summary. For complete examples and best practices, refer to the Argo Workflows documentation.
`
