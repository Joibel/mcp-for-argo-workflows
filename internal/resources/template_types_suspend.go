// Package resources implements MCP resources for Argo Workflows template type documentation.
package resources

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	suspendTemplateURI  = "argo://docs/template-types/suspend"
	suspendTemplateName = "template-types-suspend"
)

// TemplateTypesSuspendResource returns the MCP resource definition for suspend template documentation.
func TemplateTypesSuspendResource() *mcp.Resource {
	return &mcp.Resource{
		URI:         suspendTemplateURI,
		Name:        suspendTemplateName,
		Title:       "Suspend Template Type",
		Description: "Documentation for the Suspend Template Type",
		MIMEType:    "text/markdown",
	}
}

// TemplateTypesSuspendHandler returns a handler function for the suspend template type resource.
func TemplateTypesSuspendHandler() mcp.ResourceHandler {
	return func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		if req.Params.URI != suspendTemplateURI {
			return nil, mcp.ResourceNotFoundError(req.Params.URI)
		}

		content, err := readDoc("template_types_suspend.md")
		if err != nil {
			return nil, err
		}

		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:      suspendTemplateURI,
					MIMEType: "text/markdown",
					Text:     content,
				},
			},
		}, nil
	}
}
