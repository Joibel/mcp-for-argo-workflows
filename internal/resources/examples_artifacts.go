// Package resources implements MCP resources for Argo Workflows examples.
package resources

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	examplesArtifactsURI  = "argo://examples/artifacts"
	examplesArtifactsName = "examples-artifacts"
)

// ExamplesArtifactsResource returns the MCP resource definition for artifacts example.
func ExamplesArtifactsResource() *mcp.Resource {
	return &mcp.Resource{
		URI:         examplesArtifactsURI,
		Name:        examplesArtifactsName,
		Title:       "Artifacts Example",
		Description: "Artifact passing between steps with S3/GCS configuration",
		MIMEType:    "text/markdown",
	}
}

// ExamplesArtifactsHandler returns a handler function for the artifacts example resource.
func ExamplesArtifactsHandler() mcp.ResourceHandler {
	return func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		if req.Params.URI != examplesArtifactsURI {
			return nil, mcp.ResourceNotFoundError(req.Params.URI)
		}

		content, err := readDoc("examples_artifacts.md")
		if err != nil {
			return nil, err
		}

		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:      examplesArtifactsURI,
					MIMEType: "text/markdown",
					Text:     content,
				},
			},
		}, nil
	}
}
