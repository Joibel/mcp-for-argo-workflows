//go:build e2e

package e2e

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/argoproj/argo-workflows/v3/pkg/apiclient/workflow"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Joibel/mcp-for-argo-workflows/internal/tools"
)

// TestWorkflow_FullLifecycle tests the full lifecycle: submit → get → logs → wait → delete.
func TestWorkflow_FullLifecycle(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	ctx := context.Background()
	cluster := SetupE2ECluster(ctx, t)

	// Load test workflow
	manifest := LoadTestDataFile(t, "hello-world.yaml")

	// Step 1: Submit workflow
	t.Log("Submitting workflow...")
	submitHandler := tools.SubmitWorkflowHandler(cluster.ArgoClient)
	submitInput := tools.SubmitWorkflowInput{
		Namespace: cluster.ArgoNamespace,
		Manifest:  manifest,
	}

	_, submitOutput, err := submitHandler(ctx, nil, submitInput)
	require.NoError(t, err, "Failed to submit workflow")
	require.NotNil(t, submitOutput)

	workflowName := submitOutput.Name
	t.Logf("Submitted workflow: %s", workflowName)

	// Cleanup at the end (also verified explicitly below)
	defer func() {
		deleteHandler := tools.DeleteWorkflowHandler(cluster.ArgoClient)
		deleteInput := tools.DeleteWorkflowInput{
			Namespace: cluster.ArgoNamespace,
			Name:      workflowName,
		}
		_, _, _ = deleteHandler(ctx, nil, deleteInput)
	}()

	// Verify workflow was created
	assert.True(t, cluster.WorkflowExists(t, cluster.ArgoNamespace, workflowName),
		"Workflow should exist after submission")

	// Step 2: Get workflow details
	t.Log("Getting workflow details...")
	getHandler := tools.GetWorkflowHandler(cluster.ArgoClient)
	getInput := tools.GetWorkflowInput{
		Namespace: cluster.ArgoNamespace,
		Name:      workflowName,
	}

	_, getOutput, err := getHandler(ctx, nil, getInput)
	require.NoError(t, err, "Failed to get workflow")
	require.NotNil(t, getOutput)

	assert.Equal(t, workflowName, getOutput.Name)
	assert.Equal(t, cluster.ArgoNamespace, getOutput.Namespace)
	assert.NotEmpty(t, getOutput.UID)
	assert.NotEmpty(t, getOutput.Phase)

	// Step 3: Wait for workflow to complete
	t.Log("Waiting for workflow to complete...")
	finalPhase := cluster.WaitForWorkflowPhase(t, cluster.ArgoNamespace, workflowName,
		2*time.Minute, "Succeeded", "Failed", "Error")

	assert.Equal(t, "Succeeded", finalPhase, "Workflow should complete successfully")

	// Step 4: Get logs (verify logs are accessible)
	t.Log("Getting workflow logs...")
	logsHandler := tools.LogsWorkflowHandler(cluster.ArgoClient)
	logsInput := tools.LogsWorkflowInput{
		Namespace: cluster.ArgoNamespace,
		Name:      workflowName,
	}

	_, logsOutput, err := logsHandler(ctx, nil, logsInput)
	require.NoError(t, err, "Failed to get workflow logs")
	require.NotNil(t, logsOutput)
	assert.NotEmpty(t, logsOutput.Logs, "Logs should not be empty")
	// Check that at least one log entry contains the expected output
	foundHelloWorld := false
	for _, entry := range logsOutput.Logs {
		if strings.Contains(entry.Content, "Hello World") {
			foundHelloWorld = true
			break
		}
	}
	assert.True(t, foundHelloWorld, "Logs should contain expected output 'Hello World'")

	// Step 5: Delete workflow
	t.Log("Deleting workflow...")
	deleteHandler := tools.DeleteWorkflowHandler(cluster.ArgoClient)
	deleteInput := tools.DeleteWorkflowInput{
		Namespace: cluster.ArgoNamespace,
		Name:      workflowName,
	}

	_, deleteOutput, err := deleteHandler(ctx, nil, deleteInput)
	require.NoError(t, err, "Failed to delete workflow")
	require.NotNil(t, deleteOutput)

	assert.Equal(t, workflowName, deleteOutput.Name)

	// Verify workflow was deleted (give it a moment to propagate)
	time.Sleep(2 * time.Second)
	assert.False(t, cluster.WorkflowExists(t, cluster.ArgoNamespace, workflowName),
		"Workflow should be deleted")
}

