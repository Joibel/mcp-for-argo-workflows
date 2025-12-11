package argo

import (
	"context"
	"errors"
	"fmt"

	"github.com/argoproj/argo-workflows/v3/pkg/apiclient"
	"github.com/argoproj/argo-workflows/v3/pkg/apiclient/clusterworkflowtemplate"
	"github.com/argoproj/argo-workflows/v3/pkg/apiclient/cronworkflow"
	"github.com/argoproj/argo-workflows/v3/pkg/apiclient/info"
	"github.com/argoproj/argo-workflows/v3/pkg/apiclient/workflow"
	"github.com/argoproj/argo-workflows/v3/pkg/apiclient/workflowarchive"
	"github.com/argoproj/argo-workflows/v3/pkg/apiclient/workflowtemplate"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	// ErrArchivedWorkflowsNotSupported is returned when trying to access archived workflows
	// in direct Kubernetes API mode.
	ErrArchivedWorkflowsNotSupported = errors.New("archived workflows are only supported with Argo Server connection")
)

// Client wraps the Argo Workflows API client and provides service client getters.
type Client struct {
	config    *Config
	apiClient apiclient.Client
	ctx       context.Context
}

// NewClient creates a new Argo Workflows client based on the provided configuration.
// It supports two connection modes:
//   - Argo Server mode: When config.ArgoServer is set, connects via gRPC
//   - Direct K8s mode: When config.ArgoServer is empty, uses kubeconfig
func NewClient(config *Config) (*Client, error) {
	if config == nil {
		return nil, errors.New("config cannot be nil")
	}

	var opts apiclient.Opts

	if config.ArgoServer != "" {
		// Argo Server mode
		opts = apiclient.Opts{
			ArgoServerOpts: apiclient.ArgoServerOpts{
				URL:                config.ArgoServer,
				Secure:             config.Secure,
				InsecureSkipVerify: config.InsecureSkipVerify,
			},
		}

		// Add auth supplier if token is provided
		if config.ArgoToken != "" {
			opts.AuthSupplier = func() string {
				return config.ArgoToken
			}
		}
	} else {
		// Direct Kubernetes API mode
		opts = apiclient.Opts{}

		// Set kubeconfig if provided
		if config.Kubeconfig != "" {
			opts.ClientConfigSupplier = func() clientcmd.ClientConfig {
				return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
					&clientcmd.ClientConfigLoadingRules{ExplicitPath: config.Kubeconfig},
					&clientcmd.ConfigOverrides{},
				)
			}
		}
	}

	// Set context for API calls
	opts.Context = context.Background()

	ctx, apiClient, err := apiclient.NewClientFromOpts(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create Argo API client: %w", err)
	}

	return &Client{
		config:    config,
		apiClient: apiClient,
		ctx:       ctx,
	}, nil
}

// WorkflowService returns the workflow service client.
func (c *Client) WorkflowService() workflow.WorkflowServiceClient {
	return c.apiClient.NewWorkflowServiceClient()
}

// Note: The WorkflowServiceClient is returned directly without error since
// the Argo API client always returns a valid client for this service.

// CronWorkflowService returns the cron workflow service client.
func (c *Client) CronWorkflowService() (cronworkflow.CronWorkflowServiceClient, error) {
	client, err := c.apiClient.NewCronWorkflowServiceClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create cron workflow service client: %w", err)
	}
	return client, nil
}

// WorkflowTemplateService returns the workflow template service client.
func (c *Client) WorkflowTemplateService() (workflowtemplate.WorkflowTemplateServiceClient, error) {
	client, err := c.apiClient.NewWorkflowTemplateServiceClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create workflow template service client: %w", err)
	}
	return client, nil
}

// ClusterWorkflowTemplateService returns the cluster workflow template service client.
func (c *Client) ClusterWorkflowTemplateService() (clusterworkflowtemplate.ClusterWorkflowTemplateServiceClient, error) {
	client, err := c.apiClient.NewClusterWorkflowTemplateServiceClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create cluster workflow template service client: %w", err)
	}
	return client, nil
}

// ArchivedWorkflowService returns the archived workflow service client.
// This service is only available when using Argo Server mode.
// Returns ErrArchivedWorkflowsNotSupported if not in Argo Server mode.
func (c *Client) ArchivedWorkflowService() (workflowarchive.ArchivedWorkflowServiceClient, error) {
	if !c.IsArgoServerMode() {
		return nil, ErrArchivedWorkflowsNotSupported
	}

	client, err := c.apiClient.NewArchivedWorkflowServiceClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create archived workflow service client: %w", err)
	}
	return client, nil
}

// InfoService returns the info service client.
func (c *Client) InfoService() (info.InfoServiceClient, error) {
	client, err := c.apiClient.NewInfoServiceClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create info service client: %w", err)
	}
	return client, nil
}

// IsArgoServerMode returns true if the client is connected via Argo Server,
// false if using direct Kubernetes API access.
func (c *Client) IsArgoServerMode() bool {
	return c.config.ArgoServer != ""
}

// DefaultNamespace returns the default namespace configured for this client.
func (c *Client) DefaultNamespace() string {
	return c.config.Namespace
}

// Context returns the context associated with this client.
func (c *Client) Context() context.Context {
	return c.ctx
}
