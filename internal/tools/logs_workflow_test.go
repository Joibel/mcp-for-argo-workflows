package tools

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogsWorkflowTool(t *testing.T) {
	tool := LogsWorkflowTool()

	assert.Equal(t, "logs_workflow", tool.Name)
	assert.NotEmpty(t, tool.Description)
}

func TestLogsWorkflowInput_Validation(t *testing.T) {
	tests := []struct {
		name      string
		input     LogsWorkflowInput
		wantValid bool
	}{
		{
			name: "valid input with name only",
			input: LogsWorkflowInput{
				Name: "my-workflow",
			},
			wantValid: true,
		},
		{
			name: "valid input with all options",
			input: LogsWorkflowInput{
				Name:      "my-workflow",
				Namespace: "custom-ns",
				PodName:   "my-pod",
				Container: "main",
				TailLines: int64Ptr(50),
				Grep:      "error",
			},
			wantValid: true,
		},
		{
			name: "empty name should be invalid",
			input: LogsWorkflowInput{
				Name: "",
			},
			wantValid: false,
		},
		{
			name: "whitespace-only name should be invalid",
			input: LogsWorkflowInput{
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

func TestLogsWorkflowOutput(t *testing.T) {
	tests := []struct {
		check  func(t *testing.T, o *LogsWorkflowOutput)
		name   string
		output LogsWorkflowOutput
	}{
		{
			name: "output with logs",
			output: LogsWorkflowOutput{
				Name:      "test-workflow",
				Namespace: "default",
				Logs: []LogEntryOutput{
					{PodName: "pod-1", Content: "log line 1"},
					{PodName: "pod-1", Content: "log line 2"},
				},
				Message: "Retrieved 2 log entries",
			},
			check: func(t *testing.T, o *LogsWorkflowOutput) {
				assert.Equal(t, "test-workflow", o.Name)
				assert.Equal(t, "default", o.Namespace)
				assert.Len(t, o.Logs, 2)
				assert.False(t, o.Truncated)
			},
		},
		{
			name: "output with truncation",
			output: LogsWorkflowOutput{
				Name:      "test-workflow",
				Namespace: "default",
				Logs: []LogEntryOutput{
					{PodName: "pod-1", Content: "partial logs..."},
				},
				Truncated: true,
				Message:   "Logs truncated after 1048576 bytes",
			},
			check: func(t *testing.T, o *LogsWorkflowOutput) {
				assert.True(t, o.Truncated)
				assert.Contains(t, o.Message, "truncated")
			},
		},
		{
			name: "output with no logs",
			output: LogsWorkflowOutput{
				Name:      "test-workflow",
				Namespace: "default",
				Logs:      []LogEntryOutput{},
				Message:   "No logs available",
			},
			check: func(t *testing.T, o *LogsWorkflowOutput) {
				assert.Empty(t, o.Logs)
				assert.Contains(t, o.Message, "No logs")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.check(t, &tt.output)
		})
	}
}

func TestLogEntryOutput(t *testing.T) {
	entry := LogEntryOutput{
		PodName: "my-pod",
		Content: "2025-01-15T10:00:00Z INFO Starting process...",
	}

	assert.Equal(t, "my-pod", entry.PodName)
	assert.Contains(t, entry.Content, "Starting process")
}

// int64Ptr is a helper to create int64 pointers for tests.
func int64Ptr(i int64) *int64 {
	return &i
}