// TestWorkflow_SuspendResume tests suspend → resume workflow operations.
func TestWorkflow_SuspendResume(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	ctx := context.Background()
	cluster := SetupE2ECluster(ctx, t)

	// Load DAG workflow (longer running, easier to suspend)
	manifest := LoadTestDataFile(t, "dag-workflow.yaml")

	// Submit workflow
	t.Log("Submitting DAG workflow...")
	submitHandler := tools.SubmitWorkflowHandler(cluster.ArgoClient)
	submitInput := tools.SubmitWorkflowInput{
		Namespace: cluster.ArgoNamespace,
		Manifest:  manifest,
	}

	_, submitOutput, err := submitHandler(ctx, nil, submitInput)
	require.NoError(t, err, "Failed to submit workflow")

	workflowName := submitOutput.Name
	t.Logf("Submitted workflow: %s", workflowName)

	// Cleanup at the end
	defer func() {
		deleteHandler := tools.DeleteWorkflowHandler(cluster.ArgoClient)
		deleteInput := tools.DeleteWorkflowInput{
			Namespace: cluster.ArgoNamespace,
			Name:      workflowName,
		}
		_, _, _ = deleteHandler(ctx, nil, deleteInput)
	}()

	// Wait a moment for workflow to start
	time.Sleep(2 * time.Second)

	// Suspend workflow
	t.Log("Suspending workflow...")
	suspendHandler := tools.SuspendWorkflowHandler(cluster.ArgoClient)
	suspendInput := tools.SuspendWorkflowInput{
		Namespace: cluster.ArgoNamespace,
		Name:      workflowName,
	}

	_, suspendOutput, err := suspendHandler(ctx, nil, suspendInput)
	require.NoError(t, err, "Failed to suspend workflow")
	require.NotNil(t, suspendOutput)

	t.Logf("Workflow suspended, phase: %s", suspendOutput.Phase)

	// Get workflow to verify it's suspended
	wfService := cluster.ArgoClient.WorkflowService()
	wf, err := wfService.GetWorkflow(cluster.ArgoClient.Context(), &workflow.WorkflowGetRequest{
		Namespace: cluster.ArgoNamespace,
		Name:      workflowName,
	})
	require.NoError(t, err, "Failed to get workflow after suspend")
	assert.NotNil(t, wf.Spec.Suspend, "Workflow spec should have suspend set")
	assert.True(t, *wf.Spec.Suspend, "Workflow should be suspended")

	// Resume workflow
	t.Log("Resuming workflow...")
	resumeHandler := tools.ResumeWorkflowHandler(cluster.ArgoClient)
	resumeInput := tools.ResumeWorkflowInput{
		Namespace: cluster.ArgoNamespace,
		Name:      workflowName,
	}

	_, resumeOutput, err := resumeHandler(ctx, nil, resumeInput)
	require.NoError(t, err, "Failed to resume workflow")
	require.NotNil(t, resumeOutput)

	t.Logf("Workflow resumed, phase: %s", resumeOutput.Phase)

	// Wait for workflow to complete
	t.Log("Waiting for workflow to complete...")
	finalPhase := cluster.WaitForWorkflowPhase(t, cluster.ArgoNamespace, workflowName,
		2*time.Minute, "Succeeded", "Failed", "Error")

	assert.Equal(t, "Succeeded", finalPhase, "Workflow should complete successfully after resume")
}

