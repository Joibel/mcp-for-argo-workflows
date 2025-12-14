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
		RegisterTemplateTypesOverview,
		RegisterTemplateTypesContainer,
		RegisterTemplateTypesScript,
		RegisterTemplateTypesDAG,
		RegisterTemplateTypesSteps,
		RegisterTemplateTypesSuspend,
		RegisterTemplateTypesResource,
		RegisterTemplateTypesHTTP,
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

// RegisterTemplateTypesOverview registers the template types overview resource.
func RegisterTemplateTypesOverview(s *mcp.Server) {
	s.AddResource(TemplateTypesOverviewResource(), TemplateTypesOverviewHandler())
}

// RegisterTemplateTypesContainer registers the container template type resource.
func RegisterTemplateTypesContainer(s *mcp.Server) {
	s.AddResource(TemplateTypesContainerResource(), TemplateTypesContainerHandler())
}

// RegisterTemplateTypesScript registers the script template type resource.
func RegisterTemplateTypesScript(s *mcp.Server) {
	s.AddResource(TemplateTypesScriptResource(), TemplateTypesScriptHandler())
}

// RegisterTemplateTypesDAG registers the DAG template type resource.
func RegisterTemplateTypesDAG(s *mcp.Server) {
	s.AddResource(TemplateTypesDAGResource(), TemplateTypesDAGHandler())
}

// RegisterTemplateTypesSteps registers the steps template type resource.
func RegisterTemplateTypesSteps(s *mcp.Server) {
	s.AddResource(TemplateTypesStepsResource(), TemplateTypesStepsHandler())
}

// RegisterTemplateTypesSuspend registers the suspend template type resource.
func RegisterTemplateTypesSuspend(s *mcp.Server) {
	s.AddResource(TemplateTypesSuspendResource(), TemplateTypesSuspendHandler())
}

// RegisterTemplateTypesResource registers the resource template type resource.
func RegisterTemplateTypesResource(s *mcp.Server) {
	s.AddResource(TemplateTypesResourceResource(), TemplateTypesResourceHandler())
}

// RegisterTemplateTypesHTTP registers the HTTP template type resource.
func RegisterTemplateTypesHTTP(s *mcp.Server) {
	s.AddResource(TemplateTypesHTTPResource(), TemplateTypesHTTPHandler())
}
