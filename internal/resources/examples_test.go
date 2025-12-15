package resources

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Hello World Example Tests

func TestExamplesHelloWorldResource(t *testing.T) {
	resource := ExamplesHelloWorldResource()

	assert.Equal(t, "argo://examples/hello-world", resource.URI)
	assert.Equal(t, "examples-hello-world", resource.Name)
	assert.Equal(t, "Hello World Workflow Example", resource.Title)
	assert.Equal(t, "text/markdown", resource.MIMEType)
	assert.NotEmpty(t, resource.Description)
}

func TestExamplesHelloWorldHandler(t *testing.T) {
	handler := ExamplesHelloWorldHandler()
	ctx := context.Background()

	req := &mcp.ReadResourceRequest{
		Params: &mcp.ReadResourceParams{
			URI: "argo://examples/hello-world",
		},
	}

	result, err := handler(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Contents, 1)

	content := result.Contents[0]
	assert.Equal(t, "argo://examples/hello-world", content.URI)
	assert.NotEmpty(t, content.Text)
	assert.Contains(t, content.Text, "Hello World")
	assert.Contains(t, content.Text, "apiVersion")
	assert.Contains(t, content.Text, "kind: Workflow")
}

// Multi-Step Example Tests

func TestExamplesMultiStepResource(t *testing.T) {
	resource := ExamplesMultiStepResource()

	assert.Equal(t, "argo://examples/multi-step", resource.URI)
	assert.Equal(t, "examples-multi-step", resource.Name)
	assert.Equal(t, "Multi-Step Workflow Example", resource.Title)
	assert.Equal(t, "text/markdown", resource.MIMEType)
}

func TestExamplesMultiStepHandler(t *testing.T) {
	handler := ExamplesMultiStepHandler()
	ctx := context.Background()

	req := &mcp.ReadResourceRequest{
		Params: &mcp.ReadResourceParams{
			URI: "argo://examples/multi-step",
		},
	}

	result, err := handler(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Contents, 1)

	content := result.Contents[0]
	assert.Contains(t, content.Text, "Multi-Step")
	assert.Contains(t, content.Text, "steps")
	assert.Contains(t, content.Text, "outputs")
}

// DAG Diamond Example Tests

func TestExamplesDAGDiamondResource(t *testing.T) {
	resource := ExamplesDAGDiamondResource()

	assert.Equal(t, "argo://examples/dag-diamond", resource.URI)
	assert.Equal(t, "examples-dag-diamond", resource.Name)
	assert.Equal(t, "DAG Diamond Pattern Example", resource.Title)
}

func TestExamplesDAGDiamondHandler(t *testing.T) {
	handler := ExamplesDAGDiamondHandler()
	ctx := context.Background()

	req := &mcp.ReadResourceRequest{
		Params: &mcp.ReadResourceParams{
			URI: "argo://examples/dag-diamond",
		},
	}

	result, err := handler(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, result)

	content := result.Contents[0]
	assert.Contains(t, content.Text, "diamond")
	assert.Contains(t, content.Text, "dag")
	assert.Contains(t, content.Text, "dependencies")
}

// Parameters Example Tests

func TestExamplesParametersResource(t *testing.T) {
	resource := ExamplesParametersResource()

	assert.Equal(t, "argo://examples/parameters", resource.URI)
	assert.Equal(t, "examples-parameters", resource.Name)
	assert.Equal(t, "Parameters Example", resource.Title)
}

func TestExamplesParametersHandler(t *testing.T) {
	handler := ExamplesParametersHandler()
	ctx := context.Background()

	req := &mcp.ReadResourceRequest{
		Params: &mcp.ReadResourceParams{
			URI: "argo://examples/parameters",
		},
	}

	result, err := handler(ctx, req)
	require.NoError(t, err)

	content := result.Contents[0]
	assert.Contains(t, content.Text, "parameters")
	assert.Contains(t, content.Text, "arguments")
	assert.Contains(t, content.Text, "inputs")
}

// Artifacts Example Tests

func TestExamplesArtifactsResource(t *testing.T) {
	resource := ExamplesArtifactsResource()

	assert.Equal(t, "argo://examples/artifacts", resource.URI)
	assert.Equal(t, "examples-artifacts", resource.Name)
}

func TestExamplesArtifactsHandler(t *testing.T) {
	handler := ExamplesArtifactsHandler()
	ctx := context.Background()

	req := &mcp.ReadResourceRequest{
		Params: &mcp.ReadResourceParams{
			URI: "argo://examples/artifacts",
		},
	}

	result, err := handler(ctx, req)
	require.NoError(t, err)

	content := result.Contents[0]
	assert.Contains(t, content.Text, "artifacts")
	assert.Contains(t, content.Text, "S3")
	assert.Contains(t, content.Text, "path")
}

// Loops Example Tests

func TestExamplesLoopsResource(t *testing.T) {
	resource := ExamplesLoopsResource()

	assert.Equal(t, "argo://examples/loops", resource.URI)
	assert.Equal(t, "examples-loops", resource.Name)
}

func TestExamplesLoopsHandler(t *testing.T) {
	handler := ExamplesLoopsHandler()
	ctx := context.Background()

	req := &mcp.ReadResourceRequest{
		Params: &mcp.ReadResourceParams{
			URI: "argo://examples/loops",
		},
	}

	result, err := handler(ctx, req)
	require.NoError(t, err)

	content := result.Contents[0]
	assert.Contains(t, content.Text, "withItems")
	assert.Contains(t, content.Text, "withParam")
	assert.Contains(t, content.Text, "withSequence")
}

