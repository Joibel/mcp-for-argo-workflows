// Package resources implements MCP resources for Argo Workflows template type documentation.
package resources

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// TemplateTypesContainerResource returns the MCP resource definition for container template documentation.
func TemplateTypesContainerResource() *mcp.Resource {
	return &mcp.Resource{
		URI:         "argo://docs/template-types/container",
		Name:        "template-types-container",
		Title:       "Container Template Type",
		Description: "Documentation for the Container Template Type",
		MIMEType:    "text/markdown",
	}
}

// TemplateTypesContainerHandler returns a handler function for the container template type resource.
func TemplateTypesContainerHandler() mcp.ResourceHandler {
	return func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		if req.Params.URI != "argo://docs/template-types/container" {
			return nil, mcp.ResourceNotFoundError(req.Params.URI)
		}

		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:      "argo://docs/template-types/container",
					MIMEType: "text/markdown",
					Text:     containerMarkdown,
				},
			},
		}, nil
	}
}

const containerMarkdown = `# Container Template Type

Run containers with specified images, commands, and arguments.

## Key Fields

- **image** (required) - Container image to run
- **command** - Override container entrypoint  
- **args** - Arguments to pass to command
- **env** - Environment variables
- **resources** - CPU/memory requests and limits
- **volumeMounts** - Volumes to mount

## See Full Documentation

This is a summary. For complete examples and best practices, refer to the Argo Workflows documentation.
`
