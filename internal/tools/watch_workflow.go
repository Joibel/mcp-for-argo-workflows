// Package tools implements MCP tool handlers for Argo Workflows operations.
package tools

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/argoproj/argo-workflows/v3/pkg/apiclient/workflow"
	wfv1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/Joibel/mcp-for-argo-workflows/internal/argo"
)

// WatchWorkflowInput defines the input parameters for the watch_workflow tool.
type WatchWorkflowInput struct {
	// Namespace is the Kubernetes namespace (uses default if not specified).
	Namespace string `json:"namespace,omitempty" jsonschema:"description=Kubernetes namespace (uses default if not specified)"`

	// Name is the workflow name.
	Name string `json:"name" jsonschema:"description=Workflow name,required"`

	// Timeout is the maximum time to watch (e.g., '5m', '1h'). Default: no timeout.
	Timeout string `json:"timeout,omitempty" jsonschema:"description=Maximum time to watch (e.g. 5m or 1h). Default: no timeout"`
}

// WatchWorkflowOutput defines the output for the watch_workflow tool.
type WatchWorkflowOutput struct {
	// Name is the workflow name.
	Name string `json:"name"`

	// Namespace is the namespace of the workflow.
	Namespace string `json:"namespace"`

	// Phase is the final workflow phase.
	Phase string `json:"phase"`

	// Message provides additional status information.
	Message string `json:"message,omitempty"`

	// StartedAt is when the workflow started.
	StartedAt string `json:"startedAt,omitempty"`

	// FinishedAt is when the workflow finished.
	FinishedAt string `json:"finishedAt,omitempty"`

	// Duration is the workflow duration in a human-readable format.
	Duration string `json:"duration,omitempty"`

	// Progress shows completed/total nodes.
	Progress string `json:"progress,omitempty"`

	// Events is a summary of watch events received.
	Events []WatchEventSummary `json:"events,omitempty"`

	// TimedOut indicates if the watch operation timed out.
	TimedOut bool `json:"timedOut,omitempty"`
}

// WatchEventSummary provides a summary of a watch event.
type WatchEventSummary struct {
	// Type is the event type (ADDED, MODIFIED, DELETED).
	Type string `json:"type"`

	// Phase is the workflow phase at this event.
	Phase string `json:"phase"`

	// Timestamp is when the event was received.
	Timestamp string `json:"timestamp"`

	// Progress shows completed/total nodes at this event.
	Progress string `json:"progress,omitempty"`
}

// WatchWorkflowTool returns the MCP tool definition for watch_workflow.
func WatchWorkflowTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "watch_workflow",
		Description: "Watch an Argo Workflow and stream status updates until completion",
	}
}

// WatchWorkflowHandler returns a handler function for the watch_workflow tool.
func WatchWorkflowHandler(client *argo.Client) func(context.Context, *mcp.CallToolRequest, WatchWorkflowInput) (*mcp.CallToolResult, *WatchWorkflowOutput, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input WatchWorkflowInput) (*mcp.CallToolResult, *WatchWorkflowOutput, error) {
		// Validate and normalize name
		input.Name = strings.TrimSpace(input.Name)
		if input.Name == "" {
			return nil, nil, fmt.Errorf("workflow name cannot be empty")
		}

		// Determine namespace
		namespace := input.Namespace
		if namespace == "" {
			namespace = client.DefaultNamespace()
		}

		// Parse timeout if provided
		var timeout time.Duration
		if input.Timeout != "" {
			var err error
			timeout, err = time.ParseDuration(input.Timeout)
			if err != nil {
				return nil, nil, fmt.Errorf("invalid timeout format: %w", err)
			}
			if timeout <= 0 {
				return nil, nil, fmt.Errorf("invalid timeout: must be a positive duration")
			}
		}

		// Create a context with timeout or cancellation for cleanup
		var watchCtx context.Context
		var cancel context.CancelFunc
		if timeout > 0 {
			watchCtx, cancel = context.WithTimeout(ctx, timeout)
		} else {
			watchCtx, cancel = context.WithCancel(ctx)
		}
		defer cancel()

		// Build the request with field selector to watch specific workflow
		req := &workflow.WatchWorkflowsRequest{
			Namespace: namespace,
			ListOptions: &metav1.ListOptions{
				FieldSelector: fmt.Sprintf("metadata.name=%s", input.Name),
			},
		}

		// Get the workflow service client
		wfService := client.WorkflowService()

		// Start watching
		stream, err := wfService.WatchWorkflows(watchCtx, req)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to watch workflow: %w", err)
		}

		// Collect events and watch until completion
		var events []WatchEventSummary
		var lastWorkflow *wfv1.Workflow
		timedOut := false

		for {
			event, recvErr := stream.Recv()
			if errors.Is(recvErr, io.EOF) {
				break
			}
			if recvErr != nil {
				// Check if it was a timeout (handle both context and gRPC status)
				if errors.Is(recvErr, context.DeadlineExceeded) ||
					errors.Is(watchCtx.Err(), context.DeadlineExceeded) ||
					status.Code(recvErr) == codes.DeadlineExceeded {
					timedOut = true
					break
				}
				return nil, nil, fmt.Errorf("failed to receive watch event: %w", recvErr)
			}

			if event.Object == nil {
				continue
			}

			lastWorkflow = event.Object

			// Create event summary
			eventSummary := WatchEventSummary{
				Type:      event.Type,
				Phase:     string(event.Object.Status.Phase),
				Timestamp: time.Now().Format(time.RFC3339),
			}
			if event.Object.Status.Progress != "" {
				eventSummary.Progress = string(event.Object.Status.Progress)
			}
			events = append(events, eventSummary)

			// Check if workflow has completed
			if isWorkflowCompleted(event.Object.Status.Phase) {
				break
			}
		}

		// Build the output
		output := &WatchWorkflowOutput{
			Name:      input.Name,
			Namespace: namespace,
			Events:    events,
			TimedOut:  timedOut,
		}

		if lastWorkflow != nil {
			output.Phase = string(lastWorkflow.Status.Phase)
			output.Message = lastWorkflow.Status.Message
			output.Progress = string(lastWorkflow.Status.Progress)

			if !lastWorkflow.Status.StartedAt.Time.IsZero() {
				output.StartedAt = lastWorkflow.Status.StartedAt.Format(time.RFC3339)
			}
			if !lastWorkflow.Status.FinishedAt.Time.IsZero() {
				output.FinishedAt = lastWorkflow.Status.FinishedAt.Format(time.RFC3339)
			}

			// Calculate duration
			if !lastWorkflow.Status.StartedAt.Time.IsZero() {
				endTime := lastWorkflow.Status.FinishedAt.Time
				if endTime.IsZero() {
					endTime = time.Now()
				}
				duration := endTime.Sub(lastWorkflow.Status.StartedAt.Time)
				output.Duration = formatDuration(duration)
			}
		} else {
			output.Phase = "Unknown"
			output.Message = "No workflow events received"
		}

		if timedOut {
			if input.Timeout != "" {
				output.Message = fmt.Sprintf("Watch timed out after %s. Last phase: %s", input.Timeout, output.Phase)
			} else {
				output.Message = fmt.Sprintf("Watch timed out. Last phase: %s", output.Phase)
			}
		}

		return nil, output, nil
	}
}

// isWorkflowCompleted checks if a workflow phase indicates completion.
func isWorkflowCompleted(phase wfv1.WorkflowPhase) bool {
	switch phase {
	case wfv1.WorkflowSucceeded, wfv1.WorkflowFailed, wfv1.WorkflowError:
		return true
	default:
		return false
	}
}
