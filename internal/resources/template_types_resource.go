// Package resources implements MCP resources for Argo Workflows template type documentation.
package resources

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	resourceTemplateURI  = "argo://docs/template-types/resource"
	resourceTemplateName = "template-types-resource"
)

// TemplateTypesResourceResource returns the MCP resource definition for resource template documentation.
func TemplateTypesResourceResource() *mcp.Resource {
	return &mcp.Resource{
		URI:         resourceTemplateURI,
		Name:        resourceTemplateName,
		Title:       "Resource Template Type",
		Description: "Documentation for the Resource Template Type",
		MIMEType:    "text/markdown",
	}
}

// TemplateTypesResourceHandler returns a handler function for the resource template type resource.
func TemplateTypesResourceHandler() mcp.ResourceHandler {
	return func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		if req.Params.URI != resourceTemplateURI {
			return nil, mcp.ResourceNotFoundError(req.Params.URI)
		}

		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:      resourceTemplateURI,
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

## Example

` + "```yaml" + `
templates:
  - name: create-configmap
    resource:
      action: create
      manifest: |
        apiVersion: v1
        kind: ConfigMap
        metadata:
          name: my-config
        data:
          key: value
` + "```" + `

## When to Use

Resource templates are ideal for managing Kubernetes resources as part of your workflow.

## See Full Documentation

For complete examples and best practices, refer to the Argo Workflows documentation.
`
