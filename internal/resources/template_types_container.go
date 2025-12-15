// Package resources implements MCP resources for Argo Workflows template type documentation.
package resources

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	containerTemplateURI  = "argo://docs/template-types/container"
	containerTemplateName = "template-types-container"
)

// TemplateTypesContainerResource returns the MCP resource definition for container template documentation.
func TemplateTypesContainerResource() *mcp.Resource {
	return &mcp.Resource{
		URI:         containerTemplateURI,
		Name:        containerTemplateName,
		Title:       "Container Template Type",
		Description: "Documentation for the Container Template Type",
		MIMEType:    "text/markdown",
	}
}

// TemplateTypesContainerHandler returns a handler function for the container template type resource.
func TemplateTypesContainerHandler() mcp.ResourceHandler {
	return func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		if req.Params.URI != containerTemplateURI {
			return nil, mcp.ResourceNotFoundError(req.Params.URI)
		}

		content, err := readDoc("template_types_container.md")
		if err != nil {
			return nil, err
		}

		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:      containerTemplateURI,
					MIMEType: "text/markdown",
					Text:     content,
				},
			},
		}, nil
	}
}
