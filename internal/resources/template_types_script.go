// Package resources implements MCP resources for Argo Workflows template type documentation.
package resources

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// TemplateTypesScriptResource returns the MCP resource definition for script template documentation.
func TemplateTypesScriptResource() *mcp.Resource {
	return &mcp.Resource{
		URI:         "argo://docs/template-types/script",
		Name:        "template-types-script",
		Title:       "Script Template Type",
		Description: "Documentation for the Script Template Type",
		MIMEType:    "text/markdown",
	}
}

// TemplateTypesScriptHandler returns a handler function for the script template type resource.
func TemplateTypesScriptHandler() mcp.ResourceHandler {
	return func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		if req.Params.URI != "argo://docs/template-types/script" {
			return nil, mcp.ResourceNotFoundError(req.Params.URI)
		}

		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:      "argo://docs/template-types/script",
					MIMEType: "text/markdown",
					Text:     scriptMarkdown,
				},
			},
		}, nil
	}
}

const scriptMarkdown = `# Script Template Type

Execute inline scripts without needing custom container images.

## Key Fields

- **image** (required) - Container image with interpreter
- **source** (required) - Inline script code
- **command** - Script interpreter (e.g., python, bash)

## See Full Documentation

This is a summary. For complete examples and best practices, refer to the Argo Workflows documentation.
`
