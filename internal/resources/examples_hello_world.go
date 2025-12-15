// Package resources implements MCP resources for Argo Workflows examples.
package resources

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	examplesHelloWorldURI  = "argo://examples/hello-world"
	examplesHelloWorldName = "examples-hello-world"
)

// ExamplesHelloWorldResource returns the MCP resource definition for hello-world example.
func ExamplesHelloWorldResource() *mcp.Resource {
	return &mcp.Resource{
		URI:         examplesHelloWorldURI,
		Name:        examplesHelloWorldName,
		Title:       "Hello World Workflow Example",
		Description: "Simplest workflow example with a single container template",
		MIMEType:    "text/markdown",
	}
}

// ExamplesHelloWorldHandler returns a handler function for the hello-world example resource.
func ExamplesHelloWorldHandler() mcp.ResourceHandler {
	return func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		if req.Params.URI != examplesHelloWorldURI {
			return nil, mcp.ResourceNotFoundError(req.Params.URI)
		}

		content, err := readDoc("examples_hello_world.md")
		if err != nil {
			return nil, err
		}

		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:      examplesHelloWorldURI,
					MIMEType: "text/markdown",
					Text:     content,
				},
			},
		}, nil
	}
}
