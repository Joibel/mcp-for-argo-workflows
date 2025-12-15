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
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/k3s"

	"github.com/Joibel/mcp-for-argo-workflows/internal/argo"
)

// ConnectionMode specifies how the E2E tests connect to Argo Workflows.
type ConnectionMode string

const (
	// ModeKubernetesAPI uses direct Kubernetes API access (default).
	ModeKubernetesAPI ConnectionMode = "kubernetes"

	// ModeArgoServer connects via Argo Server API.
	ModeArgoServer ConnectionMode = "argo-server"
)

const (
	// ArgoVersion is the Argo Workflows version to install for E2E tests.
	ArgoVersion = "v3.6.2"

	// ArgoNamespace is the namespace where Argo Workflows is installed.
	ArgoNamespace = "argo"

	// ArgoQuickStartURL is the URL to the Argo quick-start manifest.
	ArgoQuickStartURL = "https://github.com/argoproj/argo-workflows/releases/download/" + ArgoVersion + "/quick-start-minimal.yaml"

	// ArgoServerPort is the port where Argo Server listens.
	ArgoServerPort = 2746
)

// Shared cluster state for all E2E tests.
// We maintain separate state for each connection mode to allow parallel testing.
var (
	sharedCluster     *E2ECluster
	sharedClusterOnce sync.Once
	sharedClusterErr  error
)

// GetConnectionMode returns the connection mode from the E2E_MODE environment variable.
// Defaults to ModeKubernetesAPI if not set or invalid.
func GetConnectionMode() ConnectionMode {
	mode := os.Getenv("E2E_MODE")
	switch ConnectionMode(mode) {
	case ModeArgoServer:
		return ModeArgoServer
	case ModeKubernetesAPI:
		return ModeKubernetesAPI
	default:
		return ModeKubernetesAPI
	}
}

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

	// ConnectionMode indicates how the client connects to Argo Workflows.
	ConnectionMode ConnectionMode

	// ArgoServerURL is the URL to the Argo Server (only set in ModeArgoServer).
	ArgoServerURL string

	// portForwardCmd is the kubectl port-forward process (only set in ModeArgoServer).
	portForwardCmd *exec.Cmd
}

// SetupE2ECluster returns the shared E2E cluster, creating it on first call.
// All tests share the same cluster to speed up test execution.
// The cluster is terminated when the test binary exits.
func SetupE2ECluster(ctx context.Context, t *testing.T) *E2ECluster {
	t.Helper()

	sharedClusterOnce.Do(func() {
		sharedCluster, sharedClusterErr = createSharedCluster(ctx, t)
	})

	require.NoError(t, sharedClusterErr, "Failed to create shared E2E cluster")
	require.NotNil(t, sharedCluster, "Shared cluster is nil")

	return sharedCluster
}

