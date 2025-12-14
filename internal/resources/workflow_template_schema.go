package resources

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// WorkflowTemplateSchemaResource returns the MCP resource definition for the WorkflowTemplate CRD schema.
func WorkflowTemplateSchemaResource() *mcp.Resource {
	return &mcp.Resource{
		URI:         "argo://schemas/workflow-template",
		Name:        "workflow-template-schema",
		Title:       "Argo WorkflowTemplate CRD Schema",
		Description: "Complete schema documentation for the WorkflowTemplate custom resource definition",
		MIMEType:    "text/markdown",
	}
}

// WorkflowTemplateSchemaHandler returns a handler function for the WorkflowTemplate schema resource.
func WorkflowTemplateSchemaHandler() mcp.ResourceHandler {
	return func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		if req.Params.URI != "argo://schemas/workflow-template" {
			return nil, mcp.ResourceNotFoundError(req.Params.URI)
		}

		content, err := readDoc("workflow_template_schema.md")
		if err != nil {
			return nil, err
		}

		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:      "argo://schemas/workflow-template",
					MIMEType: "text/markdown",
					Text:     content,
				},
			},
		}, nil
	}
}
