package tools

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogsWorkflowTool(t *testing.T) {
	tool := LogsWorkflowTool()

	assert.Equal(t, "logs_workflow", tool.Name)
	assert.NotEmpty(t, tool.Description)
}

func TestLogsWorkflowHandler_Validation(t *testing.T) {
	// Test handler validation directly - name validation occurs before client is used,
	// so we can pass nil and test that validation returns the expected errors.
	handler := LogsWorkflowHandler(nil)

	tests := []struct {
		input       LogsWorkflowInput
		name        string
		errContains string
		wantErr     bool
	}{
		{
			name: "empty name returns error",
			input: LogsWorkflowInput{
				Name: "",
			},
			wantErr:     true,
			errContains: "workflow name cannot be empty",
		},
		{
			name: "whitespace-only name returns error",
			input: LogsWorkflowInput{
				Name: "   ",
			},
			wantErr:     true,
			errContains: "workflow name cannot be empty",
		},
		{
			name: "whitespace-padded name returns error",
			input: LogsWorkflowInput{
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
