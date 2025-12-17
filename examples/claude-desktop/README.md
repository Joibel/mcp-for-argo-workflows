# Claude Desktop Configuration

This directory contains example configurations for using MCP for Argo Workflows with [Claude Desktop](https://claude.ai/download).

## Configuration Location

Claude Desktop reads MCP server configurations from:

- **macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
- **Windows**: `%APPDATA%\Claude\claude_desktop_config.json`

## Example Configurations

### Using Direct Kubernetes API (Recommended for Local Development)

Copy the contents of [`config-kubernetes.json`](config-kubernetes.json) to your Claude Desktop configuration:

```json
{
  "mcpServers": {
    "argo-workflows": {
      "command": "/usr/local/bin/mcp-for-argo-workflows",
      "args": ["--namespace", "argo"]
    }
  }
}
```

This configuration:
- Uses your local kubeconfig (`~/.kube/config`)
- Connects directly to the Kubernetes API
- Sets `argo` as the default namespace

### Using Argo Server

Copy the contents of [`config-argo-server.json`](config-argo-server.json) to your Claude Desktop configuration:

```json
{
  "mcpServers": {
    "argo-workflows": {
      "command": "/usr/local/bin/mcp-for-argo-workflows",
      "args": [
        "--argo-server", "argo-server.argo:2746",
        "--namespace", "argo"
      ],
      "env": {
        "ARGO_TOKEN": "Bearer eyJhbGciOiJSUzI1NiIs..."
      }
    }
  }
}
```

This configuration:
- Connects via Argo Server (required for archive operations)
- Uses token-based authentication
- Supports large workflows and workflow archive

### Using Port-Forwarded Argo Server

For development with a port-forwarded Argo Server:

```json
{
  "mcpServers": {
    "argo-workflows": {
      "command": "/usr/local/bin/mcp-for-argo-workflows",
      "args": [
        "--argo-server", "localhost:2746",
        "--argo-insecure-skip-verify",
        "--namespace", "argo"
      ]
    }
  }
}
```

First, start the port-forward in a terminal:

```bash
kubectl port-forward svc/argo-server -n argo 2746:2746
```

## Getting an Argo Token

To generate a token for Argo Server authentication:

```bash
# Using kubectl create token (Kubernetes 1.24+)
export ARGO_TOKEN="Bearer $(kubectl create token argo-server -n argo)"

# Or from an existing secret
export ARGO_TOKEN="Bearer $(kubectl get secret -n argo argo-server-token -o jsonpath='{.data.token}' | base64 -d)"
```

## Environment Variables

You can also configure the server using environment variables in the `env` block:

| Variable | Description |
|----------|-------------|
| `ARGO_SERVER` | Argo Server host:port |
| `ARGO_TOKEN` | Bearer token for authentication |
| `ARGO_NAMESPACE` | Default namespace |
| `ARGO_SECURE` | Use TLS (default: true) |
| `ARGO_INSECURE_SKIP_VERIFY` | Skip TLS verification |
| `KUBECONFIG` | Path to kubeconfig file |

## Troubleshooting

### "Failed to create Argo client"

- Verify your kubeconfig is valid: `kubectl cluster-info`
- Check RBAC permissions: `kubectl auth can-i list workflows -n argo`

### "Unauthorized" errors

- Ensure `ARGO_TOKEN` is set correctly and not expired
- Verify the token has required RBAC permissions

### "Certificate verification failed"

- For development, use `--argo-insecure-skip-verify`
- For production, ensure proper TLS certificates are configured

## Verifying the Setup

After configuring Claude Desktop:

1. Restart Claude Desktop
2. Start a new conversation
3. Ask Claude: "List all workflows in the argo namespace"

Claude should be able to query your Argo Workflows installation.
