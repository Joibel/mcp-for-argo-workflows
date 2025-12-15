// Package resources implements MCP resources for Argo Workflows examples.
package resources

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	examplesLoopsURI  = "argo://examples/loops"
	examplesLoopsName = "examples-loops"
)

// ExamplesLoopsResource returns the MCP resource definition for loops example.
func ExamplesLoopsResource() *mcp.Resource {
	return &mcp.Resource{
		URI:         examplesLoopsURI,
		Name:        examplesLoopsName,
		Title:       "Loops Example",
		Description: "withItems, withParam, and withSequence for iteration patterns",
		MIMEType:    "text/markdown",
	}
}

// ExamplesLoopsHandler returns a handler function for the loops example resource.
func ExamplesLoopsHandler() mcp.ResourceHandler {
	return func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		if req.Params.URI != examplesLoopsURI {
			return nil, mcp.ResourceNotFoundError(req.Params.URI)
		}

		content, err := readDoc("examples_loops.md")
		if err != nil {
			return nil, err
		}

		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:      examplesLoopsURI,
					MIMEType: "text/markdown",
					Text:     content,
				},
			},
		}, nil
	}
}
