// Package resources implements MCP resource handlers for Argo Workflows schema documentation.
package resources

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ResourceRegistrar is a function that registers a resource with the MCP server.
type ResourceRegistrar func(s *mcp.Server)

// AllResources returns all resource registrars in the order they should be registered.
func AllResources() []ResourceRegistrar {
	return []ResourceRegistrar{
		RegisterWorkflowSchema,
		RegisterWorkflowTemplateSchema,
		RegisterClusterWorkflowTemplateSchema,
		RegisterCronWorkflowSchema,
	}
}

// RegisterAll registers all resources with the MCP server.
func RegisterAll(s *mcp.Server) {
	for _, register := range AllResources() {
		register(s)
	}
}

// RegisterWorkflowSchema registers the Workflow schema resource.
func RegisterWorkflowSchema(s *mcp.Server) {
	s.AddResource(WorkflowSchemaResource(), WorkflowSchemaHandler())
}

// RegisterWorkflowTemplateSchema registers the WorkflowTemplate schema resource.
func RegisterWorkflowTemplateSchema(s *mcp.Server) {
	s.AddResource(WorkflowTemplateSchemaResource(), WorkflowTemplateSchemaHandler())
}

// RegisterClusterWorkflowTemplateSchema registers the ClusterWorkflowTemplate schema resource.
func RegisterClusterWorkflowTemplateSchema(s *mcp.Server) {
	s.AddResource(ClusterWorkflowTemplateSchemaResource(), ClusterWorkflowTemplateSchemaHandler())
}

// RegisterCronWorkflowSchema registers the CronWorkflow schema resource.
func RegisterCronWorkflowSchema(s *mcp.Server) {
	s.AddResource(CronWorkflowSchemaResource(), CronWorkflowSchemaHandler())
}
