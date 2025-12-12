package tools

import (
	"context"
	"testing"

	wfv1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWatchWorkflowTool(t *testing.T) {
	tool := WatchWorkflowTool()

	assert.Equal(t, "watch_workflow", tool.Name)
	assert.NotEmpty(t, tool.Description)
}

func TestWatchWorkflowHandler_NameValidation(t *testing.T) {
	// Test handler validation directly - name validation occurs before client is used,
	// so we can pass nil and test that validation returns the expected errors.
	handler := WatchWorkflowHandler(nil)

	tests := []struct {
		input       WatchWorkflowInput
		name        string
		errContains string
		wantErr     bool
	}{
		{
			name: "empty name returns error",
			input: WatchWorkflowInput{
				Name: "",
			},
			wantErr:     true,
			errContains: "workflow name cannot be empty",
		},
		{
			name: "whitespace-only name returns error",
			input: WatchWorkflowInput{
				Name: "   ",
			},
			wantErr:     true,
			errContains: "workflow name cannot be empty",
		},
		{
			name: "whitespace-padded name returns error",
			input: WatchWorkflowInput{
				Name: "  \t\n  ",
			},
			wantErr:     true,
			errContains: "workflow name cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := handler(context.Background(), nil, tt.input)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
			}
		})
	}
}

func TestWatchWorkflowHandler_TimeoutValidation(t *testing.T) {
	// Test handler timeout validation directly - timeout validation occurs before
	// client is used (after name validation), so we can pass nil client.
	handler := WatchWorkflowHandler(nil)

	tests := []struct {
		input       WatchWorkflowInput
		name        string
		errContains string
		wantErr     bool
	}{
		{
			name: "invalid timeout format returns error",
			input: WatchWorkflowInput{
				Name:    "my-workflow",
				Timeout: "invalid",
			},
			wantErr:     true,
			errContains: "invalid timeout format",
		},
		{
			name: "negative timeout returns error",
			input: WatchWorkflowInput{
				Name:    "my-workflow",
				Timeout: "-5m",
			},
			wantErr:     true,
			errContains: "invalid timeout: must be a positive duration",
		},
		{
			name: "zero timeout returns error",
			input: WatchWorkflowInput{
				Name:    "my-workflow",
				Timeout: "0s",
			},
			wantErr:     true,
			errContains: "invalid timeout: must be a positive duration",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := handler(context.Background(), nil, tt.input)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
			}
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
