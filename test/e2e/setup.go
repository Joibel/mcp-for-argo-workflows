//go:build e2e

// Package e2e contains end-to-end tests for the MCP server.
// Note: gosec security warnings are disabled for this test package as it intentionally
// uses exec.Command and file operations with test data.
package e2e

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/k3s"

	"github.com/Joibel/mcp-for-argo-workflows/internal/argo"
)

const (
	// ArgoVersion is the Argo Workflows version to install for E2E tests.
	ArgoVersion = "v3.6.2"

	// ArgoNamespace is the namespace where Argo Workflows is installed.
	ArgoNamespace = "argo"

	// ArgoQuickStartURL is the URL to the Argo quick-start manifest.
	ArgoQuickStartURL = "https://github.com/argoproj/argo-workflows/releases/download/" + ArgoVersion + "/quick-start-minimal.yaml"
)

// getProjectRoot returns the project root directory.
func getProjectRoot() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("failed to get caller information")
	}
	// Go up from test/e2e to project root
	return filepath.Join(filepath.Dir(file), "..", "..")
}

// buildBinary builds the MCP server binary and returns the path.
func buildBinary(t *testing.T) string {
	t.Helper()
	projectRoot := getProjectRoot()
	// Use test name to avoid conflicts with parallel test execution
	binaryPath := filepath.Join(projectRoot, "dist", fmt.Sprintf("mcp-for-argo-workflows-e2e-test-%s", t.Name()))

	//nolint:gosec // Building binaries in tests is expected
	buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/mcp-for-argo-workflows")
	buildCmd.Dir = projectRoot
	buildOutput, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Failed to build binary: %s", string(buildOutput))

	t.Cleanup(func() {
		if err := os.Remove(binaryPath); err != nil && !os.IsNotExist(err) {
			t.Logf("Failed to remove test binary: %v", err)
		}
	})

	return binaryPath
}

// E2ECluster represents a test cluster with Argo Workflows installed.
//
//nolint:revive // E2ECluster is clearer than Cluster for this test package
type E2ECluster struct {
	// ArgoClient is a configured Argo client for the cluster.
	ArgoClient *argo.Client

	// container is the k3s container (for cleanup).
	container *k3s.K3sContainer

	// Kubeconfig is the raw kubeconfig content.
	Kubeconfig string

	// KubeconfigPath is the path to the temporary kubeconfig file.
	KubeconfigPath string

	// ArgoNamespace is the namespace where Argo Workflows is installed.
	ArgoNamespace string
}

// SetupE2ECluster creates a k3s cluster, installs Argo Workflows, and returns
// a configured E2ECluster. The cluster is automatically torn down when the test completes.
func SetupE2ECluster(ctx context.Context, t *testing.T) *E2ECluster {
	t.Helper()

	t.Log("Starting k3s container...")

	// Start k3s container
	k3sContainer, err := k3s.Run(ctx, "rancher/k3s:v1.31.2-k3s1")
	require.NoError(t, err, "Failed to start k3s container")

	// Register cleanup to terminate container
	t.Cleanup(func() {
		t.Log("Terminating k3s container...")
		if termErr := k3sContainer.Terminate(context.Background()); termErr != nil {
			t.Logf("Failed to terminate k3s container: %v", termErr)
		}
	})

	// Get kubeconfig from container
	kubeconfig, err := k3sContainer.GetKubeConfig(ctx)
	require.NoError(t, err, "Failed to get kubeconfig from k3s")

	// Write kubeconfig to temp file
	kubeconfigFile, err := os.CreateTemp("", "e2e-kubeconfig-*.yaml")
	require.NoError(t, err, "Failed to create temp kubeconfig file")

	kubeconfigPath := kubeconfigFile.Name()
	t.Cleanup(func() {
		_ = os.Remove(kubeconfigPath) //nolint:errcheck // Cleanup is best-effort
	})

	_, err = kubeconfigFile.Write(kubeconfig)
	require.NoError(t, err, "Failed to write kubeconfig")
	err = kubeconfigFile.Close()
	require.NoError(t, err, "Failed to close kubeconfig file")

	t.Logf("Kubeconfig written to: %s", kubeconfigPath)

	// Install Argo Workflows
	t.Log("Installing Argo Workflows...")
	installArgoWorkflows(t, kubeconfigPath)

	// Wait for Argo controller to be ready
	t.Log("Waiting for Argo controller to be ready...")
	waitForArgoController(t, kubeconfigPath)

	// Create Argo client
	argoClient, err := argo.NewClient(ctx, &argo.Config{
		Kubeconfig: kubeconfigPath,
		Namespace:  ArgoNamespace,
	})
	require.NoError(t, err, "Failed to create Argo client")

	t.Log("E2E cluster setup complete")

	return &E2ECluster{
		Kubeconfig:     string(kubeconfig),
		KubeconfigPath: kubeconfigPath,
		ArgoNamespace:  ArgoNamespace,
		ArgoClient:     argoClient,
		container:      k3sContainer,
	}
}

