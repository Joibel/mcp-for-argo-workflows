// Package resources implements MCP resources for Argo Workflows template type documentation.
package resources

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// TemplateTypesStepsResource returns the MCP resource definition for steps template documentation.
func TemplateTypesStepsResource() *mcp.Resource {
	return &mcp.Resource{
		URI:         "argo://docs/template-types/steps",
		Name:        "template-types-steps",
		Title:       "Steps Template Type",
		Description: "Documentation for the Steps Template Type",
		MIMEType:    "text/markdown",
	}
}

// TemplateTypesStepsHandler returns a handler function for the steps template type resource.
func TemplateTypesStepsHandler() mcp.ResourceHandler {
	return func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		if req.Params.URI != "argo://docs/template-types/steps" {
			return nil, mcp.ResourceNotFoundError(req.Params.URI)
		}

		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:      "argo://docs/template-types/steps",
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

## See Full Documentation

This is a summary. For complete examples and best practices, refer to the Argo Workflows documentation.
`
