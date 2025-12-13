package tools

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWaitWorkflowTool(t *testing.T) {
	tool := WaitWorkflowTool()

	assert.Equal(t, "wait_workflow", tool.Name)
	assert.NotEmpty(t, tool.Description)
}

func TestWaitWorkflowHandler_NameValidation(t *testing.T) {
	// Test handler validation directly - name validation occurs before client is used,
	// so we can pass nil and test that validation returns the expected errors.
	handler := WaitWorkflowHandler(nil)

	tests := []struct {
		input       WaitWorkflowInput
		name        string
		errContains string
		wantErr     bool
	}{
		{
			name: "empty name returns error",
			input: WaitWorkflowInput{
				Name: "",
			},
			wantErr:     true,
			errContains: "workflow name cannot be empty",
		},
		{
			name: "whitespace-only name returns error",
			input: WaitWorkflowInput{
				Name: "   ",
			},
			wantErr:     true,
			errContains: "workflow name cannot be empty",
		},
		{
			name: "whitespace-padded name returns error",
			input: WaitWorkflowInput{
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

func TestWaitWorkflowHandler_TimeoutValidation(t *testing.T) {
	// Test handler timeout validation directly - timeout validation occurs before
	// client is used (after name validation), so we can pass nil client.
	handler := WaitWorkflowHandler(nil)

	tests := []struct {
		input       WaitWorkflowInput
		name        string
		errContains string
		wantErr     bool
	}{
		{
			name: "invalid timeout format returns error",
			input: WaitWorkflowInput{
				Name:    "my-workflow",
				Timeout: "invalid",
			},
			wantErr:     true,
			errContains: "invalid timeout format",
		},
		{
			name: "negative timeout returns error",
			input: WaitWorkflowInput{
				Name:    "my-workflow",
				Timeout: "-5m",
			},
			wantErr:     true,
			errContains: "invalid timeout: must be a positive duration",
		},
		{
			name: "zero timeout returns error",
			input: WaitWorkflowInput{
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

func TestWaitWorkflowOutput(t *testing.T) {
	tests := []struct {
		check  func(t *testing.T, o *WaitWorkflowOutput)
		name   string
		output WaitWorkflowOutput
	}{
		{
			name: "output has expected fields",
			output: WaitWorkflowOutput{
				Name:       "test-workflow",
				Namespace:  "default",
				Phase:      "Succeeded",
				Message:    "test message",
				StartedAt:  "2024-01-01T00:00:00Z",
				FinishedAt: "2024-01-01T00:05:00Z",
				Duration:   "5m0s",
				Progress:   "3/3",
				TimedOut:   false,
			},
			check: func(t *testing.T, o *WaitWorkflowOutput) {
				assert.Equal(t, "test-workflow", o.Name)
				assert.Equal(t, "default", o.Namespace)
				assert.Equal(t, "Succeeded", o.Phase)
				assert.Equal(t, "test message", o.Message)
				assert.Equal(t, "2024-01-01T00:00:00Z", o.StartedAt)
				assert.Equal(t, "2024-01-01T00:05:00Z", o.FinishedAt)
				assert.Equal(t, "5m0s", o.Duration)
				assert.Equal(t, "3/3", o.Progress)
				assert.False(t, o.TimedOut)
			},
		},
		{
			name: "timed out output",
			output: WaitWorkflowOutput{
				Name:      "test-workflow",
				Namespace: "default",
				Phase:     "Running",
				TimedOut:  true,
			},
			check: func(t *testing.T, o *WaitWorkflowOutput) {
				assert.True(t, o.TimedOut)
				assert.Equal(t, "Running", o.Phase)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.check(t, &tt.output)
		})
	}
}
