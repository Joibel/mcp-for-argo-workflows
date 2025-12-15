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

		content, err := readDoc("template_types_dag.md")
		if err != nil {
			return nil, err
		}

		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:      dagTemplateURI,
					MIMEType: "text/markdown",
					Text:     content,
				},
			},
		}, nil
	}
}
