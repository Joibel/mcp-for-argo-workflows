// Package resources implements MCP resources for Argo Workflows schema documentation.
package resources

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// WorkflowSchemaResource returns the MCP resource definition for the Workflow CRD schema.
func WorkflowSchemaResource() *mcp.Resource {
	return &mcp.Resource{
		URI:         "argo://schemas/workflow",
		Name:        "workflow-schema",
		Title:       "Argo Workflow CRD Schema",
		Description: "Complete schema documentation for the Workflow custom resource definition",
		MIMEType:    "text/markdown",
	}
}

// WorkflowSchemaHandler returns a handler function for the Workflow schema resource.
func WorkflowSchemaHandler() mcp.ResourceHandler {
	return func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		if req.Params.URI != "argo://schemas/workflow" {
			return nil, mcp.ResourceNotFoundError(req.Params.URI)
		}

		content, err := readDoc("workflow_schema.md")
		if err != nil {
			return nil, err
		}

		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:      "argo://schemas/workflow",
					MIMEType: "text/markdown",
					Text:     content,
				},
			},
		}, nil
	}
}
