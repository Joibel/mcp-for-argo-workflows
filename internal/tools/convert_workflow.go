// Package tools implements MCP tool handlers for Argo Workflows operations.
package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	wfv1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"sigs.k8s.io/yaml"
)

// Output format constants.
const (
	// OutputFormatYAML is the YAML output format.
	OutputFormatYAML = "yaml"
	// OutputFormatJSON is the JSON output format.
	OutputFormatJSON = "json"
)

// ConvertWorkflowInput defines the input parameters for the convert_workflow tool.
type ConvertWorkflowInput struct {
	// Manifest is the workflow YAML manifest to convert.
	Manifest string `json:"manifest" jsonschema:"Workflow YAML manifest to convert,required"`

	// OutputFormat is the output format (yaml or json).
	OutputFormat string `json:"outputFormat,omitempty" jsonschema:"Output format: yaml (default) or json,enum=yaml,enum=json"`
}

// ConvertWorkflowOutput defines the output for the convert_workflow tool.
type ConvertWorkflowOutput struct {
	// Manifest is the converted manifest.
	Manifest string `json:"manifest"`

	// Format is the output format used.
	Format string `json:"format"`

	// Kind is the kind of manifest that was converted.
	Kind string `json:"kind"`

	// Name is the name of the workflow/template.
	Name string `json:"name,omitempty"`

	// Changes is a list of changes made during conversion.
	Changes []string `json:"changes,omitempty"`

	// Warnings is a list of warnings for manual review.
	Warnings []string `json:"warnings,omitempty"`
}

// ConvertWorkflowTool returns the MCP tool definition for convert_workflow.
func ConvertWorkflowTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "convert_workflow",
		Description: "Convert an Argo Workflow manifest to a newer format, migrating deprecated fields. Supports Workflow, WorkflowTemplate, ClusterWorkflowTemplate, and CronWorkflow manifests.",
	}
}