// createSharedCluster creates the shared k3s cluster with Argo Workflows.
func createSharedCluster(ctx context.Context, t *testing.T) (*E2ECluster, error) {
	mode := GetConnectionMode()
	t.Logf("Starting shared k3s container for all E2E tests (mode: %s)...", mode)

	// Start k3s container
	k3sContainer, err := k3s.Run(ctx, "rancher/k3s:v1.31.2-k3s1")
	if err != nil {
		return nil, fmt.Errorf("failed to start k3s container: %w", err)
	}

	// Get kubeconfig from container
	kubeconfig, err := k3sContainer.GetKubeConfig(ctx)
	if err != nil {
		_ = k3sContainer.Terminate(context.Background())
		return nil, fmt.Errorf("failed to get kubeconfig from k3s: %w", err)
	}

	// Write kubeconfig to temp file
	kubeconfigFile, err := os.CreateTemp("", "e2e-kubeconfig-*.yaml")
	if err != nil {
		_ = k3sContainer.Terminate(context.Background())
		return nil, fmt.Errorf("failed to create temp kubeconfig file: %w", err)
	}

	kubeconfigPath := kubeconfigFile.Name()

	_, err = kubeconfigFile.Write(kubeconfig)
	if err != nil {
		_ = os.Remove(kubeconfigPath)
		_ = k3sContainer.Terminate(context.Background())
		return nil, fmt.Errorf("failed to write kubeconfig: %w", err)
	}
	err = kubeconfigFile.Close()
	if err != nil {
		_ = os.Remove(kubeconfigPath)
		_ = k3sContainer.Terminate(context.Background())
		return nil, fmt.Errorf("failed to close kubeconfig file: %w", err)
	}

	t.Logf("Kubeconfig written to: %s", kubeconfigPath)

	// Install Argo Workflows
	t.Log("Installing Argo Workflows...")
	if err := installArgoWorkflowsShared(t, kubeconfigPath); err != nil {
		_ = os.Remove(kubeconfigPath)
		_ = k3sContainer.Terminate(context.Background())
		return nil, fmt.Errorf("failed to install Argo Workflows: %w", err)
	}

	// Wait for Argo controller to be ready
	t.Log("Waiting for Argo controller to be ready...")
	if err := waitForArgoControllerShared(t, kubeconfigPath); err != nil {
		_ = os.Remove(kubeconfigPath)
		_ = k3sContainer.Terminate(context.Background())
		return nil, fmt.Errorf("argo controller not ready: %w", err)
	}

	cluster := &E2ECluster{
		Kubeconfig:     string(kubeconfig),
		KubeconfigPath: kubeconfigPath,
		ArgoNamespace:  ArgoNamespace,
		container:      k3sContainer,
		ConnectionMode: mode,
	}

	// Set up connection based on mode
	if mode == ModeArgoServer {
		// Wait for Argo Server to be ready
		t.Log("Waiting for Argo Server to be ready...")
		if err := waitForArgoServerShared(t, kubeconfigPath); err != nil {
			_ = os.Remove(kubeconfigPath)
			_ = k3sContainer.Terminate(context.Background())
			return nil, fmt.Errorf("argo server not ready: %w", err)
		}

		// Start port-forward to Argo Server
		t.Log("Starting port-forward to Argo Server...")
		portForwardCmd, localPort, err := startPortForward(t, kubeconfigPath)
		if err != nil {
			_ = os.Remove(kubeconfigPath)
			_ = k3sContainer.Terminate(context.Background())
			return nil, fmt.Errorf("failed to start port-forward: %w", err)
		}

		cluster.portForwardCmd = portForwardCmd
		cluster.ArgoServerURL = fmt.Sprintf("localhost:%d", localPort)

		t.Logf("Argo Server available at: %s", cluster.ArgoServerURL)

		// Create Argo client with server mode
		argoClient, err := argo.NewClient(ctx, &argo.Config{
			ArgoServer:         cluster.ArgoServerURL,
			Namespace:          ArgoNamespace,
			Secure:             false, // Local port-forward uses HTTP
			InsecureSkipVerify: true,
		})
		if err != nil {
			stopPortForward(portForwardCmd)
			_ = os.Remove(kubeconfigPath)
			_ = k3sContainer.Terminate(context.Background())
			return nil, fmt.Errorf("failed to create Argo client (server mode): %w", err)
		}
		cluster.ArgoClient = argoClient
	} else {
		// Direct Kubernetes API mode
		argoClient, err := argo.NewClient(ctx, &argo.Config{
			Kubeconfig: kubeconfigPath,
			Namespace:  ArgoNamespace,
		})
		if err != nil {
			_ = os.Remove(kubeconfigPath)
			_ = k3sContainer.Terminate(context.Background())
			return nil, fmt.Errorf("failed to create Argo client (kubernetes mode): %w", err)
		}
		cluster.ArgoClient = argoClient
	}

	t.Logf("Shared E2E cluster setup complete (mode: %s)", mode)

	return cluster, nil
}

// installArgoWorkflowsShared installs Argo Workflows in the k3s cluster (non-fatal version).
func installArgoWorkflowsShared(t *testing.T, kubeconfigPath string) error {
	t.Helper()

	// Create the argo namespace first (quick-start manifest expects it to exist)
	//nolint:gosec // Using kubectl in tests is expected
	nsCmd := exec.Command("kubectl", "create", "namespace", ArgoNamespace)
	nsCmd.Env = append(os.Environ(), "KUBECONFIG="+kubeconfigPath)
	output, err := nsCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create argo namespace: %s: %w", string(output), err)
	}

	// Download the quick-start manifest to a temp file
	manifestFile, err := os.CreateTemp("", "argo-install-*.yaml")
	if err != nil {
		return fmt.Errorf("failed to create temp manifest file: %w", err)
	}
	manifestPath := manifestFile.Name()

	// Close file before curl writes to it
	err = manifestFile.Close()
	if err != nil {
		return fmt.Errorf("failed to close manifest file: %w", err)
	}

	defer func() {
		_ = os.Remove(manifestPath) //nolint:errcheck // Cleanup is best-effort
	}()

	// Download manifest
	//nolint:gosec // Using curl to download manifests in tests is expected
	downloadCmd := exec.Command("curl", "-sSL", "-o", manifestPath, ArgoQuickStartURL)
	output, err = downloadCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to download Argo quick-start manifest: %s: %w", string(output), err)
	}

	// Apply the manifest
	//nolint:gosec // Using kubectl in tests is expected
	applyCmd := exec.Command("kubectl", "apply", "-f", manifestPath)
	applyCmd.Env = append(os.Environ(), "KUBECONFIG="+kubeconfigPath)
	output, err = applyCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to apply Argo manifest: %s: %w", string(output), err)
	}

	t.Logf("Argo Workflows %s installed", ArgoVersion)
	return nil
}

