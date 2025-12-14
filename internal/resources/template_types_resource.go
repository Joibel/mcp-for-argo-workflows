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

		content, err := readDoc("template_types_resource.md")
		if err != nil {
			return nil, err
		}

		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:      resourceTemplateURI,
					MIMEType: "text/markdown",
					Text:     content,
				},
			},
		}, nil
	}
}
