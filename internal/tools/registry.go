// Package tools implements MCP tool handlers for Argo Workflows operations.
package tools

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/Joibel/mcp-for-argo-workflows/internal/argo"
)

// ToolRegistrar is a function that registers a tool with the MCP server.
// Each tool provides its own registrar that calls mcp.AddTool with the correct types.
type ToolRegistrar func(s *mcp.Server, client argo.ClientInterface)

// AllTools returns all tool registrars in the order they should be registered.
func AllTools() []ToolRegistrar {
	return []ToolRegistrar{
		RegisterSubmitWorkflow,
		RegisterListWorkflows,
		RegisterGetWorkflow,
		RegisterDeleteWorkflow,
		RegisterWatchWorkflow,
		RegisterLogsWorkflow,
		RegisterWaitWorkflow,
		RegisterLintWorkflow,
		RegisterRetryWorkflow,
		RegisterResubmitWorkflow,
		RegisterSuspendWorkflow,
		RegisterResumeWorkflow,
		RegisterStopWorkflow,
		RegisterTerminateWorkflow,
		RegisterRenderWorkflowGraph,
		RegisterListWorkflowTemplates,
		RegisterGetWorkflowTemplate,
		RegisterCreateWorkflowTemplate,
		RegisterDeleteWorkflowTemplate,
		RegisterListClusterWorkflowTemplates,
		RegisterGetClusterWorkflowTemplate,
		RegisterCreateClusterWorkflowTemplate,
		RegisterDeleteClusterWorkflowTemplate,
		RegisterListCronWorkflows,
		RegisterGetCronWorkflow,
		RegisterCreateCronWorkflow,
		RegisterDeleteCronWorkflow,
		RegisterSuspendCronWorkflow,
		RegisterListCronWorkflows,
		RegisterGetCronWorkflow,
		RegisterCreateCronWorkflow,
	}
}

// RegisterAll registers all tools with the MCP server.
func RegisterAll(s *mcp.Server, client argo.ClientInterface) {
	for _, register := range AllTools() {
		register(s, client)
	}
}

// Individual tool registrars - these wrap mcp.AddTool with the correct type parameters.

// RegisterSubmitWorkflow registers the submit_workflow tool.
func RegisterSubmitWorkflow(s *mcp.Server, client argo.ClientInterface) {
	mcp.AddTool(s, SubmitWorkflowTool(), SubmitWorkflowHandler(client))
}

// RegisterListWorkflows registers the list_workflows tool.
func RegisterListWorkflows(s *mcp.Server, client argo.ClientInterface) {
	mcp.AddTool(s, ListWorkflowsTool(), ListWorkflowsHandler(client))
}

// RegisterGetWorkflow registers the get_workflow tool.
func RegisterGetWorkflow(s *mcp.Server, client argo.ClientInterface) {
	mcp.AddTool(s, GetWorkflowTool(), GetWorkflowHandler(client))
}

// RegisterDeleteWorkflow registers the delete_workflow tool.
func RegisterDeleteWorkflow(s *mcp.Server, client argo.ClientInterface) {
	mcp.AddTool(s, DeleteWorkflowTool(), DeleteWorkflowHandler(client))
}

// RegisterWatchWorkflow registers the watch_workflow tool.
func RegisterWatchWorkflow(s *mcp.Server, client argo.ClientInterface) {
	mcp.AddTool(s, WatchWorkflowTool(), WatchWorkflowHandler(client))
}

// RegisterLogsWorkflow registers the logs_workflow tool.
func RegisterLogsWorkflow(s *mcp.Server, client argo.ClientInterface) {
	mcp.AddTool(s, LogsWorkflowTool(), LogsWorkflowHandler(client))
}

// RegisterWaitWorkflow registers the wait_workflow tool.
func RegisterWaitWorkflow(s *mcp.Server, client argo.ClientInterface) {
	mcp.AddTool(s, WaitWorkflowTool(), WaitWorkflowHandler(client))
}

// RegisterLintWorkflow registers the lint_workflow tool.
func RegisterLintWorkflow(s *mcp.Server, client argo.ClientInterface) {
	mcp.AddTool(s, LintWorkflowTool(), LintWorkflowHandler(client))
}

// RegisterRetryWorkflow registers the retry_workflow tool.
func RegisterRetryWorkflow(s *mcp.Server, client argo.ClientInterface) {
	mcp.AddTool(s, RetryWorkflowTool(), RetryWorkflowHandler(client))
}

// RegisterResubmitWorkflow registers the resubmit_workflow tool.
func RegisterResubmitWorkflow(s *mcp.Server, client argo.ClientInterface) {
	mcp.AddTool(s, ResubmitWorkflowTool(), ResubmitWorkflowHandler(client))
}

// RegisterSuspendWorkflow registers the suspend_workflow tool.
func RegisterSuspendWorkflow(s *mcp.Server, client argo.ClientInterface) {
	mcp.AddTool(s, SuspendWorkflowTool(), SuspendWorkflowHandler(client))
}

