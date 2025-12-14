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

	// Should have exactly 4 resource registrars
	assert.Len(t, registrars, 4, "Expected 4 resource registrars")

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
