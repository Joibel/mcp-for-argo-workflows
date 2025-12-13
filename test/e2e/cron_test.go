//go:build e2e

package e2e

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestCronWorkflow_CRUD tests the full CRUD lifecycle: create → get → list → suspend → resume → delete.
//
// Note: This test is currently skipped because CronWorkflow tools have not been implemented yet.
// Once the following tools are implemented in internal/tools/, this test should be enabled:
// - create_cron_workflow
// - get_cron_workflow
// - list_cron_workflows
// - suspend_cron_workflow
// - resume_cron_workflow
// - delete_cron_workflow
func TestCronWorkflow_CRUD(t *testing.T) {
	t.Skip("CronWorkflow tools not yet implemented (PIP-74 scope: E2E infrastructure only)")

	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	ctx := context.Background()
	cluster := SetupE2ECluster(ctx, t)

	// Load test cron workflow
	manifest := LoadTestDataFile(t, "cron-workflow.yaml")
	_ = manifest

	// TODO: Implement when tools are available
	// Step 1: Create cron workflow
	// createHandler := tools.CreateCronWorkflowHandler(cluster.ArgoClient)
	// createInput := tools.CreateCronWorkflowInput{
	//     Namespace: cluster.ArgoNamespace,
	//     Manifest:  manifest,
	// }
	// _, createOutput, err := createHandler(ctx, nil, createInput)
	// require.NoError(t, err, "Failed to create cron workflow")

	// Step 2: Get cron workflow
	// getHandler := tools.GetCronWorkflowHandler(cluster.ArgoClient)
	// getInput := tools.GetCronWorkflowInput{
	//     Namespace: cluster.ArgoNamespace,
	//     Name:      cronWorkflowName,
	// }
	// _, getOutput, err := getHandler(ctx, nil, getInput)
	// require.NoError(t, err, "Failed to get cron workflow")

	// Step 3: List cron workflows
	// listHandler := tools.ListCronWorkflowsHandler(cluster.ArgoClient)
	// listInput := tools.ListCronWorkflowsInput{
	//     Namespace: cluster.ArgoNamespace,
	// }
	// _, listOutput, err := listHandler(ctx, nil, listInput)
	// require.NoError(t, err, "Failed to list cron workflows")

	// Step 4: Suspend cron workflow
	// suspendHandler := tools.SuspendCronWorkflowHandler(cluster.ArgoClient)
	// suspendInput := tools.SuspendCronWorkflowInput{
	//     Namespace: cluster.ArgoNamespace,
	//     Name:      cronWorkflowName,
	// }
	// _, suspendOutput, err := suspendHandler(ctx, nil, suspendInput)
	// require.NoError(t, err, "Failed to suspend cron workflow")

	// Step 5: Resume cron workflow
	// resumeHandler := tools.ResumeCronWorkflowHandler(cluster.ArgoClient)
	// resumeInput := tools.ResumeCronWorkflowInput{
	//     Namespace: cluster.ArgoNamespace,
	//     Name:      cronWorkflowName,
	// }
	// _, resumeOutput, err := resumeHandler(ctx, nil, resumeInput)
	// require.NoError(t, err, "Failed to resume cron workflow")

	// Step 6: Delete cron workflow
	// deleteHandler := tools.DeleteCronWorkflowHandler(cluster.ArgoClient)
	// deleteInput := tools.DeleteCronWorkflowInput{
	//     Namespace: cluster.ArgoNamespace,
	//     Name:      cronWorkflowName,
	// }
	// _, deleteOutput, err := deleteHandler(ctx, nil, deleteInput)
	// require.NoError(t, err, "Failed to delete cron workflow")

	_ = ctx
	_ = cluster
	_ = require.NoError
}

// TestCronWorkflow_Schedule tests that a cron workflow is properly scheduled.
//
// Note: This test is currently skipped because CronWorkflow tools have not been implemented yet.
func TestCronWorkflow_Schedule(t *testing.T) {
	t.Skip("CronWorkflow tools not yet implemented (PIP-74 scope: E2E infrastructure only)")

	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	// TODO: Implement when tools are available
	// This test should:
	// 1. Create a cron workflow with a frequent schedule (e.g., every minute)
	// 2. Wait for the cron workflow to trigger at least once
	// 3. List workflows to verify a workflow instance was created
	// 4. Clean up the cron workflow
}
