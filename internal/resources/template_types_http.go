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

		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:      httpTemplateURI,
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

## Example

` + "```yaml" + `
templates:
  - name: http-request
    http:
      url: "https://api.example.com/webhook"
      method: "POST"
      headers:
        - name: "Content-Type"
          value: "application/json"
      body: '{"message": "Hello from Argo"}'
      timeoutSeconds: 30
` + "```" + `

## When to Use

HTTP templates are ideal for calling REST APIs, webhooks, or external services.

## See Full Documentation

For complete examples and best practices, refer to the Argo Workflows documentation.
`
