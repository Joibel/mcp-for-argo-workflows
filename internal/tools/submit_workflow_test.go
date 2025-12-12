package tools

import (
	"testing"

	wfv1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApplyParameterOverrides(t *testing.T) {
	tests := []struct {
		wantParams []wfv1.Parameter
		workflow   *wfv1.Workflow
		name       string
		params     []string
		wantErr    bool
	}{
		{
			name: "update existing parameter",
			workflow: &wfv1.Workflow{
				Spec: wfv1.WorkflowSpec{
					Arguments: wfv1.Arguments{
						Parameters: []wfv1.Parameter{
							{Name: "message", Value: wfv1.AnyStringPtr("hello")},
						},
					},
				},
			},
			params:  []string{"message=world"},
			wantErr: false,
			wantParams: []wfv1.Parameter{
				{Name: "message", Value: wfv1.AnyStringPtr("world")},
			},
		},
		{
			name: "add new parameter",
			workflow: &wfv1.Workflow{
				Spec: wfv1.WorkflowSpec{
					Arguments: wfv1.Arguments{
						Parameters: []wfv1.Parameter{},
					},
				},
			},
			params:  []string{"newparam=newvalue"},
			wantErr: false,
			wantParams: []wfv1.Parameter{
				{Name: "newparam", Value: wfv1.AnyStringPtr("newvalue")},
			},
		},
		{
			name: "multiple parameters",
			workflow: &wfv1.Workflow{
				Spec: wfv1.WorkflowSpec{
					Arguments: wfv1.Arguments{
						Parameters: []wfv1.Parameter{
							{Name: "existing", Value: wfv1.AnyStringPtr("old")},
						},
					},
				},
			},
			params:  []string{"existing=new", "another=value"},
			wantErr: false,
			wantParams: []wfv1.Parameter{
				{Name: "existing", Value: wfv1.AnyStringPtr("new")},
				{Name: "another", Value: wfv1.AnyStringPtr("value")},
			},
		},
		{
			name: "parameter with equals in value",
			workflow: &wfv1.Workflow{
				Spec: wfv1.WorkflowSpec{
					Arguments: wfv1.Arguments{
						Parameters: []wfv1.Parameter{},
					},
				},
			},
			params:  []string{"equation=a=b+c"},
			wantErr: false,
			wantParams: []wfv1.Parameter{
				{Name: "equation", Value: wfv1.AnyStringPtr("a=b+c")},
			},
		},
		{
			name: "invalid parameter format - no equals",
			workflow: &wfv1.Workflow{
				Spec: wfv1.WorkflowSpec{
					Arguments: wfv1.Arguments{
						Parameters: []wfv1.Parameter{},
					},
				},
			},
			params:  []string{"invalid"},
			wantErr: true,
		},
		{
			name: "empty workflow arguments",
			workflow: &wfv1.Workflow{
				Spec: wfv1.WorkflowSpec{},
			},
			params:  []string{"param=value"},
			wantErr: false,
			wantParams: []wfv1.Parameter{
				{Name: "param", Value: wfv1.AnyStringPtr("value")},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := applyParameterOverrides(tt.workflow, tt.params)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			gotParams := tt.workflow.Spec.Arguments.Parameters
			require.Len(t, gotParams, len(tt.wantParams))

			for i, want := range tt.wantParams {
				got := gotParams[i]
				assert.Equal(t, want.Name, got.Name, "param[%d].Name mismatch", i)
				assert.Equal(t, want.Value.String(), got.Value.String(), "param[%d].Value mismatch", i)
			}
		})
	}
}

func TestSubmitWorkflowTool(t *testing.T) {
	tool := SubmitWorkflowTool()

	assert.Equal(t, "submit_workflow", tool.Name)
	assert.NotEmpty(t, tool.Description)
}
