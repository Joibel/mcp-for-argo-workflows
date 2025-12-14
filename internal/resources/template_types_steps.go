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

		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:      stepsTemplateURI,
					MIMEType: "text/markdown",
					Text:     stepsMarkdown,
				},
			},
		}, nil
	}
}

const stepsMarkdown = `# Steps Template Type

Define sequential execution with support for parallel step groups.

## Structure

Steps are organized into groups:
- Each group is a list of steps
- Steps within a group run in parallel
- Groups run sequentially

## Example

` + "```yaml" + `
templates:
  - name: hello-hello-hello
    steps:
      - - name: step1
          template: whalesay
          arguments:
            parameters: [{name: message, value: "hello1"}]
      - - name: step2a
          template: whalesay
          arguments:
            parameters: [{name: message, value: "hello2a"}]
        - name: step2b
          template: whalesay
          arguments:
            parameters: [{name: message, value: "hello2b"}]
` + "```" + `

## When to Use

Steps templates are ideal for simple sequential workflows with optional parallel groups.

## See Full Documentation

For complete examples and best practices, refer to the Argo Workflows documentation.
`
