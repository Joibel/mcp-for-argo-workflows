//go:build e2e

package e2e

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Joibel/mcp-for-argo-workflows/internal/tools"
)

// TestWorkflowTemplate_CRUD tests the full CRUD lifecycle: create → get → list → delete.
func TestWorkflowTemplate_CRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	ctx := context.Background()
	cluster := SetupE2ECluster(ctx, t)

	// Load test workflow template
	manifest := LoadTestDataFile(t, "workflow-template.yaml")

	// Step 1: Create workflow template
	t.Log("Creating workflow template...")
	createHandler := tools.CreateWorkflowTemplateHandler(cluster.ArgoClient)
	createInput := tools.CreateWorkflowTemplateInput{
		Namespace: cluster.ArgoNamespace,
		Manifest:  manifest,
	}

	_, createOutput, err := createHandler(ctx, nil, createInput)
	require.NoError(t, err, "Failed to create workflow template")
	require.NotNil(t, createOutput)

	templateName := createOutput.Name
	t.Logf("Created workflow template: %s", templateName)

	// Verify template was created
	assert.True(t, cluster.WorkflowTemplateExists(t, cluster.ArgoNamespace, templateName),
		"WorkflowTemplate should exist after creation")

	// Cleanup at the end
	defer func() {
		deleteHandler := tools.DeleteWorkflowTemplateHandler(cluster.ArgoClient)
		deleteInput := tools.DeleteWorkflowTemplateInput{
			Namespace: cluster.ArgoNamespace,
			Name:      templateName,
		}
		_, _, _ = deleteHandler(ctx, nil, deleteInput)
	}()

	// Step 2: Get workflow template
	t.Log("Getting workflow template...")
	getHandler := tools.GetWorkflowTemplateHandler(cluster.ArgoClient)
	getInput := tools.GetWorkflowTemplateInput{
		Namespace: cluster.ArgoNamespace,
		Name:      templateName,
	}

	_, getOutput, err := getHandler(ctx, nil, getInput)
	require.NoError(t, err, "Failed to get workflow template")
	require.NotNil(t, getOutput)

	assert.Equal(t, templateName, getOutput.Name)
	assert.Equal(t, cluster.ArgoNamespace, getOutput.Namespace)
	assert.NotEmpty(t, getOutput.CreatedAt)
	assert.NotEmpty(t, getOutput.Templates, "Should have templates")

	// Step 3: List workflow templates
	t.Log("Listing workflow templates...")
	listHandler := tools.ListWorkflowTemplatesHandler(cluster.ArgoClient)
	listInput := tools.ListWorkflowTemplatesInput{
		Namespace: cluster.ArgoNamespace,
	}

	_, listOutput, err := listHandler(ctx, nil, listInput)
	require.NoError(t, err, "Failed to list workflow templates")
	require.NotNil(t, listOutput)

	// Verify our template is in the list
	assert.NotEmpty(t, listOutput.Templates, "Should have at least one template")

	found := false
	for _, tmpl := range listOutput.Templates {
		if tmpl.Name == templateName {
			found = true
			assert.Equal(t, cluster.ArgoNamespace, tmpl.Namespace)
			break
		}
	}
	assert.True(t, found, "Created template should be in the list")

	// Step 4: Delete workflow template
	t.Log("Deleting workflow template...")
	deleteHandler := tools.DeleteWorkflowTemplateHandler(cluster.ArgoClient)
	deleteInput := tools.DeleteWorkflowTemplateInput{
		Namespace: cluster.ArgoNamespace,
		Name:      templateName,
	}

	_, deleteOutput, err := deleteHandler(ctx, nil, deleteInput)
	require.NoError(t, err, "Failed to delete workflow template")
	require.NotNil(t, deleteOutput)

	assert.Equal(t, templateName, deleteOutput.Name)

	// Verify template was deleted (give it a moment to propagate)
	time.Sleep(2 * time.Second)
	assert.False(t, cluster.WorkflowTemplateExists(t, cluster.ArgoNamespace, templateName),
		"WorkflowTemplate should be deleted")
}