// RegisterResumeWorkflow registers the resume_workflow tool.
func RegisterResumeWorkflow(s *mcp.Server, client argo.ClientInterface) {
	mcp.AddTool(s, ResumeWorkflowTool(), ResumeWorkflowHandler(client))
}

// RegisterStopWorkflow registers the stop_workflow tool.
func RegisterStopWorkflow(s *mcp.Server, client argo.ClientInterface) {
	mcp.AddTool(s, StopWorkflowTool(), StopWorkflowHandler(client))
}

// RegisterTerminateWorkflow registers the terminate_workflow tool.
func RegisterTerminateWorkflow(s *mcp.Server, client argo.ClientInterface) {
	mcp.AddTool(s, TerminateWorkflowTool(), TerminateWorkflowHandler(client))
}

// RegisterListWorkflowTemplates registers the list_workflow_templates tool.
func RegisterListWorkflowTemplates(s *mcp.Server, client argo.ClientInterface) {
	mcp.AddTool(s, ListWorkflowTemplatesTool(), ListWorkflowTemplatesHandler(client))
}

// RegisterGetWorkflowTemplate registers the get_workflow_template tool.
func RegisterGetWorkflowTemplate(s *mcp.Server, client argo.ClientInterface) {
	mcp.AddTool(s, GetWorkflowTemplateTool(), GetWorkflowTemplateHandler(client))
}

// RegisterCreateWorkflowTemplate registers the create_workflow_template tool.
func RegisterCreateWorkflowTemplate(s *mcp.Server, client argo.ClientInterface) {
	mcp.AddTool(s, CreateWorkflowTemplateTool(), CreateWorkflowTemplateHandler(client))
}

// RegisterDeleteWorkflowTemplate registers the delete_workflow_template tool.
func RegisterDeleteWorkflowTemplate(s *mcp.Server, client argo.ClientInterface) {
	mcp.AddTool(s, DeleteWorkflowTemplateTool(), DeleteWorkflowTemplateHandler(client))
}

// RegisterListClusterWorkflowTemplates registers the list_cluster_workflow_templates tool.
func RegisterListClusterWorkflowTemplates(s *mcp.Server, client argo.ClientInterface) {
	mcp.AddTool(s, ListClusterWorkflowTemplatesTool(), ListClusterWorkflowTemplatesHandler(client))
}

// RegisterRenderWorkflowGraph registers the render_workflow_graph tool.
func RegisterRenderWorkflowGraph(s *mcp.Server, client argo.ClientInterface) {
	mcp.AddTool(s, RenderWorkflowGraphTool(), RenderWorkflowGraphHandler(client))
}

// RegisterGetClusterWorkflowTemplate registers the get_cluster_workflow_template tool.
func RegisterGetClusterWorkflowTemplate(s *mcp.Server, client argo.ClientInterface) {
	mcp.AddTool(s, GetClusterWorkflowTemplateTool(), GetClusterWorkflowTemplateHandler(client))
}

// RegisterCreateClusterWorkflowTemplate registers the create_cluster_workflow_template tool.
func RegisterCreateClusterWorkflowTemplate(s *mcp.Server, client argo.ClientInterface) {
	mcp.AddTool(s, CreateClusterWorkflowTemplateTool(), CreateClusterWorkflowTemplateHandler(client))
}

// RegisterDeleteClusterWorkflowTemplate registers the delete_cluster_workflow_template tool.
func RegisterDeleteClusterWorkflowTemplate(s *mcp.Server, client argo.ClientInterface) {
	mcp.AddTool(s, DeleteClusterWorkflowTemplateTool(), DeleteClusterWorkflowTemplateHandler(client))
}

// RegisterListCronWorkflows registers the list_cron_workflows tool.
func RegisterListCronWorkflows(s *mcp.Server, client argo.ClientInterface) {
	mcp.AddTool(s, ListCronWorkflowsTool(), ListCronWorkflowsHandler(client))
}

// RegisterGetCronWorkflow registers the get_cron_workflow tool.
func RegisterGetCronWorkflow(s *mcp.Server, client argo.ClientInterface) {
	mcp.AddTool(s, GetCronWorkflowTool(), GetCronWorkflowHandler(client))
}

// RegisterCreateCronWorkflow registers the create_cron_workflow tool.
func RegisterCreateCronWorkflow(s *mcp.Server, client argo.ClientInterface) {
	mcp.AddTool(s, CreateCronWorkflowTool(), CreateCronWorkflowHandler(client))
}

// RegisterDeleteCronWorkflow registers the delete_cron_workflow tool.
func RegisterDeleteCronWorkflow(s *mcp.Server, client argo.ClientInterface) {
	mcp.AddTool(s, DeleteCronWorkflowTool(), DeleteCronWorkflowHandler(client))
}

// RegisterSuspendCronWorkflow registers the suspend_cron_workflow tool.
func RegisterSuspendCronWorkflow(s *mcp.Server, client argo.ClientInterface) {
	mcp.AddTool(s, SuspendCronWorkflowTool(), SuspendCronWorkflowHandler(client))
}
