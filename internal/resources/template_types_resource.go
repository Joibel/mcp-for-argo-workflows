// Package resources implements MCP resources for Argo Workflows template type documentation.
package resources

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// TemplateTypesResourceResource returns the MCP resource definition for resource template documentation.
func TemplateTypesResourceResource() *mcp.Resource {
	return &mcp.Resource{
		URI:         "argo://docs/template-types/resource",
		Name:        "template-types-resource",
		Title:       "Resource Template Type",
		Description: "Documentation for the Resource Template Type",
		MIMEType:    "text/markdown",
	}
}

// TemplateTypesResourceHandler returns a handler function for the resource template type resource.
func TemplateTypesResourceHandler() mcp.ResourceHandler {
	return func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		if req.Params.URI != "argo://docs/template-types/resource" {
			return nil, mcp.ResourceNotFoundError(req.Params.URI)
		}

		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:      "argo://docs/template-types/resource",
					MIMEType: "text/markdown",
					Text:     resourceMarkdown,
				},
			},
		}, nil
	}
}

const resourceMarkdown = `# Resource Template Type

Create, apply, patch, or delete Kubernetes resources.

## Key Fields

- **action** (required) - create, apply, patch, delete, or get
- **manifest** (required) - Resource YAML/JSON
- **successCondition** - Condition for success
- **failureCondition** - Condition for failure

## See Full Documentation

This is a summary. For complete examples and best practices, refer to the Argo Workflows documentation.
`
