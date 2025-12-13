//go:build e2e

package e2e

import (
	"testing"
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
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	t.Skip("CronWorkflow tools not yet implemented (PIP-74 scope: E2E infrastructure only)")
}

// TestCronWorkflow_Schedule tests that a cron workflow is properly scheduled.
//
// Note: This test is currently skipped because CronWorkflow tools have not been implemented yet.
func TestCronWorkflow_Schedule(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	t.Skip("CronWorkflow tools not yet implemented (PIP-74 scope: E2E infrastructure only)")
}
