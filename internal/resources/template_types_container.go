// Package resources implements MCP resources for Argo Workflows template type documentation.
package resources

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	containerTemplateURI  = "argo://docs/template-types/container"
	containerTemplateName = "template-types-container"
)

// TemplateTypesContainerResource returns the MCP resource definition for container template documentation.
func TemplateTypesContainerResource() *mcp.Resource {
	return &mcp.Resource{
		URI:         containerTemplateURI,
		Name:        containerTemplateName,
		Title:       "Container Template Type",
		Description: "Documentation for the Container Template Type",
		MIMEType:    "text/markdown",
	}
}

// TemplateTypesContainerHandler returns a handler function for the container template type resource.
func TemplateTypesContainerHandler() mcp.ResourceHandler {
	return func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		if req.Params.URI != containerTemplateURI {
			return nil, mcp.ResourceNotFoundError(req.Params.URI)
		}

		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:      containerTemplateURI,
					MIMEType: "text/markdown",
					Text:     containerMarkdown,
				},
			},
		}, nil
	}
}

const containerMarkdown = `# Container Template Type

Run containers with specified images, commands, and arguments.

## Key Fields

- **image** (required) - Container image to run
- **command** - Override container entrypoint
- **args** - Arguments to pass to command
- **env** - Environment variables
- **resources** - CPU/memory requests and limits
- **volumeMounts** - Volumes to mount

## Example

` + "```yaml" + `
templates:
  - name: hello
    container:
      image: alpine:latest
      command: [echo]
      args: ["Hello, World!"]
      resources:
        requests:
          memory: "64Mi"
          cpu: "100m"
        limits:
          memory: "128Mi"
          cpu: "200m"
      env:
        - name: MY_VAR
          value: "my-value"
` + "```" + `

## When to Use

Container templates are ideal when you need to run a specific Docker image with custom commands or arguments.

## See Full Documentation

For complete examples and best practices, refer to the Argo Workflows documentation.
`
