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

		content, err := readDoc("template_types_script.md")
		if err != nil {
			return nil, err
		}

		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:      scriptTemplateURI,
					MIMEType: "text/markdown",
					Text:     content,
				},
			},
		}, nil
	}
}
