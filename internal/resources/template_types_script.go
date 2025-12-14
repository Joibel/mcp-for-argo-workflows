// Package resources implements MCP resources for Argo Workflows template type documentation.
package resources

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	scriptTemplateURI  = "argo://docs/template-types/script"
	scriptTemplateName = "template-types-script"
)

// TemplateTypesScriptResource returns the MCP resource definition for script template documentation.
func TemplateTypesScriptResource() *mcp.Resource {
	return &mcp.Resource{
		URI:         scriptTemplateURI,
		Name:        scriptTemplateName,
		Title:       "Script Template Type",
		Description: "Documentation for the Script Template Type",
		MIMEType:    "text/markdown",
	}
}

// TemplateTypesScriptHandler returns a handler function for the script template type resource.
func TemplateTypesScriptHandler() mcp.ResourceHandler {
	return func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		if req.Params.URI != scriptTemplateURI {
			return nil, mcp.ResourceNotFoundError(req.Params.URI)
		}

		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:      scriptTemplateURI,
					MIMEType: "text/markdown",
					Text:     scriptMarkdown,
				},
			},
		}, nil
	}
}

const scriptMarkdown = `# Script Template Type

Execute inline scripts without needing custom container images.

## Key Fields

- **image** (required) - Container image with interpreter
- **source** (required) - Inline script code
- **command** - Script interpreter (e.g., python, bash)

## Example

` + "```yaml" + `
templates:
  - name: gen-random
    script:
      image: python:alpine
      command: [python]
      source: |
        import random
        result = random.randint(1, 100)
        print(result)
` + "```" + `

## When to Use

Script templates are ideal when you need to run inline code without building a custom container image.

## See Full Documentation

For complete examples and best practices, refer to the Argo Workflows documentation.
`
