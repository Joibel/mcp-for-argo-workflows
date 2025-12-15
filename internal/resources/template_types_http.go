// Package resources implements MCP resources for Argo Workflows template type documentation.
package resources

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	httpTemplateURI  = "argo://docs/template-types/http"
	httpTemplateName = "template-types-http"
)

// TemplateTypesHTTPResource returns the MCP resource definition for http template documentation.
func TemplateTypesHTTPResource() *mcp.Resource {
	return &mcp.Resource{
		URI:         httpTemplateURI,
		Name:        httpTemplateName,
		Title:       "HTTP Template Type",
		Description: "Documentation for the HTTP Template Type",
		MIMEType:    "text/markdown",
	}
}

// TemplateTypesHTTPHandler returns a handler function for the http template type resource.
func TemplateTypesHTTPHandler() mcp.ResourceHandler {
	return func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		if req.Params.URI != httpTemplateURI {
			return nil, mcp.ResourceNotFoundError(req.Params.URI)
		}

		content, err := readDoc("template_types_http.md")
		if err != nil {
			return nil, err
		}

		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:      httpTemplateURI,
					MIMEType: "text/markdown",
					Text:     content,
				},
			},
		}, nil
	}
}
