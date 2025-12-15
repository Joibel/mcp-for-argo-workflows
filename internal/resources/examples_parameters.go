// Package resources implements MCP resources for Argo Workflows examples.
package resources

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	examplesParametersURI  = "argo://examples/parameters"
	examplesParametersName = "examples-parameters"
)

// ExamplesParametersResource returns the MCP resource definition for parameters example.
func ExamplesParametersResource() *mcp.Resource {
	return &mcp.Resource{
		URI:         examplesParametersURI,
		Name:        examplesParametersName,
		Title:       "Parameters Example",
		Description: "Input parameters, default values, and parameter passing patterns",
		MIMEType:    "text/markdown",
	}
}

// ExamplesParametersHandler returns a handler function for the parameters example resource.
func ExamplesParametersHandler() mcp.ResourceHandler {
	return func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		if req.Params.URI != examplesParametersURI {
			return nil, mcp.ResourceNotFoundError(req.Params.URI)
		}

		content, err := readDoc("examples_parameters.md")
		if err != nil {
			return nil, err
		}

		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:      examplesParametersURI,
					MIMEType: "text/markdown",
					Text:     content,
				},
			},
		}, nil
	}
}