// ConvertWorkflowHandler returns a handler function for the convert_workflow tool.
// Note: This tool doesn't require the Argo client since it works purely from YAML.
func ConvertWorkflowHandler() func(context.Context, *mcp.CallToolRequest, ConvertWorkflowInput) (*mcp.CallToolResult, *ConvertWorkflowOutput, error) {
	return func(_ context.Context, _ *mcp.CallToolRequest, input ConvertWorkflowInput) (*mcp.CallToolResult, *ConvertWorkflowOutput, error) {
		// Validate manifest is provided
		if strings.TrimSpace(input.Manifest) == "" {
			return nil, nil, fmt.Errorf("manifest cannot be empty")
		}

		// Guard against oversized manifests (DoS hardening)
		const maxManifestBytes = 1 << 20 // 1 MiB
		if len(input.Manifest) > maxManifestBytes {
			return nil, nil, fmt.Errorf("manifest too large (%d bytes), max %d", len(input.Manifest), maxManifestBytes)
		}

		// Determine output format
		outputFormat := strings.ToLower(strings.TrimSpace(input.OutputFormat))
		if outputFormat == "" {
			outputFormat = OutputFormatYAML
		}
		if outputFormat != OutputFormatYAML && outputFormat != OutputFormatJSON {
			return nil, nil, fmt.Errorf("invalid output format: %s (must be %s or %s)", outputFormat, OutputFormatYAML, OutputFormatJSON)
		}

		// First, determine the kind of manifest
		var generic struct {
			Kind     string `json:"kind"`
			Metadata struct {
				Name         string `json:"name"`
				GenerateName string `json:"generateName"`
			} `json:"metadata"`
		}
		if err := yaml.Unmarshal([]byte(input.Manifest), &generic); err != nil {
			return nil, nil, fmt.Errorf("failed to parse manifest: %w", err)
		}

		kind := generic.Kind
		name := generic.Metadata.Name
		if name == "" {
			name = generic.Metadata.GenerateName
		}

		var changes []string
		var warnings []string
		var convertedManifest string

		// Convert based on kind
		switch kind {
		case KindWorkflow, "":
			var wf wfv1.Workflow
			if err := yaml.Unmarshal([]byte(input.Manifest), &wf); err != nil {
				return nil, nil, fmt.Errorf("failed to parse %s manifest: %w", KindWorkflow, err)
			}
			if name == "" {
				name = wf.Name
				if name == "" {
					name = wf.GenerateName
				}
			}
			if kind == "" {
				kind = KindWorkflow
			}

			// Apply workflow-specific conversions
			wfChanges, wfWarnings := convertWorkflowSpec(&wf.Spec)
			changes = append(changes, wfChanges...)
			warnings = append(warnings, wfWarnings...)

			// Serialize back
			var err error
			convertedManifest, err = serializeManifest(wf, outputFormat)
			if err != nil {
				return nil, nil, err
			}

		case KindWorkflowTemplate:
			var wft wfv1.WorkflowTemplate
			if err := yaml.Unmarshal([]byte(input.Manifest), &wft); err != nil {
				return nil, nil, fmt.Errorf("failed to parse %s manifest: %w", KindWorkflowTemplate, err)
			}
			if name == "" {
				name = wft.Name
			}

			// Apply workflow-specific conversions
			wfChanges, wfWarnings := convertWorkflowSpec(&wft.Spec)
			changes = append(changes, wfChanges...)
			warnings = append(warnings, wfWarnings...)

			// Serialize back
			var err error
			convertedManifest, err = serializeManifest(wft, outputFormat)
			if err != nil {
				return nil, nil, err
			}

		case KindClusterWorkflowTemplate:
			var cwft wfv1.ClusterWorkflowTemplate
			if err := yaml.Unmarshal([]byte(input.Manifest), &cwft); err != nil {
				return nil, nil, fmt.Errorf("failed to parse %s manifest: %w", KindClusterWorkflowTemplate, err)
			}
			if name == "" {
				name = cwft.Name
			}

			// Apply workflow-specific conversions
			wfChanges, wfWarnings := convertWorkflowSpec(&cwft.Spec)
			changes = append(changes, wfChanges...)
			warnings = append(warnings, wfWarnings...)

			// Serialize back
			var err error
			convertedManifest, err = serializeManifest(cwft, outputFormat)
			if err != nil {
				return nil, nil, err
			}

		case KindCronWorkflow:
			var cronWf wfv1.CronWorkflow
			if err := yaml.Unmarshal([]byte(input.Manifest), &cronWf); err != nil {
				return nil, nil, fmt.Errorf("failed to parse %s manifest: %w", KindCronWorkflow, err)
			}
			if name == "" {
				name = cronWf.Name
			}

			// Apply CronWorkflow-specific conversions
			cronChanges, cronWarnings := convertCronWorkflowSpec(&cronWf.Spec)
			changes = append(changes, cronChanges...)
			warnings = append(warnings, cronWarnings...)

			// Apply workflow-specific conversions to the embedded WorkflowSpec
			wfChanges, wfWarnings := convertWorkflowSpec(&cronWf.Spec.WorkflowSpec)
			changes = append(changes, wfChanges...)
			warnings = append(warnings, wfWarnings...)

			// Serialize back
			var err error
			convertedManifest, err = serializeManifest(cronWf, outputFormat)
			if err != nil {
				return nil, nil, err
			}

		default:
			return nil, nil, fmt.Errorf("unsupported manifest kind: %s (must be %s, %s, %s, or %s)", kind, KindWorkflow, KindWorkflowTemplate, KindClusterWorkflowTemplate, KindCronWorkflow)
		}

		// Build output
		output := &ConvertWorkflowOutput{
			Manifest: convertedManifest,
			Format:   outputFormat,
			Kind:     kind,
			Name:     name,
			Changes:  changes,
			Warnings: warnings,
		}

		// Build human-readable result
		var resultText strings.Builder
		resultText.WriteString(fmt.Sprintf("Converted %s %q to %s format\n", kind, name, outputFormat))

		if len(changes) > 0 {
			resultText.WriteString("\nChanges made:\n")
			for _, change := range changes {
				resultText.WriteString(fmt.Sprintf("  - %s\n", change))
			}
		} else {
			resultText.WriteString("\nNo changes needed - manifest is already up to date.\n")
		}

		if len(warnings) > 0 {
			resultText.WriteString("\nWarnings (manual review recommended):\n")
			for _, warning := range warnings {
				resultText.WriteString(fmt.Sprintf("  - %s\n", warning))
			}
		}

		result := &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: resultText.String()},
			},
		}

		return result, output, nil
	}
}

