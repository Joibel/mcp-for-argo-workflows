// Package resources implements MCP resources for Argo Workflows examples.
package resources

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	examplesConditionalsURI  = "argo://examples/conditionals"
	examplesConditionalsName = "examples-conditionals"
)

// ExamplesConditionalsResource returns the MCP resource definition for conditionals example.
func ExamplesConditionalsResource() *mcp.Resource {
	return &mcp.Resource{
		URI:         examplesConditionalsURI,
		Name:        examplesConditionalsName,
		Title:       "Conditionals Example",
		Description: "Conditional step execution using when expressions",
		MIMEType:    "text/markdown",
	}
}

// ExamplesConditionalsHandler returns a handler function for the conditionals example resource.
func ExamplesConditionalsHandler() mcp.ResourceHandler {
	return func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		if req.Params.URI != examplesConditionalsURI {
			return nil, mcp.ResourceNotFoundError(req.Params.URI)
		}

		content, err := readDoc("examples_conditionals.md")
		if err != nil {
			return nil, err
		}

		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:      examplesConditionalsURI,
					MIMEType: "text/markdown",
					Text:     content,
				},
			},
		}, nil
	}
}