// TestWorkflow_Lint tests linting valid and invalid manifests.
func TestWorkflow_Lint(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	ctx := context.Background()
	cluster := SetupE2ECluster(ctx, t)

	lintHandler := tools.LintWorkflowHandler(cluster.ArgoClient)

	//nolint:govet // Field alignment is not critical for test structs
	tests := []struct {
		name        string
		wantErr     bool
		errContains string
		manifest    string
	}{
		{
			name:     "valid hello-world workflow",
			manifest: LoadTestDataFile(t, "hello-world.yaml"),
			wantErr:  false,
		},
		{
			name:     "valid DAG workflow",
			manifest: LoadTestDataFile(t, "dag-workflow.yaml"),
			wantErr:  false,
		},
		{
			name: "invalid workflow - missing entrypoint",
			manifest: `apiVersion: argoproj.io/v1alpha1
kind: Workflow
metadata:
  generateName: invalid-
spec:
  templates:
    - name: main
      container:
        image: busybox:1.35
        command: [echo]
        args: ["hello"]
`,
			wantErr:     true,
			errContains: "entrypoint",
		},
		{
			name: "invalid workflow - bad image",
			manifest: `apiVersion: argoproj.io/v1alpha1
kind: Workflow
metadata:
  generateName: invalid-
spec:
  entrypoint: main
  templates:
    - name: main
      container:
        image: ""
        command: [echo]
`,
			wantErr:     true,
			errContains: "image",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lintInput := tools.LintWorkflowInput{
				Namespace: cluster.ArgoNamespace,
				Manifest:  tt.manifest,
			}

			_, lintOutput, err := lintHandler(ctx, nil, lintInput)

			if tt.wantErr {
				// Lint should return an output with validation errors, not a Go error
				require.NoError(t, err, "Lint handler should not return Go error")
				require.NotNil(t, lintOutput)
				assert.False(t, lintOutput.Valid, "Manifest should be invalid")
				assert.NotEmpty(t, lintOutput.Errors, "Should have validation errors")

				if tt.errContains != "" {
					found := false
					for _, errMsg := range lintOutput.Errors {
						if strings.Contains(errMsg, tt.errContains) {
							found = true
							break
						}
					}
					assert.True(t, found, "Expected error containing %q, got: %v", tt.errContains, lintOutput.Errors)
				}
			} else {
				require.NoError(t, err, "Lint handler should not return error")
				require.NotNil(t, lintOutput)
				assert.True(t, lintOutput.Valid, "Manifest should be valid")
				assert.Empty(t, lintOutput.Errors, "Should have no validation errors")
			}
		})
	}
}

// TestWorkflow_Submit tests the submit_workflow tool handler.
func TestWorkflow_Submit(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	ctx := context.Background()
	cluster := SetupE2ECluster(ctx, t)

	// Load test workflow
	manifest := LoadTestDataFile(t, "hello-world.yaml")

	// Submit workflow using the tool handler
	t.Log("Testing submit_workflow tool...")
	submitHandler := tools.SubmitWorkflowHandler(cluster.ArgoClient)
	submitInput := tools.SubmitWorkflowInput{
		Namespace: cluster.ArgoNamespace,
		Manifest:  manifest,
	}

	result, submitOutput, err := submitHandler(ctx, nil, submitInput)
	require.NoError(t, err, "submit_workflow should not return error")
	require.NotNil(t, result)
	require.NotNil(t, submitOutput)

	workflowName := submitOutput.Name

	// Cleanup at the end
	defer func() {
		deleteHandler := tools.DeleteWorkflowHandler(cluster.ArgoClient)
		deleteInput := tools.DeleteWorkflowInput{
			Namespace: cluster.ArgoNamespace,
			Name:      workflowName,
		}
		_, _, _ = deleteHandler(ctx, nil, deleteInput)
	}()

	// Verify submit output fields
	assert.NotEmpty(t, submitOutput.Name, "Name should be set")
	assert.Equal(t, cluster.ArgoNamespace, submitOutput.Namespace, "Namespace should match")
	assert.NotEmpty(t, submitOutput.UID, "UID should be set")
	assert.NotEmpty(t, submitOutput.Phase, "Phase should be set")
	assert.NotEmpty(t, submitOutput.CreatedAt, "CreatedAt should be set")

	// Verify workflow actually exists and is running
	assert.True(t, cluster.WorkflowExists(t, cluster.ArgoNamespace, workflowName),
		"Workflow should exist after submission")

	// Verify it starts running (phase should be Pending or Running initially)
	assert.Contains(t, []string{"Pending", "Running"}, submitOutput.Phase,
		"Workflow should start in Pending or Running phase")
}

