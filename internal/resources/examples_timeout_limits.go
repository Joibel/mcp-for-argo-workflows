// Package resources implements MCP resources for Argo Workflows examples.
package resources

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	examplesTimeoutLimitsURI  = "argo://examples/timeout-limits"
	examplesTimeoutLimitsName = "examples-timeout-limits"
)

// ExamplesTimeoutLimitsResource returns the MCP resource definition for timeout-limits example.
func ExamplesTimeoutLimitsResource() *mcp.Resource {
	return &mcp.Resource{
		URI:         examplesTimeoutLimitsURI,
		Name:        examplesTimeoutLimitsName,
		Title:       "Timeout and Limits Example",
		Description: "activeDeadlineSeconds and template-level timeout configurations",
		MIMEType:    "text/markdown",
	}
}

// ExamplesTimeoutLimitsHandler returns a handler function for the timeout-limits example resource.
func ExamplesTimeoutLimitsHandler() mcp.ResourceHandler {
	return func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		if req.Params.URI != examplesTimeoutLimitsURI {
			return nil, mcp.ResourceNotFoundError(req.Params.URI)
		}

		content, err := readDoc("examples_timeout_limits.md")
		if err != nil {
			return nil, err
		}

		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:      examplesTimeoutLimitsURI,
					MIMEType: "text/markdown",
					Text:     content,
				},
			},
		}, nil
	}
}
