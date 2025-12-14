// Package resources implements MCP resources for Argo Workflows template type documentation.
package resources

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	overviewTemplateURI  = "argo://docs/template-types"
	overviewTemplateName = "template-types-overview"
)

// TemplateTypesOverviewResource returns the MCP resource definition for overview template documentation.
func TemplateTypesOverviewResource() *mcp.Resource {
	return &mcp.Resource{
		URI:         overviewTemplateURI,
		Name:        overviewTemplateName,
		Title:       "Argo Workflows Template Types Overview",
		Description: "Documentation for the Argo Workflows Template Types Overview",
		MIMEType:    "text/markdown",
	}
}

// TemplateTypesOverviewHandler returns a handler function for the overview template type resource.
func TemplateTypesOverviewHandler() mcp.ResourceHandler {
	return func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		if req.Params.URI != overviewTemplateURI {
			return nil, mcp.ResourceNotFoundError(req.Params.URI)
		}

		content, err := readDoc("template_types_overview.md")
		if err != nil {
			return nil, err
		}

		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:      overviewTemplateURI,
					MIMEType: "text/markdown",
					Text:     content,
				},
			},
		}, nil
	}
}
