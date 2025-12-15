// Package tools implements MCP tool handlers for Argo Workflows operations.
package tools

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/argoproj/argo-workflows/v3/pkg/apiclient/cronworkflow"
	wfv1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"sigs.k8s.io/yaml"

	"github.com/Joibel/mcp-for-argo-workflows/internal/argo"
)

// CreateCronWorkflowInput defines the input parameters for the create_cron_workflow tool.
type CreateCronWorkflowInput struct {
	// Manifest is the YAML manifest of the CronWorkflow to create (required).
	Manifest string `json:"manifest" jsonschema:"CronWorkflow YAML manifest,required"`

	// Namespace is the Kubernetes namespace (uses default if not specified).
	Namespace string `json:"namespace,omitempty" jsonschema:"Kubernetes namespace (uses default if not specified)"`
}

// CreateCronWorkflowOutput defines the output for the create_cron_workflow tool.
type CreateCronWorkflowOutput struct {
	// Labels are the labels on the created CronWorkflow.
	Labels map[string]string `json:"labels,omitempty"`

	// Annotations are the annotations on the created CronWorkflow.
	Annotations map[string]string `json:"annotations,omitempty"`

	// CreatedAt is when the CronWorkflow was created.
	CreatedAt string `json:"createdAt"`

	// Name is the name of the created CronWorkflow.
	Name string `json:"name"`

	// Namespace is the namespace of the created CronWorkflow.
	Namespace string `json:"namespace"`

	// Schedule is the cron schedule expression.
	Schedule string `json:"schedule"`

	// Timezone is the timezone for the schedule.
	Timezone string `json:"timezone,omitempty"`

	// Entrypoint is the workflow entrypoint.
	Entrypoint string `json:"entrypoint,omitempty"`

	// ConcurrencyPolicy defines how to treat concurrent executions.
	ConcurrencyPolicy string `json:"concurrencyPolicy,omitempty"`

	// Suspended indicates whether the CronWorkflow is suspended.
	Suspended bool `json:"suspended"`
}

// CreateCronWorkflowTool returns the MCP tool definition for create_cron_workflow.
func CreateCronWorkflowTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "create_cron_workflow",
		Description: "Create a new CronWorkflow from a YAML manifest",
	}
}

// CreateCronWorkflowHandler returns a handler function for the create_cron_workflow tool.
func CreateCronWorkflowHandler(client argo.ClientInterface) func(context.Context, *mcp.CallToolRequest, CreateCronWorkflowInput) (*mcp.CallToolResult, *CreateCronWorkflowOutput, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input CreateCronWorkflowInput) (*mcp.CallToolResult, *CreateCronWorkflowOutput, error) {
		// Validate manifest is provided
		if strings.TrimSpace(input.Manifest) == "" {
			return nil, nil, fmt.Errorf("manifest cannot be empty")
		}

		// Guard against oversized manifests (DoS hardening)
		const maxManifestBytes = 1 << 20 // 1 MiB
		if len(input.Manifest) > maxManifestBytes {
			return nil, nil, fmt.Errorf("manifest exceeds maximum size of %d bytes", maxManifestBytes)
		}

		// Parse the YAML manifest
		var cronWf wfv1.CronWorkflow
		if err := yaml.UnmarshalStrict([]byte(input.Manifest), &cronWf); err != nil {
			return nil, nil, fmt.Errorf("failed to parse CronWorkflow manifest: %w", err)
		}

		// Validate that the manifest is a CronWorkflow
		if cronWf.Kind != "" && cronWf.Kind != "CronWorkflow" {
			return nil, nil, fmt.Errorf("manifest kind must be CronWorkflow, got %q", cronWf.Kind)
		}

		// Validate name
		if cronWf.Name == "" {
			return nil, nil, fmt.Errorf("CronWorkflow name is required in manifest")
		}

		// Validate schedule
		if cronWf.Spec.Schedule == "" {
			return nil, nil, fmt.Errorf("CronWorkflow schedule is required in manifest")
		}

		// Resolve namespace - prefer input namespace, then manifest namespace, then default
		namespace := input.Namespace
		if namespace == "" {
			namespace = cronWf.Namespace
		}
		namespace = ResolveNamespace(namespace, client)
		cronWf.Namespace = namespace

		// Get the cron workflow service client
		cronService, err := client.CronWorkflowService()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get cron workflow service: %w", err)
		}

		// Create the cron workflow
		created, err := cronService.CreateCronWorkflow(ctx, &cronworkflow.CreateCronWorkflowRequest{
			Namespace:    namespace,
			CronWorkflow: &cronWf,
		})
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create cron workflow: %w", err)
		}

		// Build output
		output := &CreateCronWorkflowOutput{
			Name:              created.Name,
			Namespace:         created.Namespace,
			Schedule:          created.Spec.Schedule,
			Timezone:          created.Spec.Timezone,
			ConcurrencyPolicy: string(created.Spec.ConcurrencyPolicy),
			Suspended:         created.Spec.Suspend,
			Labels:            created.Labels,
			Annotations:       created.Annotations,
		}

		// Format creation timestamp
		if !created.CreationTimestamp.IsZero() {
			output.CreatedAt = created.CreationTimestamp.Format(time.RFC3339)
		}

		// Get entrypoint if available
		if created.Spec.WorkflowSpec.Entrypoint != "" {
			output.Entrypoint = created.Spec.WorkflowSpec.Entrypoint
		}

		// Build human-readable result
		resultText := fmt.Sprintf("Created CronWorkflow %q in namespace %q", output.Name, output.Namespace)
		resultText += fmt.Sprintf("\nSchedule: %s", output.Schedule)
		if output.Timezone != "" {
			resultText += fmt.Sprintf(" (%s)", output.Timezone)
		}
		if output.Suspended {
			resultText += "\nStatus: Suspended"
		} else {
			resultText += "\nStatus: Active"
		}

		result := &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: resultText},
			},
		}

		return result, output, nil
	}
}