// TestWorkflow_Get tests the get_workflow tool handler.
func TestWorkflow_Get(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	ctx := context.Background()
	cluster := SetupE2ECluster(ctx, t)

	// Submit a workflow first
	manifest := LoadTestDataFile(t, "hello-world.yaml")
	submitHandler := tools.SubmitWorkflowHandler(cluster.ArgoClient)
	submitInput := tools.SubmitWorkflowInput{
		Namespace: cluster.ArgoNamespace,
		Manifest:  manifest,
	}

	_, submitOutput, err := submitHandler(ctx, nil, submitInput)
	require.NoError(t, err, "Failed to submit workflow")

	workflowName := submitOutput.Name

	// Cleanup at the end
	defer func() {
		deleteHandler := tools.DeleteWorkflowHandler(cluster.ArgoClient)
		deleteInput := tools.DeleteWorkflowInput{
			Namespace: cluster.ArgoNamespace,
			Name:      workflowName,
		}
		_, _, _ = deleteHandler(ctx, nil, deleteInput)
	}()

	// Test get_workflow tool handler
	t.Log("Testing get_workflow tool...")
	getHandler := tools.GetWorkflowHandler(cluster.ArgoClient)
	getInput := tools.GetWorkflowInput{
		Namespace: cluster.ArgoNamespace,
		Name:      workflowName,
	}

	result, getOutput, err := getHandler(ctx, nil, getInput)
	require.NoError(t, err, "get_workflow should not return error")
	require.NotNil(t, result)
	require.NotNil(t, getOutput)

	// Verify all expected fields are present
	assert.Equal(t, workflowName, getOutput.Name, "Name should match")
	assert.Equal(t, cluster.ArgoNamespace, getOutput.Namespace, "Namespace should match")
	assert.NotEmpty(t, getOutput.UID, "UID should be set")
	assert.NotEmpty(t, getOutput.Phase, "Phase should be set")
	assert.NotEmpty(t, getOutput.CreatedAt, "CreatedAt should be set")
	assert.NotEmpty(t, getOutput.Entrypoint, "Entrypoint should be set")

	// Wait for completion and verify final state
	cluster.WaitForWorkflowPhase(t, cluster.ArgoNamespace, workflowName,
		2*time.Minute, "Succeeded", "Failed", "Error")

	// Get again after completion
	result, getOutput, err = getHandler(ctx, nil, getInput)
	require.NoError(t, err, "get_workflow should not return error after completion")
	require.NotNil(t, getOutput)

	assert.Equal(t, "Succeeded", getOutput.Phase, "Phase should be Succeeded")
	assert.NotEmpty(t, getOutput.StartedAt, "StartedAt should be set after completion")
	assert.NotEmpty(t, getOutput.FinishedAt, "FinishedAt should be set after completion")
}

// TestWorkflow_Delete tests the delete_workflow tool handler.
func TestWorkflow_Delete(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	ctx := context.Background()
	cluster := SetupE2ECluster(ctx, t)

	// Submit a workflow first
	manifest := LoadTestDataFile(t, "hello-world.yaml")
	submitHandler := tools.SubmitWorkflowHandler(cluster.ArgoClient)
	submitInput := tools.SubmitWorkflowInput{
		Namespace: cluster.ArgoNamespace,
		Manifest:  manifest,
	}

	_, submitOutput, err := submitHandler(ctx, nil, submitInput)
	require.NoError(t, err, "Failed to submit workflow")

	workflowName := submitOutput.Name

	// Verify workflow exists
	assert.True(t, cluster.WorkflowExists(t, cluster.ArgoNamespace, workflowName),
		"Workflow should exist after submission")

	// Test delete_workflow tool handler
	t.Log("Testing delete_workflow tool...")
	deleteHandler := tools.DeleteWorkflowHandler(cluster.ArgoClient)
	deleteInput := tools.DeleteWorkflowInput{
		Namespace: cluster.ArgoNamespace,
		Name:      workflowName,
	}

	result, deleteOutput, err := deleteHandler(ctx, nil, deleteInput)
	require.NoError(t, err, "delete_workflow should not return error")
	require.NotNil(t, result)
	require.NotNil(t, deleteOutput)

	// Verify delete output
	assert.Equal(t, workflowName, deleteOutput.Name, "Deleted workflow name should match")
	assert.Equal(t, cluster.ArgoNamespace, deleteOutput.Namespace, "Namespace should match")
	assert.True(t, deleteOutput.Deleted, "Deleted flag should be true")

	// Give deletion time to propagate
	time.Sleep(2 * time.Second)

	// Verify workflow is removed from list
	assert.False(t, cluster.WorkflowExists(t, cluster.ArgoNamespace, workflowName),
		"Workflow should not exist after deletion")

	// Verify it's not in list_workflows output
	listHandler := tools.ListWorkflowsHandler(cluster.ArgoClient)
	namespace := cluster.ArgoNamespace
	listInput := tools.ListWorkflowsInput{
		Namespace: &namespace,
	}

	_, listOutput, err := listHandler(ctx, nil, listInput)
	require.NoError(t, err, "list_workflows should not return error")

	for _, wf := range listOutput.Workflows {
		assert.NotEqual(t, workflowName, wf.Name, "Deleted workflow should not appear in list")
	}
}

