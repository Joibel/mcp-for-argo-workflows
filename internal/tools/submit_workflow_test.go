package tools

import (
	"testing"

	wfv1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
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
			if (err != nil) != tt.wantErr {
				t.Errorf("applyParameterOverrides() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			gotParams := tt.workflow.Spec.Arguments.Parameters
			if len(gotParams) != len(tt.wantParams) {
				t.Errorf("got %d parameters, want %d", len(gotParams), len(tt.wantParams))
				return
			}

			for i, want := range tt.wantParams {
				got := gotParams[i]
				if got.Name != want.Name {
					t.Errorf("param[%d].Name = %q, want %q", i, got.Name, want.Name)
				}
				if got.Value.String() != want.Value.String() {
					t.Errorf("param[%d].Value = %q, want %q", i, got.Value.String(), want.Value.String())
				}
			}
		})
	}
}

func TestSubmitWorkflowTool(t *testing.T) {
	tool := SubmitWorkflowTool()

	if tool.Name != "submit_workflow" {
		t.Errorf("tool.Name = %q, want %q", tool.Name, "submit_workflow")
	}

	if tool.Description == "" {
		t.Error("tool.Description should not be empty")
	}
}
