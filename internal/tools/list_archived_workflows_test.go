package tools

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/argoproj/argo-workflows/v3/pkg/apiclient/workflowarchive"
	wfv1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/Joibel/mcp-for-argo-workflows/internal/argo"
	"github.com/Joibel/mcp-for-argo-workflows/internal/argo/mocks"
)

func TestListArchivedWorkflowsTool(t *testing.T) {
	tool := ListArchivedWorkflowsTool()

	assert.Equal(t, "list_archived_workflows", tool.Name)
	assert.NotEmpty(t, tool.Description)
	assert.Contains(t, tool.Description, "archive")
	assert.Contains(t, tool.Description, "Argo Server")
}

func TestListArchivedWorkflowsInput(t *testing.T) {
	// Test default values
	input := ListArchivedWorkflowsInput{}
	assert.Empty(t, input.Namespace)
	assert.Empty(t, input.Labels)
	assert.Zero(t, input.Limit)

	// Test with values
	input2 := ListArchivedWorkflowsInput{
		Namespace: "test-namespace",
		Labels:    "app=test",
		Limit:     100,
	}
	assert.Equal(t, "test-namespace", input2.Namespace)
	assert.Equal(t, "app=test", input2.Labels)
	assert.Equal(t, int64(100), input2.Limit)
}

func TestArchivedWorkflowSummary(t *testing.T) {
	summary := ArchivedWorkflowSummary{
		UID:        "abc-123",
		Name:       "test-workflow",
		Namespace:  "default",
		Phase:      "Succeeded",
		CreatedAt:  "2025-01-01T00:00:00Z",
		FinishedAt: "2025-01-01T00:05:00Z",
		Message:    "Workflow completed",
	}

	assert.Equal(t, "abc-123", summary.UID)
	assert.Equal(t, "test-workflow", summary.Name)
	assert.Equal(t, "default", summary.Namespace)
	assert.Equal(t, "Succeeded", summary.Phase)
	assert.NotEmpty(t, summary.CreatedAt)
	assert.NotEmpty(t, summary.FinishedAt)
	assert.NotEmpty(t, summary.Message)
}

func TestListArchivedWorkflowsOutput(t *testing.T) {
	output := ListArchivedWorkflowsOutput{
		Workflows: []ArchivedWorkflowSummary{
			{UID: "uid-1", Name: "wf-1", Namespace: "default", Phase: "Succeeded"},
			{UID: "uid-2", Name: "wf-2", Namespace: "default", Phase: "Failed"},
		},
		Total: 2,
	}

	assert.Len(t, output.Workflows, 2)
	assert.Equal(t, 2, output.Total)
	assert.Equal(t, "wf-1", output.Workflows[0].Name)
	assert.Equal(t, "wf-2", output.Workflows[1].Name)
}

// newMockArchivedWorkflowService creates a new mock archived workflow service client.
func newMockArchivedWorkflowService(t *testing.T) *mocks.MockArchivedWorkflowServiceClient {
	t.Helper()
	m := &mocks.MockArchivedWorkflowServiceClient{}
	m.Test(t)
	return m
}

