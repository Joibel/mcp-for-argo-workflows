// Package tools implements MCP tool handlers for Argo Workflows operations.
package tools

import (
	"context"
	"fmt"
	"time"

	"github.com/argoproj/argo-workflows/v3/pkg/apiclient/workflowarchive"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/Joibel/mcp-for-argo-workflows/internal/argo"
)

// ListArchivedWorkflowsInput defines the input parameters for the list_archived_workflows tool.
type ListArchivedWorkflowsInput struct {
	Namespace string `json:"namespace,omitempty" jsonschema:"Kubernetes namespace filter"`
	Labels    string `json:"labels,omitempty" jsonschema:"Label selector (e.g. 'app=myapp,env=prod')"`
	Limit     int64  `json:"limit,omitempty" jsonschema:"Maximum number of results"`
}

// ArchivedWorkflowSummary represents a concise summary of an archived workflow.
type ArchivedWorkflowSummary struct {
	// UID is the unique identifier of the archived workflow.
	UID string `json:"uid"`

	// Name is the workflow name.
	Name string `json:"name"`

	// Namespace is the namespace where the workflow existed.
	Namespace string `json:"namespace"`

	// Phase is the final workflow phase.
	Phase string `json:"phase"`

	// CreatedAt is when the workflow was created.
	CreatedAt string `json:"createdAt"`

	// FinishedAt is when the workflow finished.
	FinishedAt string `json:"finishedAt,omitempty"`

	// Message provides additional status information.
	Message string `json:"message,omitempty"`
}

// ListArchivedWorkflowsOutput defines the output for the list_archived_workflows tool.
type ListArchivedWorkflowsOutput struct {
	// Workflows is the list of archived workflow summaries.
	Workflows []ArchivedWorkflowSummary `json:"workflows"`

	// Total is the total number of archived workflows matching the criteria.
	Total int `json:"total"`
}

// ListArchivedWorkflowsTool returns the MCP tool definition for list_archived_workflows.
func ListArchivedWorkflowsTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "list_archived_workflows",
		Description: "List workflows from the workflow archive. Requires Argo Server connection (not available in direct K8s mode).",
	}
}

// ListArchivedWorkflowsHandler returns a handler function for the list_archived_workflows tool.
func ListArchivedWorkflowsHandler(client argo.ClientInterface) func(context.Context, *mcp.CallToolRequest, ListArchivedWorkflowsInput) (*mcp.CallToolResult, *ListArchivedWorkflowsOutput, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input ListArchivedWorkflowsInput) (*mcp.CallToolResult, *ListArchivedWorkflowsOutput, error) {
		// Get the archived workflow service client
		archiveService, err := client.ArchivedWorkflowService()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get archived workflow service: %w", err)
		}

		// Build list options
		listOpts := &metav1.ListOptions{}

		// Apply label selector
		if input.Labels != "" {
			listOpts.LabelSelector = input.Labels
		}

		// Apply limit
		if input.Limit > 0 {
			listOpts.Limit = input.Limit
		}

		// Build request
		req := &workflowarchive.ListArchivedWorkflowsRequest{
			Namespace:   input.Namespace,
			ListOptions: listOpts,
		}

		// List archived workflows
		listResp, err := archiveService.ListArchivedWorkflows(ctx, req)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to list archived workflows: %w", err)
		}

		// Convert to summaries
		summaries := make([]ArchivedWorkflowSummary, 0, len(listResp.Items))
		for _, wf := range listResp.Items {
			summary := ArchivedWorkflowSummary{
				UID:       string(wf.UID),
				Name:      wf.Name,
				Namespace: wf.Namespace,
				Phase:     string(wf.Status.Phase),
				Message:   wf.Status.Message,
			}

			// Format timestamps
			if !wf.CreationTimestamp.IsZero() {
				summary.CreatedAt = wf.CreationTimestamp.Format(time.RFC3339)
			}
			if !wf.Status.FinishedAt.IsZero() {
				summary.FinishedAt = wf.Status.FinishedAt.Format(time.RFC3339)
			}

			summaries = append(summaries, summary)
		}

		// Build output
		output := &ListArchivedWorkflowsOutput{
			Workflows: summaries,
			Total:     len(summaries),
		}

		// Build human-readable result
		resultText := fmt.Sprintf("Found %d archived workflow(s)", output.Total)
		if input.Namespace != "" {
			resultText += fmt.Sprintf(" in namespace %q", input.Namespace)
		} else {
			resultText += " across all namespaces"
		}

		result := &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: resultText},
			},
		}

		return result, output, nil
	}
}