// waitForArgoControllerShared waits for the Argo controller deployment to be ready (non-fatal version).
func waitForArgoControllerShared(t *testing.T, kubeconfigPath string) error {
	t.Helper()

	// Wait for the argo namespace to exist
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for Argo controller to be ready")
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
				return nil
			}

			// Log the error but continue waiting
			t.Logf("Waiting for Argo controller... %s", string(output))
		}
	}
}

// waitForArgoServerShared waits for the Argo Server deployment to be ready (non-fatal version).
func waitForArgoServerShared(t *testing.T, kubeconfigPath string) error {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for Argo Server to be ready")
		case <-ticker.C:
			// Check if the deployment is ready
			//nolint:gosec // Using kubectl in tests is expected
			cmd := exec.Command("kubectl", "wait", "--for=condition=available",
				"--timeout=5s",
				"-n", ArgoNamespace,
				"deployment/argo-server")
			cmd.Env = append(os.Environ(), "KUBECONFIG="+kubeconfigPath)
			output, err := cmd.CombinedOutput()

			if err == nil {
				t.Log("Argo Server is ready")
				return nil
			}

			// Log the error but continue waiting
			t.Logf("Waiting for Argo Server... %s", string(output))
		}
	}
}

// startPortForward starts a kubectl port-forward to the Argo Server and returns the command and local port.
func startPortForward(t *testing.T, kubeconfigPath string) (*exec.Cmd, int, error) {
	t.Helper()

	// Use a fixed local port to avoid conflicts - we'll use 0 and parse the output
	// But kubectl port-forward doesn't easily give us the port, so we use a known port
	localPort := ArgoServerPort

	//nolint:gosec // Using kubectl in tests is expected
	cmd := exec.Command("kubectl", "port-forward",
		"-n", ArgoNamespace,
		"svc/argo-server",
		fmt.Sprintf("%d:%d", localPort, ArgoServerPort))
	cmd.Env = append(os.Environ(), "KUBECONFIG="+kubeconfigPath)

	// Start the port-forward in the background
	if err := cmd.Start(); err != nil {
		return nil, 0, fmt.Errorf("failed to start port-forward: %w", err)
	}

	// Wait for the port-forward to be ready by checking if we can connect
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			stopPortForward(cmd)
			return nil, 0, fmt.Errorf("timeout waiting for port-forward to be ready")
		case <-ticker.C:
			// Try to connect to the port to verify it's ready
			//nolint:gosec // Using curl in tests is expected
			checkCmd := exec.Command("curl", "-s", "-o", "/dev/null", "-w", "%{http_code}",
				fmt.Sprintf("http://localhost:%d/api/v1/info", localPort))
			output, err := checkCmd.CombinedOutput()
			if err == nil && (string(output) == "200" || string(output) == "401") {
				// 200 or 401 means the server is responding
				t.Logf("Port-forward is ready (HTTP status: %s)", string(output))
				return cmd, localPort, nil
			}
			t.Logf("Waiting for port-forward... (status: %s, err: %v)", string(output), err)
		}
	}
}

// stopPortForward stops the kubectl port-forward process.
func stopPortForward(cmd *exec.Cmd) {
	if cmd != nil && cmd.Process != nil {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
	}
}

// installArgoWorkflows installs Argo Workflows in the k3s cluster.
func installArgoWorkflows(t *testing.T, kubeconfigPath string) {
	t.Helper()
	err := installArgoWorkflowsShared(t, kubeconfigPath)
	require.NoError(t, err, "Failed to install Argo Workflows")
}

// waitForArgoController waits for the Argo controller deployment to be ready.
func waitForArgoController(t *testing.T, kubeconfigPath string) {
	t.Helper()
	err := waitForArgoControllerShared(t, kubeconfigPath)
	require.NoError(t, err, "Argo controller not ready")
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

// ClusterWorkflowTemplateExists checks if a cluster workflow template exists.
func (c *E2ECluster) ClusterWorkflowTemplateExists(t *testing.T, name string) bool {
	t.Helper()

	cmd := exec.Command("kubectl", "get", "clusterworkflowtemplate", name)
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
