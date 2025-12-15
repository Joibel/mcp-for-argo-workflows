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
		RegisterExamplesHelloWorld,
		RegisterExamplesMultiStep,
		RegisterExamplesDAGDiamond,
		RegisterExamplesParameters,
		RegisterExamplesArtifacts,
		RegisterExamplesLoops,
		RegisterExamplesConditionals,
		RegisterExamplesRetries,
		RegisterExamplesTimeoutLimits,
		RegisterExamplesResourceManagement,
		RegisterExamplesVolumes,
		RegisterExamplesExitHandlers,
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

// RegisterExamplesHelloWorld registers the hello-world example resource.
func RegisterExamplesHelloWorld(s *mcp.Server) {
	s.AddResource(ExamplesHelloWorldResource(), ExamplesHelloWorldHandler())
}

// RegisterExamplesMultiStep registers the multi-step example resource.
func RegisterExamplesMultiStep(s *mcp.Server) {
	s.AddResource(ExamplesMultiStepResource(), ExamplesMultiStepHandler())
}

// RegisterExamplesDAGDiamond registers the dag-diamond example resource.
func RegisterExamplesDAGDiamond(s *mcp.Server) {
	s.AddResource(ExamplesDAGDiamondResource(), ExamplesDAGDiamondHandler())
}

// RegisterExamplesParameters registers the parameters example resource.
func RegisterExamplesParameters(s *mcp.Server) {
	s.AddResource(ExamplesParametersResource(), ExamplesParametersHandler())
}

// RegisterExamplesArtifacts registers the artifacts example resource.
func RegisterExamplesArtifacts(s *mcp.Server) {
	s.AddResource(ExamplesArtifactsResource(), ExamplesArtifactsHandler())
}

// RegisterExamplesLoops registers the loops example resource.
func RegisterExamplesLoops(s *mcp.Server) {
	s.AddResource(ExamplesLoopsResource(), ExamplesLoopsHandler())
}

// RegisterExamplesConditionals registers the conditionals example resource.
func RegisterExamplesConditionals(s *mcp.Server) {
	s.AddResource(ExamplesConditionalsResource(), ExamplesConditionalsHandler())
}

// RegisterExamplesRetries registers the retries example resource.
func RegisterExamplesRetries(s *mcp.Server) {
	s.AddResource(ExamplesRetriesResource(), ExamplesRetriesHandler())
}

// RegisterExamplesTimeoutLimits registers the timeout-limits example resource.
func RegisterExamplesTimeoutLimits(s *mcp.Server) {
	s.AddResource(ExamplesTimeoutLimitsResource(), ExamplesTimeoutLimitsHandler())
}

// RegisterExamplesResourceManagement registers the resource-management example resource.
func RegisterExamplesResourceManagement(s *mcp.Server) {
	s.AddResource(ExamplesResourceManagementResource(), ExamplesResourceManagementHandler())
}

// RegisterExamplesVolumes registers the volumes example resource.
func RegisterExamplesVolumes(s *mcp.Server) {
	s.AddResource(ExamplesVolumesResource(), ExamplesVolumesHandler())
}

// RegisterExamplesExitHandlers registers the exit-handlers example resource.
func RegisterExamplesExitHandlers(s *mcp.Server) {
	s.AddResource(ExamplesExitHandlersResource(), ExamplesExitHandlersHandler())
}