func TestListArchivedWorkflowsHandler(t *testing.T) {
	testTime := time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC)
	finishedTime := metav1.Time{Time: testTime.Add(5 * time.Minute)}

	tests := []struct {
		setupMock func(*mocks.MockArchivedWorkflowServiceClient)
		validate  func(*testing.T, *ListArchivedWorkflowsOutput, *mcp.CallToolResult)
		name      string
		input     ListArchivedWorkflowsInput
		wantErr   bool
	}{
		{
			name: "success - list archived workflows in namespace",
			input: ListArchivedWorkflowsInput{
				Namespace: "default",
			},
			setupMock: func(m *mocks.MockArchivedWorkflowServiceClient) {
				m.On("ListArchivedWorkflows", mock.Anything, mock.MatchedBy(func(req *workflowarchive.ListArchivedWorkflowsRequest) bool {
					return req.Namespace == "default"
				})).Return(
					&wfv1.WorkflowList{
						Items: []wfv1.Workflow{
							{
								ObjectMeta: metav1.ObjectMeta{
									Name:              "archived-wf-1",
									Namespace:         "default",
									UID:               types.UID("uid-001"),
									CreationTimestamp: metav1.Time{Time: testTime},
								},
								Status: wfv1.WorkflowStatus{
									Phase:      wfv1.WorkflowSucceeded,
									FinishedAt: finishedTime,
									Message:    "Workflow completed successfully",
								},
							},
							{
								ObjectMeta: metav1.ObjectMeta{
									Name:              "archived-wf-2",
									Namespace:         "default",
									UID:               types.UID("uid-002"),
									CreationTimestamp: metav1.Time{Time: testTime.Add(time.Hour)},
								},
								Status: wfv1.WorkflowStatus{
									Phase:      wfv1.WorkflowFailed,
									FinishedAt: metav1.Time{Time: testTime.Add(time.Hour + 10*time.Minute)},
									Message:    "Workflow failed",
								},
							},
						},
					},
					nil,
				)
			},
			wantErr: false,
			validate: func(t *testing.T, output *ListArchivedWorkflowsOutput, result *mcp.CallToolResult) {
				assert.Equal(t, 2, output.Total)
				assert.Len(t, output.Workflows, 2)
				assert.Equal(t, "archived-wf-1", output.Workflows[0].Name)
				assert.Equal(t, "default", output.Workflows[0].Namespace)
				assert.Equal(t, "uid-001", output.Workflows[0].UID)
				assert.Equal(t, "Succeeded", output.Workflows[0].Phase)
				assert.NotEmpty(t, output.Workflows[0].FinishedAt)
				assert.Equal(t, "Workflow completed successfully", output.Workflows[0].Message)
				assert.Equal(t, "archived-wf-2", output.Workflows[1].Name)
				assert.Equal(t, "Failed", output.Workflows[1].Phase)
				require.NotNil(t, result)
				require.Len(t, result.Content, 1)
				text, ok := result.Content[0].(*mcp.TextContent)
				require.True(t, ok)
				assert.Contains(t, text.Text, "Found 2 archived workflow(s)")
				assert.Contains(t, text.Text, "default")
			},
		},
		{
			name: "success - list archived workflows with labels",
			input: ListArchivedWorkflowsInput{
				Namespace: "default",
				Labels:    "app=myapp",
			},
			setupMock: func(m *mocks.MockArchivedWorkflowServiceClient) {
				m.On("ListArchivedWorkflows", mock.Anything, mock.MatchedBy(func(req *workflowarchive.ListArchivedWorkflowsRequest) bool {
					return req.Namespace == "default" && req.ListOptions != nil && req.ListOptions.LabelSelector == "app=myapp"
				})).Return(
					&wfv1.WorkflowList{
						Items: []wfv1.Workflow{
							{
								ObjectMeta: metav1.ObjectMeta{
									Name:              "archived-myapp",
									Namespace:         "default",
									UID:               types.UID("uid-003"),
									CreationTimestamp: metav1.Time{Time: testTime},
								},
								Status: wfv1.WorkflowStatus{
									Phase:      wfv1.WorkflowSucceeded,
									FinishedAt: finishedTime,
								},
							},
						},
					},
					nil,
				)
			},
			wantErr: false,
			validate: func(t *testing.T, output *ListArchivedWorkflowsOutput, _ *mcp.CallToolResult) {
				assert.Equal(t, 1, output.Total)
				assert.Equal(t, "archived-myapp", output.Workflows[0].Name)
			},
		},
		{
			name: "success - list with limit",
			input: ListArchivedWorkflowsInput{
				Namespace: "default",
				Limit:     10,
			},
			setupMock: func(m *mocks.MockArchivedWorkflowServiceClient) {
				m.On("ListArchivedWorkflows", mock.Anything, mock.MatchedBy(func(req *workflowarchive.ListArchivedWorkflowsRequest) bool {
					return req.ListOptions != nil && req.ListOptions.Limit == 10
				})).Return(
					&wfv1.WorkflowList{
						Items: []wfv1.Workflow{
							{
								ObjectMeta: metav1.ObjectMeta{
									Name:              "archived-limited",
									Namespace:         "default",
									UID:               types.UID("uid-004"),
									CreationTimestamp: metav1.Time{Time: testTime},
								},
								Status: wfv1.WorkflowStatus{
									Phase: wfv1.WorkflowSucceeded,
								},
							},
						},
					},
					nil,
				)
			},
			wantErr: false,
			validate: func(t *testing.T, output *ListArchivedWorkflowsOutput, _ *mcp.CallToolResult) {
				assert.Equal(t, 1, output.Total)
			},
		},
		{
			name:  "success - list all namespaces",
			input: ListArchivedWorkflowsInput{},
			setupMock: func(m *mocks.MockArchivedWorkflowServiceClient) {
				m.On("ListArchivedWorkflows", mock.Anything, mock.MatchedBy(func(req *workflowarchive.ListArchivedWorkflowsRequest) bool {
					return req.Namespace == ""
				})).Return(
					&wfv1.WorkflowList{
						Items: []wfv1.Workflow{
							{
								ObjectMeta: metav1.ObjectMeta{
									Name:              "wf-ns1",
									Namespace:         "namespace1",
									UID:               types.UID("uid-005"),
									CreationTimestamp: metav1.Time{Time: testTime},
								},
								Status: wfv1.WorkflowStatus{Phase: wfv1.WorkflowSucceeded},
							},
							{
								ObjectMeta: metav1.ObjectMeta{
									Name:              "wf-ns2",
									Namespace:         "namespace2",
									UID:               types.UID("uid-006"),
									CreationTimestamp: metav1.Time{Time: testTime},
								},
								Status: wfv1.WorkflowStatus{Phase: wfv1.WorkflowSucceeded},
							},
						},
					},
					nil,
				)
			},
			wantErr: false,
			validate: func(t *testing.T, output *ListArchivedWorkflowsOutput, result *mcp.CallToolResult) {
				assert.Equal(t, 2, output.Total)
				assert.Equal(t, "namespace1", output.Workflows[0].Namespace)
				assert.Equal(t, "namespace2", output.Workflows[1].Namespace)
				require.NotNil(t, result)
				require.Len(t, result.Content, 1)
				text, ok := result.Content[0].(*mcp.TextContent)
				require.True(t, ok)
				assert.Contains(t, text.Text, "across all namespaces")
			},
		},
		{
			name: "success - empty result",
			input: ListArchivedWorkflowsInput{
				Namespace: "empty-namespace",
			},
			setupMock: func(m *mocks.MockArchivedWorkflowServiceClient) {
				m.On("ListArchivedWorkflows", mock.Anything, mock.Anything).Return(
					&wfv1.WorkflowList{
						Items: []wfv1.Workflow{},
					},
					nil,
				)
			},
			wantErr: false,
			validate: func(t *testing.T, output *ListArchivedWorkflowsOutput, result *mcp.CallToolResult) {
				assert.Equal(t, 0, output.Total)
				assert.Empty(t, output.Workflows)
				require.NotNil(t, result)
				require.Len(t, result.Content, 1)
				text, ok := result.Content[0].(*mcp.TextContent)
				require.True(t, ok)
				assert.Contains(t, text.Text, "Found 0 archived workflow(s)")
			},
		},
		{
			name: "error - API failure",
			input: ListArchivedWorkflowsInput{
				Namespace: "default",
			},
			setupMock: func(m *mocks.MockArchivedWorkflowServiceClient) {
				m.On("ListArchivedWorkflows", mock.Anything, mock.Anything).Return(
					nil,
					errors.New("connection refused"),
				)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock client and service
			mockClient := newMockClient(t, "argo", true)
			mockService := newMockArchivedWorkflowService(t)
			mockClient.SetArchivedWorkflowService(mockService)

			// Setup mock expectations
			tt.setupMock(mockService)

			// Verify mock expectations even on error paths
			defer mockService.AssertExpectations(t)

			// Create handler and call it
			handler := ListArchivedWorkflowsHandler(mockClient)
			ctx := context.Background()
			req := &mcp.CallToolRequest{}

			result, output, err := handler(ctx, req, tt.input)

			// Validate results
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, output)
			tt.validate(t, output, result)
		})
	}
}

func TestListArchivedWorkflowsHandler_ArchivedServiceError(t *testing.T) {
	// Test when ArchivedWorkflowService returns an error (not in Argo Server mode)
	mockClient := newMockClient(t, "argo", false)
	mockClient.On("ArchivedWorkflowService").Return(nil, argo.ErrArchivedWorkflowsNotSupported)

	handler := ListArchivedWorkflowsHandler(mockClient)
	ctx := context.Background()
	req := &mcp.CallToolRequest{}
	input := ListArchivedWorkflowsInput{Namespace: "default"}

	_, _, err := handler(ctx, req, input)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "archived workflow service")
}
