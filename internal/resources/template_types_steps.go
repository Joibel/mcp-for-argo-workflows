// Package resources implements MCP resources for Argo Workflows template type documentation.
package resources

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	stepsTemplateURI  = "argo://docs/template-types/steps"
	stepsTemplateName = "template-types-steps"
)

// TemplateTypesStepsResource returns the MCP resource definition for steps template documentation.
func TemplateTypesStepsResource() *mcp.Resource {
	return &mcp.Resource{
		URI:         stepsTemplateURI,
		Name:        stepsTemplateName,
		Title:       "Steps Template Type",
		Description: "Documentation for the Steps Template Type",
		MIMEType:    "text/markdown",
	}
}

// TemplateTypesStepsHandler returns a handler function for the steps template type resource.
func TemplateTypesStepsHandler() mcp.ResourceHandler {
	return func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		if req.Params.URI != stepsTemplateURI {
			return nil, mcp.ResourceNotFoundError(req.Params.URI)
		}

		content, err := readDoc("template_types_steps.md")
		if err != nil {
			return nil, err
		}

		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:      stepsTemplateURI,
					MIMEType: "text/markdown",
					Text:     content,
				},
			},
		}, nil
	}
}
