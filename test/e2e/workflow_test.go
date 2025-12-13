//go:build e2e

package e2e

import (
	"context"
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
	assert.Contains(t, logsOutput.Logs, "Hello World", "Logs should contain expected output")

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
	wf, err := wfService.GetWorkflow(ctx, &workflow.WorkflowGetRequest{
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
						if assert.Contains(t, errMsg, tt.errContains) {
							found = true
							break
						}
					}
					assert.True(t, found, "Should contain error about %s", tt.errContains)
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
