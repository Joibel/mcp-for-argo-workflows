// Package resources implements MCP resources for Argo Workflows examples.
package resources

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	examplesDAGDiamondURI  = "argo://examples/dag-diamond"
	examplesDAGDiamondName = "examples-dag-diamond"
)

// ExamplesDAGDiamondResource returns the MCP resource definition for dag-diamond example.
func ExamplesDAGDiamondResource() *mcp.Resource {
	return &mcp.Resource{
		URI:         examplesDAGDiamondURI,
		Name:        examplesDAGDiamondName,
		Title:       "DAG Diamond Pattern Example",
		Description: "Classic diamond DAG with fan-out and fan-in pattern",
		MIMEType:    "text/markdown",
	}
}

// ExamplesDAGDiamondHandler returns a handler function for the dag-diamond example resource.
func ExamplesDAGDiamondHandler() mcp.ResourceHandler {
	return func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		if req.Params.URI != examplesDAGDiamondURI {
			return nil, mcp.ResourceNotFoundError(req.Params.URI)
		}

		content, err := readDoc("examples_dag_diamond.md")
		if err != nil {
			return nil, err
		}

		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:      examplesDAGDiamondURI,
					MIMEType: "text/markdown",
					Text:     content,
				},
			},
		}, nil
	}
}
