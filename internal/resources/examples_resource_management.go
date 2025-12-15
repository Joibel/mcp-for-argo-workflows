// Package resources implements MCP resources for Argo Workflows examples.
package resources

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	examplesResourceManagementURI  = "argo://examples/resource-management"
	examplesResourceManagementName = "examples-resource-management"
)

// ExamplesResourceManagementResource returns the MCP resource definition for resource-management example.
func ExamplesResourceManagementResource() *mcp.Resource {
	return &mcp.Resource{
		URI:         examplesResourceManagementURI,
		Name:        examplesResourceManagementName,
		Title:       "Resource Management Example",
		Description: "CPU/memory requests and limits, pod priority, and resource optimization",
		MIMEType:    "text/markdown",
	}
}

// ExamplesResourceManagementHandler returns a handler function for the resource-management example resource.
func ExamplesResourceManagementHandler() mcp.ResourceHandler {
	return func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		if req.Params.URI != examplesResourceManagementURI {
			return nil, mcp.ResourceNotFoundError(req.Params.URI)
		}

		content, err := readDoc("examples_resource_management.md")
		if err != nil {
			return nil, err
		}

		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:      examplesResourceManagementURI,
					MIMEType: "text/markdown",
					Text:     content,
				},
			},
		}, nil
	}
}
