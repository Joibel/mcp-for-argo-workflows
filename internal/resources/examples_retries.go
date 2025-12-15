// Package resources implements MCP resources for Argo Workflows examples.
package resources

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	examplesRetriesURI  = "argo://examples/retries"
	examplesRetriesName = "examples-retries"
)

// ExamplesRetriesResource returns the MCP resource definition for retries example.
func ExamplesRetriesResource() *mcp.Resource {
	return &mcp.Resource{
		URI:         examplesRetriesURI,
		Name:        examplesRetriesName,
		Title:       "Retries Example",
		Description: "Retry strategies and retryPolicy configuration",
		MIMEType:    "text/markdown",
	}
}

// ExamplesRetriesHandler returns a handler function for the retries example resource.
func ExamplesRetriesHandler() mcp.ResourceHandler {
	return func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		if req.Params.URI != examplesRetriesURI {
			return nil, mcp.ResourceNotFoundError(req.Params.URI)
		}

		content, err := readDoc("examples_retries.md")
		if err != nil {
			return nil, err
		}

		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:      examplesRetriesURI,
					MIMEType: "text/markdown",
					Text:     content,
				},
			},
		}, nil
	}
}
