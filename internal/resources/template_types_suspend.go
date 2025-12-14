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

		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:      suspendTemplateURI,
					MIMEType: "text/markdown",
					Text:     suspendMarkdown,
				},
			},
		}, nil
	}
}

const suspendMarkdown = `# Suspend Template Type

Pause workflow execution for approval gates or timed delays.

## Key Fields

- **duration** - How long to suspend (e.g., "10s", "5m", "1h")
  - If omitted, suspends indefinitely until manually resumed

## Example

` + "```yaml" + `
templates:
  - name: approve
    suspend:
      duration: "0"  # Suspend indefinitely until resumed

  - name: delay
    suspend:
      duration: "20s"  # Auto-resume after 20 seconds
` + "```" + `

## When to Use

Suspend templates are ideal for manual approval gates or introducing delays in your workflow.

## See Full Documentation

For complete examples and best practices, refer to the Argo Workflows documentation.
`