// TestWorkflow_Logs tests the logs_workflow tool handler.
func TestWorkflow_Logs(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	ctx := context.Background()
	cluster := SetupE2ECluster(ctx, t)

	// Submit workflow
	manifest := LoadTestDataFile(t, "hello-world.yaml")
	submitHandler := tools.SubmitWorkflowHandler(cluster.ArgoClient)
	submitInput := tools.SubmitWorkflowInput{
		Namespace: cluster.ArgoNamespace,
		Manifest:  manifest,
	}

	_, submitOutput, err := submitHandler(ctx, nil, submitInput)
	require.NoError(t, err, "Failed to submit workflow")

	workflowName := submitOutput.Name

	// Cleanup at the end
	defer func() {
		deleteHandler := tools.DeleteWorkflowHandler(cluster.ArgoClient)
		deleteInput := tools.DeleteWorkflowInput{
			Namespace: cluster.ArgoNamespace,
			Name:      workflowName,
		}
		_, _, _ = deleteHandler(ctx, nil, deleteInput)
	}()

	// Wait for workflow to complete (so logs are available)
	finalPhase := cluster.WaitForWorkflowPhase(t, cluster.ArgoNamespace, workflowName,
		2*time.Minute, "Succeeded", "Failed", "Error")
	require.Equal(t, "Succeeded", finalPhase, "Workflow should succeed")

	// Test logs_workflow tool handler
	t.Log("Testing logs_workflow tool...")
	logsHandler := tools.LogsWorkflowHandler(cluster.ArgoClient)
	logsInput := tools.LogsWorkflowInput{
		Namespace: cluster.ArgoNamespace,
		Name:      workflowName,
	}

	result, logsOutput, err := logsHandler(ctx, nil, logsInput)
	require.NoError(t, err, "logs_workflow should not return error")
	require.NotNil(t, result)
	require.NotNil(t, logsOutput)

	// Verify logs output
	assert.NotEmpty(t, logsOutput.Logs, "Logs should not be empty")

	// Verify log entries have expected structure
	for _, entry := range logsOutput.Logs {
		assert.NotEmpty(t, entry.PodName, "Log entry should have PodName")
		assert.NotEmpty(t, entry.Content, "Log entry should have Content")
	}

	// Check that logs contain expected output from hello-world workflow
	foundHelloWorld := false
	for _, entry := range logsOutput.Logs {
		if strings.Contains(entry.Content, "Hello World") {
			foundHelloWorld = true
			break
		}
	}
	assert.True(t, foundHelloWorld, "Logs should contain 'Hello World' output")
}

// TestWorkflow_Logs_WithGrep tests logs_workflow with grep filtering.
func TestWorkflow_Logs_WithGrep(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	ctx := context.Background()
	cluster := SetupE2ECluster(ctx, t)

	// Submit workflow
	manifest := LoadTestDataFile(t, "hello-world.yaml")
	submitHandler := tools.SubmitWorkflowHandler(cluster.ArgoClient)
	submitInput := tools.SubmitWorkflowInput{
		Namespace: cluster.ArgoNamespace,
		Manifest:  manifest,
	}

	_, submitOutput, err := submitHandler(ctx, nil, submitInput)
	require.NoError(t, err, "Failed to submit workflow")

	workflowName := submitOutput.Name

	// Cleanup at the end
	defer func() {
		deleteHandler := tools.DeleteWorkflowHandler(cluster.ArgoClient)
		deleteInput := tools.DeleteWorkflowInput{
			Namespace: cluster.ArgoNamespace,
			Name:      workflowName,
		}
		_, _, _ = deleteHandler(ctx, nil, deleteInput)
	}()

	// Wait for workflow to complete
	cluster.WaitForWorkflowPhase(t, cluster.ArgoNamespace, workflowName,
		2*time.Minute, "Succeeded", "Failed", "Error")

	// Test logs_workflow with grep filter
	t.Log("Testing logs_workflow with grep filter...")
	logsHandler := tools.LogsWorkflowHandler(cluster.ArgoClient)
	grepPattern := "Hello"
	logsInput := tools.LogsWorkflowInput{
		Namespace: cluster.ArgoNamespace,
		Name:      workflowName,
		Grep:      &grepPattern,
	}

	result, logsOutput, err := logsHandler(ctx, nil, logsInput)
	require.NoError(t, err, "logs_workflow with grep should not return error")
	require.NotNil(t, result)
	require.NotNil(t, logsOutput)

	// Verify grep returned at least some results
	assert.NotEmpty(t, logsOutput.Logs, "Grep filter should have matched at least one log entry")

	// Verify filtered logs contain the grep pattern
	for _, entry := range logsOutput.Logs {
		assert.Contains(t, entry.Content, "Hello",
			"Filtered log entries should contain grep pattern")
	}
}