// installArgoWorkflows installs Argo Workflows in the k3s cluster.
func installArgoWorkflows(t *testing.T, kubeconfigPath string) {
	t.Helper()

	// Create the argo namespace first (quick-start manifest expects it to exist)
	//nolint:gosec // Using kubectl in tests is expected
	nsCmd := exec.Command("kubectl", "create", "namespace", ArgoNamespace)
	nsCmd.Env = append(os.Environ(), "KUBECONFIG="+kubeconfigPath)
	output, err := nsCmd.CombinedOutput()
	require.NoError(t, err, "Failed to create argo namespace: %s", string(output))

	// Download the quick-start manifest to a temp file
	manifestFile, err := os.CreateTemp("", "argo-install-*.yaml")
	require.NoError(t, err, "Failed to create temp manifest file")
	manifestPath := manifestFile.Name()

	// Close file before curl writes to it
	err = manifestFile.Close()
	require.NoError(t, err, "Failed to close manifest file")

	defer func() {
		_ = os.Remove(manifestPath) //nolint:errcheck // Cleanup is best-effort
	}()

	// Download manifest
	//nolint:gosec // Using curl to download manifests in tests is expected
	downloadCmd := exec.Command("curl", "-sSL", "-o", manifestPath, ArgoQuickStartURL)
	output, err = downloadCmd.CombinedOutput()
	require.NoError(t, err, "Failed to download Argo quick-start manifest: %s", string(output))

	// Apply the manifest
	//nolint:gosec // Using kubectl in tests is expected
	applyCmd := exec.Command("kubectl", "apply", "-f", manifestPath)
	applyCmd.Env = append(os.Environ(), "KUBECONFIG="+kubeconfigPath)
	output, err = applyCmd.CombinedOutput()
	require.NoError(t, err, "Failed to apply Argo manifest: %s", string(output))

	t.Logf("Argo Workflows %s installed", ArgoVersion)
}

// waitForArgoController waits for the Argo controller deployment to be ready.
func waitForArgoController(t *testing.T, kubeconfigPath string) {
	t.Helper()

	// Wait for the argo namespace to exist
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			t.Fatal("Timeout waiting for Argo controller to be ready")
		case <-ticker.C:
			// Check if the deployment is ready
			cmd := exec.Command("kubectl", "wait", "--for=condition=available",
				"--timeout=5s",
				"-n", ArgoNamespace,
				"deployment/workflow-controller")
			cmd.Env = append(os.Environ(), "KUBECONFIG="+kubeconfigPath)
			output, err := cmd.CombinedOutput()

			if err == nil {
				t.Log("Argo controller is ready")
				return
			}

			// Log the error but continue waiting
			t.Logf("Waiting for Argo controller... %s", string(output))
		}
	}
}

// LoadTestDataFile reads a test data file from the testdata directory.
func LoadTestDataFile(t *testing.T, filename string) string {
	t.Helper()

	projectRoot := getProjectRoot()
	path := filepath.Join(projectRoot, "test", "e2e", "testdata", filename)

	//nolint:gosec // Reading test data files is expected
	data, err := os.ReadFile(path)
	require.NoError(t, err, "Failed to read test data file %s", filename)

	return string(data)
}

// WaitForWorkflowPhase polls the workflow status until it reaches one of the expected phases.
// Returns the final phase or fails the test if timeout is reached.
func (c *E2ECluster) WaitForWorkflowPhase(t *testing.T, namespace, name string, timeout time.Duration, phases ...string) string {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	phaseSet := make(map[string]bool)
	for _, p := range phases {
		phaseSet[p] = true
	}

	for {
		select {
		case <-ctx.Done():
			t.Fatalf("Timeout waiting for workflow %s/%s to reach phase %v", namespace, name, phases)
		case <-ticker.C:
			// Get workflow status
			cmd := exec.Command("kubectl", "get", "workflow", name,
				"-n", namespace,
				"-o", "jsonpath={.status.phase}")
			cmd.Env = append(os.Environ(), "KUBECONFIG="+c.KubeconfigPath)
			output, err := cmd.CombinedOutput()

			if err != nil {
				t.Logf("Error getting workflow status: %s", string(output))
				continue
			}

			phase := string(output)
			if phaseSet[phase] {
				t.Logf("Workflow %s/%s reached phase: %s", namespace, name, phase)
				return phase
			}

			t.Logf("Workflow %s/%s current phase: %s (waiting for %v)", namespace, name, phase, phases)
		}
	}
}

// WorkflowExists checks if a workflow exists in the cluster.
func (c *E2ECluster) WorkflowExists(t *testing.T, namespace, name string) bool {
	t.Helper()

	cmd := exec.Command("kubectl", "get", "workflow", name, "-n", namespace)
	cmd.Env = append(os.Environ(), "KUBECONFIG="+c.KubeconfigPath)
	err := cmd.Run()

	return err == nil
}

// WorkflowTemplateExists checks if a workflow template exists in the cluster.
func (c *E2ECluster) WorkflowTemplateExists(t *testing.T, namespace, name string) bool {
	t.Helper()

	cmd := exec.Command("kubectl", "get", "workflowtemplate", name, "-n", namespace)
	cmd.Env = append(os.Environ(), "KUBECONFIG="+c.KubeconfigPath)
	err := cmd.Run()

	return err == nil
}

// CronWorkflowExists checks if a cron workflow exists in the cluster.
func (c *E2ECluster) CronWorkflowExists(t *testing.T, namespace, name string) bool {
	t.Helper()

	cmd := exec.Command("kubectl", "get", "cronworkflow", name, "-n", namespace)
	cmd.Env = append(os.Environ(), "KUBECONFIG="+c.KubeconfigPath)
	err := cmd.Run()

	return err == nil
}

// GetWorkflowPhase returns the current phase of a workflow.
func (c *E2ECluster) GetWorkflowPhase(t *testing.T, namespace, name string) (string, error) {
	t.Helper()

	cmd := exec.Command("kubectl", "get", "workflow", name,
		"-n", namespace,
		"-o", "jsonpath={.status.phase}")
	cmd.Env = append(os.Environ(), "KUBECONFIG="+c.KubeconfigPath)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return "", fmt.Errorf("failed to get workflow phase: %w: %s", err, string(output))
	}

	return string(output), nil
}
