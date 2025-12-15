package resources

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ClusterWorkflowTemplateSchemaResource returns the MCP resource definition for the ClusterWorkflowTemplate CRD schema.
func ClusterWorkflowTemplateSchemaResource() *mcp.Resource {
	return &mcp.Resource{
		URI:         "argo://schemas/cluster-workflow-template",
		Name:        "cluster-workflow-template-schema",
		Title:       "Argo ClusterWorkflowTemplate CRD Schema",
		Description: "Complete schema documentation for the ClusterWorkflowTemplate custom resource definition",
		MIMEType:    "text/markdown",
	}
}

// ClusterWorkflowTemplateSchemaHandler returns a handler function for the ClusterWorkflowTemplate schema resource.
func ClusterWorkflowTemplateSchemaHandler() mcp.ResourceHandler {
	return func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		if req.Params.URI != "argo://schemas/cluster-workflow-template" {
			return nil, mcp.ResourceNotFoundError(req.Params.URI)
		}

		content, err := readDoc("cluster_workflow_template_schema.md")
		if err != nil {
			return nil, err
		}

		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:      "argo://schemas/cluster-workflow-template",
					MIMEType: "text/markdown",
					Text:     content,
				},
			},
		}, nil
	}
}