// TestWorkflow_List tests listing workflows.
func TestWorkflow_List(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	ctx := context.Background()
	cluster := SetupE2ECluster(ctx, t)

	// Submit a workflow first
	manifest := LoadTestDataFile(t, "hello-world.yaml")
	submitHandler := tools.SubmitWorkflowHandler(cluster.ArgoClient)
	submitInput := tools.SubmitWorkflowInput{
		Namespace: cluster.ArgoNamespace,
		Manifest:  manifest,
	}

	_, submitOutput, err := submitHandler(ctx, nil, submitInput)
	require.NoError(t, err, "Failed to submit workflow")

	workflowName := submitOutput.Name
	defer func() {
		deleteHandler := tools.DeleteWorkflowHandler(cluster.ArgoClient)
		deleteInput := tools.DeleteWorkflowInput{
			Namespace: cluster.ArgoNamespace,
			Name:      workflowName,
		}
		_, _, _ = deleteHandler(ctx, nil, deleteInput)
	}()

	// List workflows
	t.Log("Listing workflows...")
	listHandler := tools.ListWorkflowsHandler(cluster.ArgoClient)
	namespace := cluster.ArgoNamespace
	listInput := tools.ListWorkflowsInput{
		Namespace: &namespace,
	}

	_, listOutput, err := listHandler(ctx, nil, listInput)
	require.NoError(t, err, "Failed to list workflows")
	require.NotNil(t, listOutput)

	// Verify our workflow is in the list
	assert.NotEmpty(t, listOutput.Workflows, "Should have at least one workflow")

	found := false
	for _, wf := range listOutput.Workflows {
		if wf.Name == workflowName {
			found = true
			assert.Equal(t, cluster.ArgoNamespace, wf.Namespace)
			break
		}
	}
	assert.True(t, found, "Submitted workflow should be in the list")
}

// TestWorkflow_WaitWorkflow tests the wait_workflow tool handler directly.
func TestWorkflow_WaitWorkflow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	ctx := context.Background()
	cluster := SetupE2ECluster(ctx, t)

	// Load test workflow
	manifest := LoadTestDataFile(t, "hello-world.yaml")

	// Submit workflow
	t.Log("Submitting workflow...")
	submitHandler := tools.SubmitWorkflowHandler(cluster.ArgoClient)
	submitInput := tools.SubmitWorkflowInput{
		Namespace: cluster.ArgoNamespace,
		Manifest:  manifest,
	}

	_, submitOutput, err := submitHandler(ctx, nil, submitInput)
	require.NoError(t, err, "Failed to submit workflow")
	require.NotNil(t, submitOutput)

	workflowName := submitOutput.Name
	t.Logf("Submitted workflow: %s", workflowName)

	// Cleanup at the end
	defer func() {
		deleteHandler := tools.DeleteWorkflowHandler(cluster.ArgoClient)
		deleteInput := tools.DeleteWorkflowInput{
			Namespace: cluster.ArgoNamespace,
			Name:      workflowName,
		}
		_, _, _ = deleteHandler(ctx, nil, deleteInput)
	}()

	// Use wait_workflow tool handler to wait for completion
	t.Log("Waiting for workflow using wait_workflow tool...")
	waitHandler := tools.WaitWorkflowHandler(cluster.ArgoClient)
	waitInput := tools.WaitWorkflowInput{
		Namespace: cluster.ArgoNamespace,
		Name:      workflowName,
		Timeout:   "2m",
	}

	result, waitOutput, err := waitHandler(ctx, nil, waitInput)
	require.NoError(t, err, "wait_workflow should not return error")
	require.NotNil(t, result)
	require.NotNil(t, waitOutput)

	// Verify the wait output
	assert.Equal(t, workflowName, waitOutput.Name, "Workflow name should match")
	assert.Equal(t, cluster.ArgoNamespace, waitOutput.Namespace, "Namespace should match")
	assert.Equal(t, "Succeeded", waitOutput.Phase, "Workflow should succeed")
	assert.False(t, waitOutput.TimedOut, "Wait should not have timed out")
	assert.NotEmpty(t, waitOutput.StartedAt, "StartedAt should be set")
	assert.NotEmpty(t, waitOutput.FinishedAt, "FinishedAt should be set")
	assert.NotEmpty(t, waitOutput.Duration, "Duration should be calculated")
}

