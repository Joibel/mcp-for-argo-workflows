// Package resources implements MCP resources for Argo Workflows template type documentation.
package resources

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	dagTemplateURI  = "argo://docs/template-types/dag"
	dagTemplateName = "template-types-dag"
)

// TemplateTypesDAGResource returns the MCP resource definition for dag template documentation.
func TemplateTypesDAGResource() *mcp.Resource {
	return &mcp.Resource{
		URI:         dagTemplateURI,
		Name:        dagTemplateName,
		Title:       "DAG Template Type",
		Description: "Documentation for the DAG Template Type",
		MIMEType:    "text/markdown",
	}
}

// TemplateTypesDAGHandler returns a handler function for the dag template type resource.
func TemplateTypesDAGHandler() mcp.ResourceHandler {
	return func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		if req.Params.URI != dagTemplateURI {
			return nil, mcp.ResourceNotFoundError(req.Params.URI)
		}

		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:      dagTemplateURI,
					MIMEType: "text/markdown",
					Text:     dagMarkdown,
				},
			},
		}, nil
	}
}

const dagMarkdown = `# DAG Template Type

Define tasks with explicit dependencies for maximum parallelism.

## Key Fields

- **tasks** (required) - List of tasks to execute
- Each task has:
  - **name** - Task name
  - **template** - Template to execute
  - **dependencies** - Other tasks to wait for

## Example

` + "```yaml" + `
templates:
  - name: diamond
    dag:
      tasks:
        - name: A
          template: echo
          arguments:
            parameters: [{name: message, value: "A"}]
        - name: B
          dependencies: [A]
          template: echo
          arguments:
            parameters: [{name: message, value: "B"}]
        - name: C
          dependencies: [A]
          template: echo
          arguments:
            parameters: [{name: message, value: "C"}]
        - name: D
          dependencies: [B, C]
          template: echo
          arguments:
            parameters: [{name: message, value: "D"}]
` + "```" + `

## When to Use

DAG templates are ideal when tasks have complex dependencies and you want maximum parallelism.

## See Full Documentation

For complete examples and best practices, refer to the Argo Workflows documentation.
`
