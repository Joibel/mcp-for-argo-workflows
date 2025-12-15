// Package resources implements MCP resources for Argo Workflows examples.
package resources

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	examplesMultiStepURI  = "argo://examples/multi-step"
	examplesMultiStepName = "examples-multi-step"
)

// ExamplesMultiStepResource returns the MCP resource definition for multi-step example.
func ExamplesMultiStepResource() *mcp.Resource {
	return &mcp.Resource{
		URI:         examplesMultiStepURI,
		Name:        examplesMultiStepName,
		Title:       "Multi-Step Workflow Example",
		Description: "Sequential steps with data passing between steps",
		MIMEType:    "text/markdown",
	}
}

// ExamplesMultiStepHandler returns a handler function for the multi-step example resource.
func ExamplesMultiStepHandler() mcp.ResourceHandler {
	return func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		if req.Params.URI != examplesMultiStepURI {
			return nil, mcp.ResourceNotFoundError(req.Params.URI)
		}

		content, err := readDoc("examples_multi_step.md")
		if err != nil {
			return nil, err
		}

		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:      examplesMultiStepURI,
					MIMEType: "text/markdown",
					Text:     content,
				},
			},
		}, nil
	}
}