// Conditionals Example Tests

func TestExamplesConditionalsResource(t *testing.T) {
	resource := ExamplesConditionalsResource()

	assert.Equal(t, "argo://examples/conditionals", resource.URI)
	assert.Equal(t, "examples-conditionals", resource.Name)
}

func TestExamplesConditionalsHandler(t *testing.T) {
	handler := ExamplesConditionalsHandler()
	ctx := context.Background()

	req := &mcp.ReadResourceRequest{
		Params: &mcp.ReadResourceParams{
			URI: "argo://examples/conditionals",
		},
	}

	result, err := handler(ctx, req)
	require.NoError(t, err)

	content := result.Contents[0]
	assert.Contains(t, content.Text, "when")
	assert.Contains(t, content.Text, "conditional")
}

// Retries Example Tests

func TestExamplesRetriesResource(t *testing.T) {
	resource := ExamplesRetriesResource()

	assert.Equal(t, "argo://examples/retries", resource.URI)
	assert.Equal(t, "examples-retries", resource.Name)
}

func TestExamplesRetriesHandler(t *testing.T) {
	handler := ExamplesRetriesHandler()
	ctx := context.Background()

	req := &mcp.ReadResourceRequest{
		Params: &mcp.ReadResourceParams{
			URI: "argo://examples/retries",
		},
	}

	result, err := handler(ctx, req)
	require.NoError(t, err)

	content := result.Contents[0]
	assert.Contains(t, content.Text, "retryStrategy")
	assert.Contains(t, content.Text, "backoff")
	assert.Contains(t, content.Text, "limit")
}

// Timeout Limits Example Tests

func TestExamplesTimeoutLimitsResource(t *testing.T) {
	resource := ExamplesTimeoutLimitsResource()

	assert.Equal(t, "argo://examples/timeout-limits", resource.URI)
	assert.Equal(t, "examples-timeout-limits", resource.Name)
}

func TestExamplesTimeoutLimitsHandler(t *testing.T) {
	handler := ExamplesTimeoutLimitsHandler()
	ctx := context.Background()

	req := &mcp.ReadResourceRequest{
		Params: &mcp.ReadResourceParams{
			URI: "argo://examples/timeout-limits",
		},
	}

	result, err := handler(ctx, req)
	require.NoError(t, err)

	content := result.Contents[0]
	assert.Contains(t, content.Text, "activeDeadlineSeconds")
	assert.Contains(t, content.Text, "timeout")
}

// Resource Management Example Tests

func TestExamplesResourceManagementResource(t *testing.T) {
	resource := ExamplesResourceManagementResource()

	assert.Equal(t, "argo://examples/resource-management", resource.URI)
	assert.Equal(t, "examples-resource-management", resource.Name)
}

func TestExamplesResourceManagementHandler(t *testing.T) {
	handler := ExamplesResourceManagementHandler()
	ctx := context.Background()

	req := &mcp.ReadResourceRequest{
		Params: &mcp.ReadResourceParams{
			URI: "argo://examples/resource-management",
		},
	}

	result, err := handler(ctx, req)
	require.NoError(t, err)

	content := result.Contents[0]
	assert.Contains(t, content.Text, "resources")
	assert.Contains(t, content.Text, "requests")
	assert.Contains(t, content.Text, "limits")
	assert.Contains(t, content.Text, "cpu")
	assert.Contains(t, content.Text, "memory")
}

// Volumes Example Tests

func TestExamplesVolumesResource(t *testing.T) {
	resource := ExamplesVolumesResource()

	assert.Equal(t, "argo://examples/volumes", resource.URI)
	assert.Equal(t, "examples-volumes", resource.Name)
}

func TestExamplesVolumesHandler(t *testing.T) {
	handler := ExamplesVolumesHandler()
	ctx := context.Background()

	req := &mcp.ReadResourceRequest{
		Params: &mcp.ReadResourceParams{
			URI: "argo://examples/volumes",
		},
	}

	result, err := handler(ctx, req)
	require.NoError(t, err)

	content := result.Contents[0]
	assert.Contains(t, content.Text, "volumes")
	assert.Contains(t, content.Text, "volumeMounts")
	assert.Contains(t, content.Text, "PersistentVolumeClaim")
}

// Exit Handlers Example Tests

func TestExamplesExitHandlersResource(t *testing.T) {
	resource := ExamplesExitHandlersResource()

	assert.Equal(t, "argo://examples/exit-handlers", resource.URI)
	assert.Equal(t, "examples-exit-handlers", resource.Name)
}

func TestExamplesExitHandlersHandler(t *testing.T) {
	handler := ExamplesExitHandlersHandler()
	ctx := context.Background()

	req := &mcp.ReadResourceRequest{
		Params: &mcp.ReadResourceParams{
			URI: "argo://examples/exit-handlers",
		},
	}

	result, err := handler(ctx, req)
	require.NoError(t, err)

	content := result.Contents[0]
	assert.Contains(t, content.Text, "onExit")
	assert.Contains(t, content.Text, "exit")
	assert.Contains(t, content.Text, "workflow.status")
}
