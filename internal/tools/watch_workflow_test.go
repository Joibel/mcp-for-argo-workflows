package tools

import (
	"strings"
	"testing"

	wfv1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/stretchr/testify/assert"
)

func TestWatchWorkflowTool(t *testing.T) {
	tool := WatchWorkflowTool()

	assert.Equal(t, "watch_workflow", tool.Name)
	assert.NotEmpty(t, tool.Description)
}

func TestWatchWorkflowInput_Validation(t *testing.T) {
	tests := []struct {
		name      string
		input     WatchWorkflowInput
		wantValid bool
	}{
		{
			name: "valid input with name only",
			input: WatchWorkflowInput{
				Name: "my-workflow",
			},
			wantValid: true,
		},
		{
			name: "valid input with namespace",
			input: WatchWorkflowInput{
				Name:      "my-workflow",
				Namespace: "custom-ns",
			},
			wantValid: true,
		},
		{
			name: "valid input with timeout",
			input: WatchWorkflowInput{
				Name:    "my-workflow",
				Timeout: "5m",
			},
			wantValid: true,
		},
		{
			name: "empty name should be invalid",
			input: WatchWorkflowInput{
				Name: "",
			},
			wantValid: false,
		},
		{
			name: "whitespace-only name should be invalid",
			input: WatchWorkflowInput{
				Name: "   ",
			},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := strings.TrimSpace(tt.input.Name) != ""
			assert.Equal(t, tt.wantValid, isValid)
		})
	}
}

func TestWatchWorkflowOutput(t *testing.T) {
	tests := []struct {
		check  func(t *testing.T, o *WatchWorkflowOutput)
		name   string
		output WatchWorkflowOutput
	}{
		{
			name: "successful completion",
			output: WatchWorkflowOutput{
				Name:       "test-workflow",
				Namespace:  "default",
				Phase:      "Succeeded",
				StartedAt:  "2025-01-15T10:00:00Z",
				FinishedAt: "2025-01-15T10:05:30Z",
				Duration:   "5m30s",
				Progress:   "3/3",
				Events: []WatchEventSummary{
					{Type: "ADDED", Phase: "Pending", Timestamp: "2025-01-15T10:00:00Z"},
					{Type: "MODIFIED", Phase: "Running", Timestamp: "2025-01-15T10:00:01Z", Progress: "1/3"},
					{Type: "MODIFIED", Phase: "Succeeded", Timestamp: "2025-01-15T10:05:30Z", Progress: "3/3"},
				},
			},
			check: func(t *testing.T, o *WatchWorkflowOutput) {
				assert.Equal(t, "test-workflow", o.Name)
				assert.Equal(t, "Succeeded", o.Phase)
				assert.Len(t, o.Events, 3)
				assert.False(t, o.TimedOut)
			},
		},
		{
			name: "timed out",
			output: WatchWorkflowOutput{
				Name:      "test-workflow",
				Namespace: "default",
				Phase:     "Running",
				Progress:  "1/3",
				TimedOut:  true,
				Message:   "Watch timed out after 5m. Last phase: Running",
			},
			check: func(t *testing.T, o *WatchWorkflowOutput) {
				assert.True(t, o.TimedOut)
				assert.Contains(t, o.Message, "timed out")
			},
		},
		{
			name: "failed workflow",
			output: WatchWorkflowOutput{
				Name:      "test-workflow",
				Namespace: "default",
				Phase:     "Failed",
				Message:   "Step failed: exit code 1",
			},
			check: func(t *testing.T, o *WatchWorkflowOutput) {
				assert.Equal(t, "Failed", o.Phase)
				assert.Contains(t, o.Message, "failed")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.check(t, &tt.output)
		})
	}
}

func TestWatchEventSummary(t *testing.T) {
	event := WatchEventSummary{
		Type:      "MODIFIED",
		Phase:     "Running",
		Timestamp: "2025-01-15T10:00:00Z",
		Progress:  "2/5",
	}

	assert.Equal(t, "MODIFIED", event.Type)
	assert.Equal(t, "Running", event.Phase)
	assert.Equal(t, "2/5", event.Progress)
	assert.Equal(t, "2025-01-15T10:00:00Z", event.Timestamp)
}

func TestIsWorkflowCompleted(t *testing.T) {
	tests := []struct {
		phase    wfv1.WorkflowPhase
		name     string
		expected bool
	}{
		{
			name:     "succeeded is completed",
			phase:    wfv1.WorkflowSucceeded,
			expected: true,
		},
		{
			name:     "failed is completed",
			phase:    wfv1.WorkflowFailed,
			expected: true,
		},
		{
			name:     "error is completed",
			phase:    wfv1.WorkflowError,
			expected: true,
		},
		{
			name:     "running is not completed",
			phase:    wfv1.WorkflowRunning,
			expected: false,
		},
		{
			name:     "pending is not completed",
			phase:    wfv1.WorkflowPending,
			expected: false,
		},
		{
			name:     "empty phase is not completed",
			phase:    "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isWorkflowCompleted(tt.phase)
			assert.Equal(t, tt.expected, result)
		})
	}
}
