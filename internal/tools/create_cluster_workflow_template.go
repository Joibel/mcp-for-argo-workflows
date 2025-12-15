// Package tools implements MCP tool handlers for Argo Workflows operations.
package tools

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/argoproj/argo-workflows/v3/pkg/apiclient/clusterworkflowtemplate"
	wfv1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"sigs.k8s.io/yaml"

	"github.com/Joibel/mcp-for-argo-workflows/internal/argo"
)

// CreateClusterWorkflowTemplateInput defines the input parameters for the create_cluster_workflow_template tool.
type CreateClusterWorkflowTemplateInput struct {
	// Manifest is the ClusterWorkflowTemplate YAML manifest.
	Manifest string `json:"manifest" jsonschema:"ClusterWorkflowTemplate YAML manifest,required"`
}

// CreateClusterWorkflowTemplateOutput defines the output for the create_cluster_workflow_template tool.
type CreateClusterWorkflowTemplateOutput struct {
	// Name is the created cluster workflow template name.
	Name string `json:"name"`

	// CreatedAt is when the cluster workflow template was created.
	CreatedAt string `json:"createdAt,omitempty"`
}

// CreateClusterWorkflowTemplateTool returns the MCP tool definition for create_cluster_workflow_template.
func CreateClusterWorkflowTemplateTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "create_cluster_workflow_template",
		Description: "Create a ClusterWorkflowTemplate from a YAML manifest. Run lint_workflow first to validate the manifest.",
	}
}

// CreateClusterWorkflowTemplateHandler returns a handler function for the create_cluster_workflow_template tool.
func CreateClusterWorkflowTemplateHandler(client argo.ClientInterface) func(context.Context, *mcp.CallToolRequest, CreateClusterWorkflowTemplateInput) (*mcp.CallToolResult, *CreateClusterWorkflowTemplateOutput, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input CreateClusterWorkflowTemplateInput) (*mcp.CallToolResult, *CreateClusterWorkflowTemplateOutput, error) {
		// Validate manifest is provided
		if strings.TrimSpace(input.Manifest) == "" {
			return nil, nil, fmt.Errorf("manifest cannot be empty")
		}

		// Guard against oversized manifests (DoS hardening)
		const maxManifestBytes = 1 << 20 // 1 MiB
		if len(input.Manifest) > maxManifestBytes {
			return nil, nil, fmt.Errorf("manifest too large (%d bytes), max %d", len(input.Manifest), maxManifestBytes)
		}

		// Parse the YAML manifest into a ClusterWorkflowTemplate object
		var cwft wfv1.ClusterWorkflowTemplate
		if err := yaml.Unmarshal([]byte(input.Manifest), &cwft); err != nil {
			return nil, nil, fmt.Errorf("failed to parse cluster workflow template manifest: %w", err)
		}

		// Validate that the manifest is a ClusterWorkflowTemplate
		if cwft.Kind != "" && cwft.Kind != "ClusterWorkflowTemplate" {
			return nil, nil, fmt.Errorf("manifest must be a ClusterWorkflowTemplate, got %q", cwft.Kind)
		}

		// Get the cluster workflow template service client
		cwftService, err := client.ClusterWorkflowTemplateService()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get cluster workflow template service: %w", err)
		}

		// Create the cluster workflow template
		createdCwft, err := cwftService.CreateClusterWorkflowTemplate(ctx, &clusterworkflowtemplate.ClusterWorkflowTemplateCreateRequest{
			Template: &cwft,
		})
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create cluster workflow template: %w", err)
		}

		// Build the output
		output := &CreateClusterWorkflowTemplateOutput{
			Name: createdCwft.Name,
		}

		// Format timestamp
		if !createdCwft.CreationTimestamp.IsZero() {
			output.CreatedAt = createdCwft.CreationTimestamp.Format(time.RFC3339)
		}

		// Build human-readable result
		resultText := fmt.Sprintf("ClusterWorkflowTemplate %q created", output.Name)

		return TextResult(resultText), output, nil
	}
}
