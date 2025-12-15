package resources

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// CronWorkflowSchemaResource returns the MCP resource definition for the CronWorkflow CRD schema.
func CronWorkflowSchemaResource() *mcp.Resource {
	return &mcp.Resource{
		URI:         "argo://schemas/cron-workflow",
		Name:        "cron-workflow-schema",
		Title:       "Argo CronWorkflow CRD Schema",
		Description: "Complete schema documentation for the CronWorkflow custom resource definition",
		MIMEType:    "text/markdown",
	}
}

// CronWorkflowSchemaHandler returns a handler function for the CronWorkflow schema resource.
func CronWorkflowSchemaHandler() mcp.ResourceHandler {
	return func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		if req.Params.URI != "argo://schemas/cron-workflow" {
			return nil, mcp.ResourceNotFoundError(req.Params.URI)
		}

		content, err := readDoc("cron_workflow_schema.md")
		if err != nil {
			return nil, err
		}

		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:      "argo://schemas/cron-workflow",
					MIMEType: "text/markdown",
					Text:     content,
				},
			},
		}, nil
	}
}
