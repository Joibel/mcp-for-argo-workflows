// Package tools implements MCP tool handlers for Argo Workflows operations.
package tools

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/argoproj/argo-workflows/v3/pkg/apiclient/workflowtemplate"
	wfv1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"sigs.k8s.io/yaml"

	"github.com/Joibel/mcp-for-argo-workflows/internal/argo"
)

// CreateWorkflowTemplateInput defines the input parameters for the create_workflow_template tool.
type CreateWorkflowTemplateInput struct {
	// Namespace is the Kubernetes namespace (uses default if not specified).
	Namespace string `json:"namespace,omitempty" jsonschema:"Kubernetes namespace (uses default if not specified)"`

	// Manifest is the WorkflowTemplate YAML manifest.
	Manifest string `json:"manifest" jsonschema:"WorkflowTemplate YAML manifest,required"`
}

// CreateWorkflowTemplateOutput defines the output for the create_workflow_template tool.
type CreateWorkflowTemplateOutput struct {
	// Name is the created workflow template name.
	Name string `json:"name"`

	// Namespace is the namespace where the workflow template was created.
	Namespace string `json:"namespace"`

	// CreatedAt is when the workflow template was created.
	CreatedAt string `json:"createdAt,omitempty"`
}

// CreateWorkflowTemplateTool returns the MCP tool definition for create_workflow_template.
func CreateWorkflowTemplateTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "create_workflow_template",
		Description: "Create a WorkflowTemplate from a YAML manifest. Run lint_workflow first to validate the manifest.",
	}
}

// CreateWorkflowTemplateHandler returns a handler function for the create_workflow_template tool.
func CreateWorkflowTemplateHandler(client argo.ClientInterface) func(context.Context, *mcp.CallToolRequest, CreateWorkflowTemplateInput) (*mcp.CallToolResult, *CreateWorkflowTemplateOutput, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input CreateWorkflowTemplateInput) (*mcp.CallToolResult, *CreateWorkflowTemplateOutput, error) {
		// Validate manifest is provided
		if strings.TrimSpace(input.Manifest) == "" {
			return nil, nil, fmt.Errorf("manifest cannot be empty")
		}

		// Guard against oversized manifests (DoS hardening)
		const maxManifestBytes = 1 << 20 // 1 MiB
		if len(input.Manifest) > maxManifestBytes {
			return nil, nil, fmt.Errorf("manifest too large (%d bytes), max %d", len(input.Manifest), maxManifestBytes)
		}

		// Parse the YAML manifest into a WorkflowTemplate object
		var wft wfv1.WorkflowTemplate
		if err := yaml.Unmarshal([]byte(input.Manifest), &wft); err != nil {
			return nil, nil, fmt.Errorf("failed to parse workflow template manifest: %w", err)
		}

		// Validate that the manifest is a WorkflowTemplate
		if wft.Kind != "" && wft.Kind != "WorkflowTemplate" {
			return nil, nil, fmt.Errorf("manifest must be a WorkflowTemplate, got %q", wft.Kind)
		}

		// Determine namespace
		namespace := ResolveNamespace(input.Namespace, client)
		wft.Namespace = namespace

		// Get the workflow template service client
		wftService, err := client.WorkflowTemplateService()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get workflow template service: %w", err)
		}

		// Create the workflow template
		createdWft, err := wftService.CreateWorkflowTemplate(ctx, &workflowtemplate.WorkflowTemplateCreateRequest{
			Namespace: namespace,
			Template:  &wft,
		})
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create workflow template: %w", err)
		}

		// Build the output
		output := &CreateWorkflowTemplateOutput{
			Name:      createdWft.Name,
			Namespace: createdWft.Namespace,
		}

		// Format timestamp
		if !createdWft.CreationTimestamp.IsZero() {
			output.CreatedAt = createdWft.CreationTimestamp.Format(time.RFC3339)
		}

		// Build human-readable result
		resultText := fmt.Sprintf("WorkflowTemplate %q created in namespace %q", output.Name, output.Namespace)

		return TextResult(resultText), output, nil
	}
}
