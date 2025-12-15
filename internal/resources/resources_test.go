package resources

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorkflowSchemaResource(t *testing.T) {
	resource := WorkflowSchemaResource()

	assert.Equal(t, "argo://schemas/workflow", resource.URI)
	assert.Equal(t, "workflow-schema", resource.Name)
	assert.Equal(t, "Argo Workflow CRD Schema", resource.Title)
	assert.Equal(t, "text/markdown", resource.MIMEType)
	assert.NotEmpty(t, resource.Description)
}

func TestWorkflowSchemaHandler(t *testing.T) {
	handler := WorkflowSchemaHandler()
	ctx := context.Background()

	t.Run("valid URI", func(t *testing.T) {
		req := &mcp.ReadResourceRequest{
			Params: &mcp.ReadResourceParams{
				URI: "argo://schemas/workflow",
			},
		}

		result, err := handler(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, result)
		require.Len(t, result.Contents, 1)

		content := result.Contents[0]
		assert.Equal(t, "argo://schemas/workflow", content.URI)
		assert.Equal(t, "text/markdown", content.MIMEType)
		assert.NotEmpty(t, content.Text)
		assert.Contains(t, content.Text, "# Workflow CRD Schema")
		assert.Contains(t, content.Text, "apiVersion")
		assert.Contains(t, content.Text, "kind")
		assert.Contains(t, content.Text, "spec")
	})

	t.Run("invalid URI", func(t *testing.T) {
		req := &mcp.ReadResourceRequest{
			Params: &mcp.ReadResourceParams{
				URI: "argo://schemas/invalid",
			},
		}

		result, err := handler(ctx, req)
		require.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestWorkflowTemplateSchemaResource(t *testing.T) {
	resource := WorkflowTemplateSchemaResource()

	assert.Equal(t, "argo://schemas/workflow-template", resource.URI)
	assert.Equal(t, "workflow-template-schema", resource.Name)
	assert.Equal(t, "Argo WorkflowTemplate CRD Schema", resource.Title)
	assert.Equal(t, "text/markdown", resource.MIMEType)
}

func TestWorkflowTemplateSchemaHandler(t *testing.T) {
	handler := WorkflowTemplateSchemaHandler()
	ctx := context.Background()

	req := &mcp.ReadResourceRequest{
		Params: &mcp.ReadResourceParams{
			URI: "argo://schemas/workflow-template",
		},
	}

	result, err := handler(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Contents, 1)

	content := result.Contents[0]
	assert.Equal(t, "argo://schemas/workflow-template", content.URI)
	assert.NotEmpty(t, content.Text)
	assert.Contains(t, content.Text, "# WorkflowTemplate CRD Schema")
	assert.Contains(t, content.Text, "No Status")
	assert.Contains(t, content.Text, "Reusable")
}

func TestClusterWorkflowTemplateSchemaResource(t *testing.T) {
	resource := ClusterWorkflowTemplateSchemaResource()

	assert.Equal(t, "argo://schemas/cluster-workflow-template", resource.URI)
	assert.Equal(t, "cluster-workflow-template-schema", resource.Name)
	assert.Equal(t, "Argo ClusterWorkflowTemplate CRD Schema", resource.Title)
	assert.Equal(t, "text/markdown", resource.MIMEType)
}

func TestClusterWorkflowTemplateSchemaHandler(t *testing.T) {
	handler := ClusterWorkflowTemplateSchemaHandler()
	ctx := context.Background()

	req := &mcp.ReadResourceRequest{
		Params: &mcp.ReadResourceParams{
			URI: "argo://schemas/cluster-workflow-template",
		},
	}

	result, err := handler(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Contents, 1)

	content := result.Contents[0]
	assert.Equal(t, "argo://schemas/cluster-workflow-template", content.URI)
	assert.NotEmpty(t, content.Text)
	assert.Contains(t, content.Text, "# ClusterWorkflowTemplate CRD Schema")
	assert.Contains(t, content.Text, "Cluster-Scoped")
	assert.Contains(t, content.Text, "RBAC")
}

func TestCronWorkflowSchemaResource(t *testing.T) {
	resource := CronWorkflowSchemaResource()

	assert.Equal(t, "argo://schemas/cron-workflow", resource.URI)
	assert.Equal(t, "cron-workflow-schema", resource.Name)
	assert.Equal(t, "Argo CronWorkflow CRD Schema", resource.Title)
	assert.Equal(t, "text/markdown", resource.MIMEType)
}

func TestCronWorkflowSchemaHandler(t *testing.T) {
	handler := CronWorkflowSchemaHandler()
	ctx := context.Background()

	req := &mcp.ReadResourceRequest{
		Params: &mcp.ReadResourceParams{
			URI: "argo://schemas/cron-workflow",
		},
	}

	result, err := handler(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Contents, 1)

	content := result.Contents[0]
	assert.Equal(t, "argo://schemas/cron-workflow", content.URI)
	assert.NotEmpty(t, content.Text)
	assert.Contains(t, content.Text, "# CronWorkflow CRD Schema")
	assert.Contains(t, content.Text, "schedule")
	assert.Contains(t, content.Text, "timezone")
	assert.Contains(t, content.Text, "concurrencyPolicy")
}

func TestAllResources(t *testing.T) {
	registrars := AllResources()

	// Should have exactly 24 resource registrars (4 schema + 8 template types + 12 examples)
	assert.Len(t, registrars, 24, "Expected 24 resource registrars")

	// Verify all are not nil
	for i, registrar := range registrars {
		assert.NotNil(t, registrar, "Registrar at index %d should not be nil", i)
	}
}

func TestRegisterAll(t *testing.T) {
	// Create a test MCP server
	implementation := &mcp.Implementation{
		Name:    "test-server",
		Version: "1.0.0",
	}
	server := mcp.NewServer(implementation, nil)

	// Register all resources
	RegisterAll(server)

	// The server should now have resources registered
	// This is a smoke test to ensure registration doesn't panic
}

// Template Types Overview Tests

func TestTemplateTypesOverviewResource(t *testing.T) {
	resource := TemplateTypesOverviewResource()

	assert.Equal(t, "argo://docs/template-types", resource.URI)
	assert.Equal(t, "template-types-overview", resource.Name)
	assert.Equal(t, "Argo Workflows Template Types Overview", resource.Title)
	assert.Equal(t, "text/markdown", resource.MIMEType)
	assert.NotEmpty(t, resource.Description)
}

func TestTemplateTypesOverviewHandler(t *testing.T) {
	handler := TemplateTypesOverviewHandler()
	ctx := context.Background()

	t.Run("valid URI", func(t *testing.T) {
		req := &mcp.ReadResourceRequest{
			Params: &mcp.ReadResourceParams{
				URI: "argo://docs/template-types",
			},
		}

		result, err := handler(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, result)
		require.Len(t, result.Contents, 1)

		content := result.Contents[0]
		assert.Equal(t, "argo://docs/template-types", content.URI)
		assert.Equal(t, "text/markdown", content.MIMEType)
		assert.NotEmpty(t, content.Text)
		assert.Contains(t, content.Text, "Template Types Overview")
		assert.Contains(t, content.Text, "Container")
		assert.Contains(t, content.Text, "Script")
		assert.Contains(t, content.Text, "DAG")
		assert.Contains(t, content.Text, "Steps")
	})

	t.Run("invalid URI", func(t *testing.T) {
		req := &mcp.ReadResourceRequest{
			Params: &mcp.ReadResourceParams{
				URI: "argo://docs/invalid",
			},
		}

		result, err := handler(ctx, req)
		require.Error(t, err)
		assert.Nil(t, result)
	})
}

// Container Template Tests

func TestTemplateTypesContainerResource(t *testing.T) {
	resource := TemplateTypesContainerResource()

	assert.Equal(t, "argo://docs/template-types/container", resource.URI)
	assert.Equal(t, "template-types-container", resource.Name)
	assert.Equal(t, "Container Template Type", resource.Title)
	assert.Equal(t, "text/markdown", resource.MIMEType)
}

func TestTemplateTypesContainerHandler(t *testing.T) {
	handler := TemplateTypesContainerHandler()
	ctx := context.Background()

	req := &mcp.ReadResourceRequest{
		Params: &mcp.ReadResourceParams{
			URI: "argo://docs/template-types/container",
		},
	}

	result, err := handler(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Contents, 1)

	content := result.Contents[0]
	assert.Equal(t, "argo://docs/template-types/container", content.URI)
	assert.NotEmpty(t, content.Text)
	assert.Contains(t, content.Text, "Container Template Type")
	assert.Contains(t, content.Text, "image")
	assert.Contains(t, content.Text, "command")
	assert.Contains(t, content.Text, "args")
}

// Script Template Tests

func TestTemplateTypesScriptResource(t *testing.T) {
	resource := TemplateTypesScriptResource()

	assert.Equal(t, "argo://docs/template-types/script", resource.URI)
	assert.Equal(t, "template-types-script", resource.Name)
	assert.Equal(t, "Script Template Type", resource.Title)
	assert.Equal(t, "text/markdown", resource.MIMEType)
}

func TestTemplateTypesScriptHandler(t *testing.T) {
	handler := TemplateTypesScriptHandler()
	ctx := context.Background()

	req := &mcp.ReadResourceRequest{
		Params: &mcp.ReadResourceParams{
			URI: "argo://docs/template-types/script",
		},
	}

	result, err := handler(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Contents, 1)

	content := result.Contents[0]
	assert.Equal(t, "argo://docs/template-types/script", content.URI)
	assert.NotEmpty(t, content.Text)
	assert.Contains(t, content.Text, "Script Template Type")
	assert.Contains(t, content.Text, "source")
	assert.Contains(t, content.Text, "inline")
}

// DAG Template Tests

func TestTemplateTypesDAGResource(t *testing.T) {
	resource := TemplateTypesDAGResource()

	assert.Equal(t, "argo://docs/template-types/dag", resource.URI)
	assert.Equal(t, "template-types-dag", resource.Name)
	assert.Equal(t, "DAG Template Type", resource.Title)
	assert.Equal(t, "text/markdown", resource.MIMEType)
}

func TestTemplateTypesDAGHandler(t *testing.T) {
	handler := TemplateTypesDAGHandler()
	ctx := context.Background()

	req := &mcp.ReadResourceRequest{
		Params: &mcp.ReadResourceParams{
			URI: "argo://docs/template-types/dag",
		},
	}

	result, err := handler(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Contents, 1)

	content := result.Contents[0]
	assert.Equal(t, "argo://docs/template-types/dag", content.URI)
	assert.NotEmpty(t, content.Text)
	assert.Contains(t, content.Text, "DAG Template Type")
	assert.Contains(t, content.Text, "dependencies")
	assert.Contains(t, content.Text, "tasks")
}

// Steps Template Tests

func TestTemplateTypesStepsResource(t *testing.T) {
	resource := TemplateTypesStepsResource()

	assert.Equal(t, "argo://docs/template-types/steps", resource.URI)
	assert.Equal(t, "template-types-steps", resource.Name)
	assert.Equal(t, "Steps Template Type", resource.Title)
	assert.Equal(t, "text/markdown", resource.MIMEType)
}

func TestTemplateTypesStepsHandler(t *testing.T) {
	handler := TemplateTypesStepsHandler()
	ctx := context.Background()

	req := &mcp.ReadResourceRequest{
		Params: &mcp.ReadResourceParams{
			URI: "argo://docs/template-types/steps",
		},
	}

	result, err := handler(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Contents, 1)

	content := result.Contents[0]
	assert.Equal(t, "argo://docs/template-types/steps", content.URI)
	assert.NotEmpty(t, content.Text)
	assert.Contains(t, content.Text, "Steps Template Type")
	assert.Contains(t, content.Text, "sequential")
	assert.Contains(t, content.Text, "parallel")
}

// Suspend Template Tests

func TestTemplateTypesSuspendResource(t *testing.T) {
	resource := TemplateTypesSuspendResource()

	assert.Equal(t, "argo://docs/template-types/suspend", resource.URI)
	assert.Equal(t, "template-types-suspend", resource.Name)
	assert.Equal(t, "Suspend Template Type", resource.Title)
	assert.Equal(t, "text/markdown", resource.MIMEType)
}

func TestTemplateTypesSuspendHandler(t *testing.T) {
	handler := TemplateTypesSuspendHandler()
	ctx := context.Background()

	req := &mcp.ReadResourceRequest{
		Params: &mcp.ReadResourceParams{
			URI: "argo://docs/template-types/suspend",
		},
	}

	result, err := handler(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Contents, 1)

	content := result.Contents[0]
	assert.Equal(t, "argo://docs/template-types/suspend", content.URI)
	assert.NotEmpty(t, content.Text)
	assert.Contains(t, content.Text, "Suspend Template Type")
	assert.Contains(t, content.Text, "duration")
	assert.Contains(t, content.Text, "approval")
}

// Resource Template Tests

func TestTemplateTypesResourceResource(t *testing.T) {
	resource := TemplateTypesResourceResource()

	assert.Equal(t, "argo://docs/template-types/resource", resource.URI)
	assert.Equal(t, "template-types-resource", resource.Name)
	assert.Equal(t, "Resource Template Type", resource.Title)
	assert.Equal(t, "text/markdown", resource.MIMEType)
}

func TestTemplateTypesResourceHandler(t *testing.T) {
	handler := TemplateTypesResourceHandler()
	ctx := context.Background()

	req := &mcp.ReadResourceRequest{
		Params: &mcp.ReadResourceParams{
			URI: "argo://docs/template-types/resource",
		},
	}

	result, err := handler(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Contents, 1)

	content := result.Contents[0]
	assert.Equal(t, "argo://docs/template-types/resource", content.URI)
	assert.NotEmpty(t, content.Text)
	assert.Contains(t, content.Text, "Resource Template Type")
	assert.Contains(t, content.Text, "manifest")
	assert.Contains(t, content.Text, "action")
	assert.Contains(t, content.Text, "Kubernetes")
}

// HTTP Template Tests

func TestTemplateTypesHTTPResource(t *testing.T) {
	resource := TemplateTypesHTTPResource()

	assert.Equal(t, "argo://docs/template-types/http", resource.URI)
	assert.Equal(t, "template-types-http", resource.Name)
	assert.Equal(t, "HTTP Template Type", resource.Title)
	assert.Equal(t, "text/markdown", resource.MIMEType)
}

func TestTemplateTypesHTTPHandler(t *testing.T) {
	handler := TemplateTypesHTTPHandler()
	ctx := context.Background()

	req := &mcp.ReadResourceRequest{
		Params: &mcp.ReadResourceParams{
			URI: "argo://docs/template-types/http",
		},
	}

	result, err := handler(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Contents, 1)

	content := result.Contents[0]
	assert.Equal(t, "argo://docs/template-types/http", content.URI)
	assert.NotEmpty(t, content.Text)
	assert.Contains(t, content.Text, "HTTP Template Type")
	assert.Contains(t, content.Text, "url")
	assert.Contains(t, content.Text, "method")
	assert.Contains(t, content.Text, "HTTP requests")
}
