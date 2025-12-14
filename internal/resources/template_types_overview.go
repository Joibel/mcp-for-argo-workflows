// Package resources implements MCP resources for Argo Workflows template type documentation.
package resources

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	overviewTemplateURI  = "argo://docs/template-types"
	overviewTemplateName = "template-types-overview"
)

// TemplateTypesOverviewResource returns the MCP resource definition for overview template documentation.
func TemplateTypesOverviewResource() *mcp.Resource {
	return &mcp.Resource{
		URI:         overviewTemplateURI,
		Name:        overviewTemplateName,
		Title:       "Argo Workflows Template Types Overview",
		Description: "Documentation for the Argo Workflows Template Types Overview",
		MIMEType:    "text/markdown",
	}
}

// TemplateTypesOverviewHandler returns a handler function for the overview template type resource.
func TemplateTypesOverviewHandler() mcp.ResourceHandler {
	return func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		if req.Params.URI != overviewTemplateURI {
			return nil, mcp.ResourceNotFoundError(req.Params.URI)
		}

		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:      overviewTemplateURI,
					MIMEType: "text/markdown",
					Text:     overviewMarkdown,
				},
			},
		}, nil
	}
}

const overviewMarkdown = `# Argo Workflows Template Types Overview

Templates are the building blocks of Argo Workflows. Each template defines a unit of work and must specify exactly ONE template type.

## Available Resources

For detailed documentation on each template type, use these resources:

- **argo://docs/template-types/container** - Container template documentation
- **argo://docs/template-types/script** - Script template documentation
- **argo://docs/template-types/dag** - DAG template documentation
- **argo://docs/template-types/steps** - Steps template documentation
- **argo://docs/template-types/suspend** - Suspend template documentation
- **argo://docs/template-types/resource** - Resource template documentation
- **argo://docs/template-types/http** - HTTP template documentation

## Template Types Quick Reference

Argo Workflows supports several template types:

### Execution Templates
- **Container** - Run containers with specific images and commands
- **Script** - Execute inline code/scripts
- **Resource** - Manage Kubernetes resources
- **HTTP** - Make HTTP requests
- **Suspend** - Pause workflow execution

### Orchestration Templates
- **Steps** - Sequential and parallel step groups
- **DAG** - Tasks with explicit dependencies

## Choosing the Right Template Type

- Use **Container** for running existing Docker images
- Use **Script** for inline code that doesn't need a custom image
- Use **Resource** for creating/managing Kubernetes resources
- Use **HTTP** for calling REST APIs
- Use **Suspend** for approval gates or delays
- Use **Steps** for simple sequential/parallel workflows
- Use **DAG** for complex dependencies and maximum parallelism
`
