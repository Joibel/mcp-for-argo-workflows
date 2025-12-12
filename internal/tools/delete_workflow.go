// Package tools implements MCP tool handlers for Argo Workflows operations.
package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/argoproj/argo-workflows/v3/pkg/apiclient/workflow"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/Joibel/mcp-for-argo-workflows/internal/argo"
)

// DeleteWorkflowInput defines the input parameters for the delete_workflow tool.
type DeleteWorkflowInput struct {
	// Namespace is the Kubernetes namespace (uses default if not specified).
	Namespace string `json:"namespace,omitempty" jsonschema:"Kubernetes namespace (uses default if not specified)"`

	// Name is the workflow name.
	Name string `json:"name" jsonschema:"Workflow name,required"`

	// Force indicates whether to force deletion without waiting for graceful termination.
	Force bool `json:"force,omitempty" jsonschema:"Force deletion without waiting for graceful termination"`
}

// DeleteWorkflowOutput defines the output for the delete_workflow tool.
type DeleteWorkflowOutput struct {
	// Name is the deleted workflow name.
	Name string `json:"name"`

	// Namespace is the namespace where the workflow was deleted.
	Namespace string `json:"namespace"`

	// Message provides confirmation of the deletion.
	Message string `json:"message"`
}

// DeleteWorkflowTool returns the MCP tool definition for delete_workflow.
func DeleteWorkflowTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "delete_workflow",
		Description: "Delete an Argo Workflow",
	}
}

// DeleteWorkflowHandler returns a handler function for the delete_workflow tool.
func DeleteWorkflowHandler(client *argo.Client) func(context.Context, *mcp.CallToolRequest, DeleteWorkflowInput) (*mcp.CallToolResult, *DeleteWorkflowOutput, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input DeleteWorkflowInput) (*mcp.CallToolResult, *DeleteWorkflowOutput, error) {
		// Validate and normalize name
		input.Name = strings.TrimSpace(input.Name)
		if input.Name == "" {
			return nil, nil, fmt.Errorf("workflow name cannot be empty")
		}

		// Determine namespace (trim for consistency with name validation)
		namespace := strings.TrimSpace(input.Namespace)
		if namespace == "" {
			namespace = client.DefaultNamespace()
		}

		// Get the workflow service client
		wfService := client.WorkflowService()

		// Delete the workflow (use client.Context() which contains the KubeClient)
		_, err := wfService.DeleteWorkflow(client.Context(), &workflow.WorkflowDeleteRequest{
			Namespace: namespace,
			Name:      input.Name,
			Force:     input.Force,
		})
		if err != nil {
			return nil, nil, fmt.Errorf("failed to delete workflow: %w", err)
		}

		// Build the output
		output := &DeleteWorkflowOutput{
			Name:      input.Name,
			Namespace: namespace,
			Message:   fmt.Sprintf("Workflow %q deleted successfully", input.Name),
		}

		return nil, output, nil
	}
}