// TestWorkflow_WaitWorkflow_Timeout tests that wait_workflow times out correctly.
func TestWorkflow_WaitWorkflow_Timeout(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	ctx := context.Background()
	cluster := SetupE2ECluster(ctx, t)

	// Load DAG workflow (takes longer to complete)
	manifest := LoadTestDataFile(t, "dag-workflow.yaml")

	// Submit workflow
	t.Log("Submitting DAG workflow...")
	submitHandler := tools.SubmitWorkflowHandler(cluster.ArgoClient)
	submitInput := tools.SubmitWorkflowInput{
		Namespace: cluster.ArgoNamespace,
		Manifest:  manifest,
	}

	_, submitOutput, err := submitHandler(ctx, nil, submitInput)
	require.NoError(t, err, "Failed to submit workflow")
	require.NotNil(t, submitOutput)

	workflowName := submitOutput.Name
	t.Logf("Submitted workflow: %s", workflowName)

	// Cleanup at the end
	defer func() {
		// Terminate the workflow to clean up
		terminateHandler := tools.TerminateWorkflowHandler(cluster.ArgoClient)
		terminateInput := tools.TerminateWorkflowInput{
			Namespace: cluster.ArgoNamespace,
			Name:      workflowName,
		}
		_, _, _ = terminateHandler(ctx, nil, terminateInput)

		// Then delete it
		deleteHandler := tools.DeleteWorkflowHandler(cluster.ArgoClient)
		deleteInput := tools.DeleteWorkflowInput{
			Namespace: cluster.ArgoNamespace,
			Name:      workflowName,
		}
		_, _, _ = deleteHandler(ctx, nil, deleteInput)
	}()

	// Wait for workflow to start (so we're not timing out on a non-existent workflow)
	time.Sleep(2 * time.Second)

	// Use wait_workflow with very short timeout
	t.Log("Waiting for workflow with short timeout...")
	waitHandler := tools.WaitWorkflowHandler(cluster.ArgoClient)
	waitInput := tools.WaitWorkflowInput{
		Namespace: cluster.ArgoNamespace,
		Name:      workflowName,
		Timeout:   "3s", // Very short timeout - workflow won't complete in time
	}

	result, waitOutput, err := waitHandler(ctx, nil, waitInput)
	require.NoError(t, err, "wait_workflow should not return Go error on timeout")
	require.NotNil(t, result)
	require.NotNil(t, waitOutput)

	// Verify timeout behavior
	assert.Equal(t, workflowName, waitOutput.Name, "Workflow name should match")
	assert.True(t, waitOutput.TimedOut, "Wait should have timed out")
	assert.Contains(t, waitOutput.Message, "Timed out", "Message should indicate timeout")
}