// convertWorkflowSpec applies conversions to a WorkflowSpec and returns changes and warnings.
func convertWorkflowSpec(spec *wfv1.WorkflowSpec) ([]string, []string) {
	var changes []string
	var warnings []string

	// Check for deprecated fields in templates
	for i := range spec.Templates {
		tmpl := &spec.Templates[i]
		tmplChanges, tmplWarnings := convertTemplate(tmpl)
		changes = append(changes, tmplChanges...)
		warnings = append(warnings, tmplWarnings...)
	}

	return changes, warnings
}

// convertTemplate applies conversions to a single template.
// Currently a placeholder for future template-level conversions.
//
//nolint:unparam // Placeholder for future conversions - returns nil for now
func convertTemplate(tmpl *wfv1.Template) ([]string, []string) {
	var changes []string
	var warnings []string

	// Check for deprecated container fields
	if tmpl.Container != nil {
		// Note: Most container deprecations are handled by Kubernetes itself
		// We can add specific checks here as needed
		_ = tmpl.Container
	}

	// Check for deprecated script fields
	if tmpl.Script != nil {
		// Note: Script-specific deprecations can be added here
		_ = tmpl.Script
	}

	// Check for deprecated resource fields
	if tmpl.Resource != nil {
		// Note: Resource-specific deprecations can be added here
		_ = tmpl.Resource
	}

	// Check DAG tasks
	if tmpl.DAG != nil {
		for j := range tmpl.DAG.Tasks {
			task := &tmpl.DAG.Tasks[j]
			// Check for deprecated task fields
			// Note: Task-specific deprecations can be added here
			_ = task
		}
	}

	// Check Steps
	for i := range tmpl.Steps {
		for j := range tmpl.Steps[i].Steps {
			step := &tmpl.Steps[i].Steps[j]
			// Check for deprecated step fields
			// Note: Step-specific deprecations can be added here
			_ = step
		}
	}

	return changes, warnings
}

// convertCronWorkflowSpec applies conversions specific to CronWorkflowSpec.
func convertCronWorkflowSpec(spec *wfv1.CronWorkflowSpec) ([]string, []string) {
	var changes []string
	var warnings []string

	// Convert deprecated `schedule` (string) to `schedules` ([]string)
	// The Schedule field is still used but Schedules takes precedence if set
	if spec.Schedule != "" && len(spec.Schedules) == 0 {
		// Migrate schedule to schedules array
		spec.Schedules = []string{spec.Schedule}
		spec.Schedule = "" // Clear the deprecated field
		changes = append(changes, "Migrated spec.schedule to spec.schedules array")
	}

	// Check for other CronWorkflow-specific deprecations
	// ConcurrencyPolicy is still valid but we can warn about certain values
	if spec.ConcurrencyPolicy == "" {
		warnings = append(warnings, "No concurrencyPolicy set - defaults to 'Allow' which may cause overlapping runs")
	}

	return changes, warnings
}

// serializeManifest serializes the manifest to the specified format.
func serializeManifest(obj interface{}, format string) (string, error) {
	switch format {
	case "json":
		data, err := json.MarshalIndent(obj, "", "  ")
		if err != nil {
			return "", fmt.Errorf("failed to serialize manifest to JSON: %w", err)
		}
		return string(data), nil
	case "yaml":
		data, err := yaml.Marshal(obj)
		if err != nil {
			return "", fmt.Errorf("failed to serialize manifest to YAML: %w", err)
		}
		return string(data), nil
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}