// TestWorkflowTemplate_SubmitWithRef tests creating a template and submitting a workflow that references it.
func TestWorkflowTemplate_SubmitWithRef(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	ctx := context.Background()
	cluster := SetupE2ECluster(ctx, t)

	// Step 1: Create workflow template
	t.Log("Creating workflow template...")
	templateManifest := LoadTestDataFile(t, "workflow-template.yaml")

	createHandler := tools.CreateWorkflowTemplateHandler(cluster.ArgoClient)
	createInput := tools.CreateWorkflowTemplateInput{
		Namespace: cluster.ArgoNamespace,
		Manifest:  templateManifest,
	}

	_, createOutput, err := createHandler(ctx, nil, createInput)
	require.NoError(t, err, "Failed to create workflow template")

	templateName := createOutput.Name
	t.Logf("Created workflow template: %s", templateName)

	// Cleanup template at the end
	defer func() {
		deleteHandler := tools.DeleteWorkflowTemplateHandler(cluster.ArgoClient)
		deleteInput := tools.DeleteWorkflowTemplateInput{
			Namespace: cluster.ArgoNamespace,
			Name:      templateName,
		}
		_, _, _ = deleteHandler(ctx, nil, deleteInput)
	}()

	// Step 2: Submit a workflow that references the template
	t.Log("Submitting workflow from template...")
	workflowManifest := `apiVersion: argoproj.io/v1alpha1
kind: Workflow
metadata:
  generateName: from-template-
spec:
  workflowTemplateRef:
    name: test-template
  arguments:
    parameters:
      - name: message
        value: "Hello from workflow using template"
`

	submitHandler := tools.SubmitWorkflowHandler(cluster.ArgoClient)
	submitInput := tools.SubmitWorkflowInput{
		Namespace: cluster.ArgoNamespace,
		Manifest:  workflowManifest,
	}

	_, submitOutput, err := submitHandler(ctx, nil, submitInput)
	require.NoError(t, err, "Failed to submit workflow from template")

	workflowName := submitOutput.Name
	t.Logf("Submitted workflow: %s", workflowName)

	// Cleanup workflow at the end
	defer func() {
		deleteWorkflowHandler := tools.DeleteWorkflowHandler(cluster.ArgoClient)
		deleteWorkflowInput := tools.DeleteWorkflowInput{
			Namespace: cluster.ArgoNamespace,
			Name:      workflowName,
		}
		_, _, _ = deleteWorkflowHandler(ctx, nil, deleteWorkflowInput)
	}()

	// Step 3: Wait for workflow to complete
	t.Log("Waiting for workflow to complete...")
	finalPhase := cluster.WaitForWorkflowPhase(t, cluster.ArgoNamespace, workflowName,
		2*time.Minute, "Succeeded", "Failed", "Error")

	assert.Equal(t, "Succeeded", finalPhase, "Workflow should complete successfully")

	// Step 4: Verify logs contain the custom message
	t.Log("Verifying workflow output...")
	logsHandler := tools.LogsWorkflowHandler(cluster.ArgoClient)
	logsInput := tools.LogsWorkflowInput{
		Namespace: cluster.ArgoNamespace,
		Name:      workflowName,
	}

	_, logsOutput, err := logsHandler(ctx, nil, logsInput)
	require.NoError(t, err, "Failed to get workflow logs")
	require.NotNil(t, logsOutput)

	assert.Contains(t, logsOutput.Logs, "Hello from workflow using template",
		"Logs should contain the custom message parameter")
}

// TestWorkflowTemplate_GetConsistency tests that getting a template returns consistent data.
func TestWorkflowTemplate_GetConsistency(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	ctx := context.Background()
	cluster := SetupE2ECluster(ctx, t)

	// Create initial template
	t.Log("Creating workflow template...")
	manifest := LoadTestDataFile(t, "workflow-template.yaml")

	createHandler := tools.CreateWorkflowTemplateHandler(cluster.ArgoClient)
	createInput := tools.CreateWorkflowTemplateInput{
		Namespace: cluster.ArgoNamespace,
		Manifest:  manifest,
	}

	_, createOutput, err := createHandler(ctx, nil, createInput)
	require.NoError(t, err, "Failed to create workflow template")

	templateName := createOutput.Name
	t.Logf("Created workflow template: %s", templateName)

	// Cleanup at the end
	defer func() {
		deleteHandler := tools.DeleteWorkflowTemplateHandler(cluster.ArgoClient)
		deleteInput := tools.DeleteWorkflowTemplateInput{
			Namespace: cluster.ArgoNamespace,
			Name:      templateName,
		}
		_, _, _ = deleteHandler(ctx, nil, deleteInput)
	}()

	// Get the template
	getHandler := tools.GetWorkflowTemplateHandler(cluster.ArgoClient)
	getInput := tools.GetWorkflowTemplateInput{
		Namespace: cluster.ArgoNamespace,
		Name:      templateName,
	}

	_, getOutput1, err := getHandler(ctx, nil, getInput)
	require.NoError(t, err, "Failed to get workflow template")

	originalCreatedAt := getOutput1.CreatedAt

	// Wait a moment to ensure timestamps would differ
	time.Sleep(1 * time.Second)

	// Note: Update is typically done via kubectl apply or direct API calls
	// For now, we just verify the template exists and is stable
	_, getOutput2, err := getHandler(ctx, nil, getInput)
	require.NoError(t, err, "Failed to get workflow template again")

	// Verify the template is consistent
	assert.Equal(t, originalCreatedAt, getOutput2.CreatedAt,
		"CreatedAt should remain the same")
	assert.Equal(t, templateName, getOutput2.Name)
}