// TestWorkflow_WatchWorkflow tests the watch_workflow tool handler.
func TestWorkflow_WatchWorkflow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	ctx := context.Background()
	cluster := SetupE2ECluster(ctx, t)

	// Load test workflow
	manifest := LoadTestDataFile(t, "hello-world.yaml")

	// Submit workflow
	t.Log("Submitting workflow...")
	submitHandler := tools.SubmitWorkflowHandler(cluster.ArgoClient)
	submitInput := tools.SubmitWorkflowInput{
		Namespace: cluster.ArgoNamespace,
		Manifest:  manifest,
	}

	_, submitOutput, err := submitHandler(ctx, nil, submitInput)
	require.NoError(t, err, "Failed to submit workflow")
	require.NotNil(t, submitOutput)

	workflowName := submitOutput.Name
	t.Logf("Submitted workflow: %s", workflowName)

	// Cleanup at the end
	defer func() {
		deleteHandler := tools.DeleteWorkflowHandler(cluster.ArgoClient)
		deleteInput := tools.DeleteWorkflowInput{
			Namespace: cluster.ArgoNamespace,
			Name:      workflowName,
		}
		_, _, _ = deleteHandler(ctx, nil, deleteInput)
	}()

	// Use watch_workflow tool handler to watch until completion
	t.Log("Watching workflow using watch_workflow tool...")
	watchHandler := tools.WatchWorkflowHandler(cluster.ArgoClient)
	watchInput := tools.WatchWorkflowInput{
		Namespace: cluster.ArgoNamespace,
		Name:      workflowName,
		Timeout:   "2m",
	}

	result, watchOutput, err := watchHandler(ctx, nil, watchInput)
	require.NoError(t, err, "watch_workflow should not return error")
	require.NotNil(t, result)
	require.NotNil(t, watchOutput)

	// Verify the watch output
	assert.Equal(t, workflowName, watchOutput.Name, "Workflow name should match")
	assert.Equal(t, cluster.ArgoNamespace, watchOutput.Namespace, "Namespace should match")
	assert.Equal(t, "Succeeded", watchOutput.Phase, "Workflow should succeed")
	assert.False(t, watchOutput.TimedOut, "Watch should not have timed out")
	assert.NotEmpty(t, watchOutput.StartedAt, "StartedAt should be set")
	assert.NotEmpty(t, watchOutput.FinishedAt, "FinishedAt should be set")
	assert.NotEmpty(t, watchOutput.Duration, "Duration should be calculated")

	// Watch-specific: verify events were collected
	assert.NotEmpty(t, watchOutput.Events, "Watch should have collected events")

	// Verify event structure
	for _, event := range watchOutput.Events {
		assert.NotEmpty(t, event.Type, "Event type should be set")
		assert.NotEmpty(t, event.Phase, "Event phase should be set")
		assert.NotEmpty(t, event.Timestamp, "Event timestamp should be set")
	}

	// Verify we have at least one event with Succeeded phase
	foundSucceeded := false
	for _, event := range watchOutput.Events {
		if event.Phase == "Succeeded" {
			foundSucceeded = true
			break
		}
	}
	assert.True(t, foundSucceeded, "Should have captured a Succeeded event")
}

// TestWorkflow_WatchWorkflow_Timeout tests that watch_workflow times out and captures events.
func TestWorkflow_WatchWorkflow_Timeout(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	ctx := context.Background()
	cluster := SetupE2ECluster(ctx, t)

	// Load DAG workflow (takes longer to complete)
	manifest := LoadTestDataFile(t, "dag-workflow.yaml")

	// Submit workflow
	t.Log("Submitting DAG workflow...")
	submitHandler := tools.SubmitWorkflowHandler(cluster.ArgoClient)
	submitInput := tools.SubmitWorkflowInput{
		Namespace: cluster.ArgoNamespace,
		Manifest:  manifest,
	}

	_, submitOutput, err := submitHandler(ctx, nil, submitInput)
	require.NoError(t, err, "Failed to submit workflow")
	require.NotNil(t, submitOutput)

	workflowName := submitOutput.Name
	t.Logf("Submitted workflow: %s", workflowName)

	// Cleanup at the end
	defer func() {
		// Terminate the workflow to clean up
		terminateHandler := tools.TerminateWorkflowHandler(cluster.ArgoClient)
		terminateInput := tools.TerminateWorkflowInput{
			Namespace: cluster.ArgoNamespace,
			Name:      workflowName,
		}
		_, _, _ = terminateHandler(ctx, nil, terminateInput)

		// Then delete it
		deleteHandler := tools.DeleteWorkflowHandler(cluster.ArgoClient)
		deleteInput := tools.DeleteWorkflowInput{
			Namespace: cluster.ArgoNamespace,
			Name:      workflowName,
		}
		_, _, _ = deleteHandler(ctx, nil, deleteInput)
	}()

	// Wait for workflow to start (so we capture some events)
	time.Sleep(2 * time.Second)

	// Use watch_workflow with short timeout
	t.Log("Watching workflow with short timeout...")
	watchHandler := tools.WatchWorkflowHandler(cluster.ArgoClient)
	watchInput := tools.WatchWorkflowInput{
		Namespace: cluster.ArgoNamespace,
		Name:      workflowName,
		Timeout:   "5s", // Short timeout - workflow won't complete in time
	}

	result, watchOutput, err := watchHandler(ctx, nil, watchInput)
	require.NoError(t, err, "watch_workflow should not return Go error on timeout")
	require.NotNil(t, result)
	require.NotNil(t, watchOutput)

	// Verify timeout behavior
	assert.Equal(t, workflowName, watchOutput.Name, "Workflow name should match")
	assert.True(t, watchOutput.TimedOut, "Watch should have timed out")
	assert.Contains(t, watchOutput.Message, "Watch timed out", "Message should indicate timeout")

	// Watch should still have captured some events before timing out
	assert.NotEmpty(t, watchOutput.Events, "Watch should have captured events before timeout")
}
