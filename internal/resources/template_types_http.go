// Package resources implements MCP resources for Argo Workflows template type documentation.
package resources

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// TemplateTypesHTTPResource returns the MCP resource definition for http template documentation.
func TemplateTypesHTTPResource() *mcp.Resource {
	return &mcp.Resource{
		URI:         "argo://docs/template-types/http",
		Name:        "template-types-http",
		Title:       "HTTP Template Type",
		Description: "Documentation for the HTTP Template Type",
		MIMEType:    "text/markdown",
	}
}

// TemplateTypesHTTPHandler returns a handler function for the http template type resource.
func TemplateTypesHTTPHandler() mcp.ResourceHandler {
	return func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		if req.Params.URI != "argo://docs/template-types/http" {
			return nil, mcp.ResourceNotFoundError(req.Params.URI)
		}

		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:      "argo://docs/template-types/http",
					MIMEType: "text/markdown",
					Text:     httpMarkdown,
				},
			},
		}, nil
	}
}

const httpMarkdown = `# HTTP Template Type

Make HTTP requests as workflow steps.

## Key Fields

- **url** (required) - URL to request
- **method** - HTTP method (default: GET)
- **headers** - HTTP headers
- **body** - Request body
- **timeoutSeconds** - Request timeout

## See Full Documentation

This is a summary. For complete examples and best practices, refer to the Argo Workflows documentation.
`
