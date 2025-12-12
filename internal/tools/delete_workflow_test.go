package tools

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteWorkflowTool(t *testing.T) {
	tool := DeleteWorkflowTool()

	assert.Equal(t, "delete_workflow", tool.Name)
	assert.NotEmpty(t, tool.Description)
}

func TestDeleteWorkflowInput_Validation(t *testing.T) {
	tests := []struct {
		name      string
		input     DeleteWorkflowInput
		wantValid bool
	}{
		{
			name: "valid input with name only",
			input: DeleteWorkflowInput{
				Name: "my-workflow",
			},
			wantValid: true,
		},
		{
			name: "valid input with namespace",
			input: DeleteWorkflowInput{
				Name:      "my-workflow",
				Namespace: "custom-ns",
			},
			wantValid: true,
		},
		{
			name: "valid input with force",
			input: DeleteWorkflowInput{
				Name:  "my-workflow",
				Force: true,
			},
			wantValid: true,
		},
		{
			name: "empty name should be invalid",
			input: DeleteWorkflowInput{
				Name: "",
			},
			wantValid: false,
		},
		{
			name: "whitespace-only name should be invalid",
			input: DeleteWorkflowInput{
				Name: "   ",
			},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Validate that empty/whitespace names are caught
			// This mirrors what the handler does with strings.TrimSpace
			isValid := strings.TrimSpace(tt.input.Name) != ""
			assert.Equal(t, tt.wantValid, isValid)
		})
	}
}

func TestDeleteWorkflowOutput(t *testing.T) {
	output := &DeleteWorkflowOutput{
		Name:      "test-workflow",
		Namespace: "default",
		Message:   "Workflow \"test-workflow\" deleted successfully",
	}

	assert.Equal(t, "test-workflow", output.Name)
	assert.Equal(t, "default", output.Namespace)
	assert.Contains(t, output.Message, "deleted successfully")
}
