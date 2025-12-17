# Cursor Configuration

This directory contains example configurations for using MCP for Argo Workflows with [Cursor](https://cursor.sh/).

## Configuration Location

Cursor reads MCP server configurations from its settings. You can configure MCP servers through:

1. **Settings UI**: Open Cursor Settings > MCP
2. **JSON Configuration**: Edit Cursor's settings JSON file

## Example Configurations

### Using Direct Kubernetes API

Copy the contents of [`mcp-kubernetes.json`](mcp-kubernetes.json) to your Cursor MCP settings:

```json
{
  "mcp.servers": {
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

Copy the contents of [`mcp-argo-server.json`](mcp-argo-server.json) to your Cursor MCP settings:

```json
{
  "mcp.servers": {
    "argo-workflows": {
      "command": "/usr/local/bin/mcp-for-argo-workflows",
      "args": [
        "--argo-server", "argo-server.argo:2746",
        "--namespace", "argo"
      ],
      "env": {
        "ARGO_TOKEN": "Bearer <your-token-here>"
      }
    }
  }
}
```

This configuration:
- Connects via Argo Server (required for archive operations)
- Uses token-based authentication
- Supports large workflows and workflow archive

## Getting an Argo Token

To generate a token for Argo Server authentication:

```bash
# Using kubectl create token (Kubernetes 1.24+)
export ARGO_TOKEN="Bearer $(kubectl create token argo-server -n argo)"

# Or from an existing secret
export ARGO_TOKEN="Bearer $(kubectl get secret -n argo argo-server-token -o jsonpath='{.data.token}' | base64 -d)"
```

## Environment Variables

You can configure the server using environment variables in the `env` block:

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

### Server not appearing in Cursor

- Restart Cursor after changing MCP configuration
- Check Cursor's Developer Tools console for errors

## Verifying the Setup

After configuring Cursor:

1. Restart Cursor
2. Open a project
3. Ask Cursor AI: "List all workflows in the argo namespace"

Cursor should be able to query your Argo Workflows installation.
