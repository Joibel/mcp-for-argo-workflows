# MCP for Argo Workflows

MCP (Model Context Protocol) server for [Argo Workflows](https://argoproj.github.io/argo-workflows/), enabling AI assistants like Claude to interact with Argo Workflows via standardized tools.

## What is MCP?

The [Model Context Protocol](https://modelcontextprotocol.io/) is an open standard that allows AI assistants to securely interact with external tools and data sources. This server exposes Argo Workflows operations as MCP tools, enabling AI assistants to:

- Submit and manage workflows
- Monitor workflow status and logs
- Manage workflow templates and cron workflows
- Query the workflow archive

## Features

### Connection Modes

- **Direct Kubernetes API** — Connect directly to the Kubernetes API using kubeconfig
- **Argo Server** — Connect via Argo Server for full feature support including workflow archive

### Transport Modes

- **stdio** — For local clients like Claude Desktop and Cursor
- **HTTP/SSE** — For remote client connections

### Supported Operations

*Coming soon - see the Linear project for implementation status*

- Workflow lifecycle: submit, list, get, delete, logs, watch, wait
- Workflow control: suspend, resume, stop, terminate, retry, resubmit
- Validation: lint workflow manifests
- WorkflowTemplates: list, get, create, delete
- ClusterWorkflowTemplates: list, get, create, delete
- CronWorkflows: list, get, create, delete, suspend, resume
- Archived workflows: list, get, delete, resubmit, retry (Argo Server only)
- Node operations: get, set

## Quick Start

*Coming soon*

```bash
# Download the binary
# Configure your MCP client to use the server
```

## Configuration

| Environment Variable | CLI Flag | Description |
|---------------------|----------|-------------|
| `ARGO_SERVER` | `--argo-server` | Argo Server host:port (omit for direct K8s) |
| `ARGO_TOKEN` | `--argo-token` | Bearer token for Argo Server auth |
| `ARGO_NAMESPACE` | `--namespace` | Default namespace for operations |
| `MCP_TRANSPORT` | `--transport` | `stdio` (default) or `http` |
| `MCP_HTTP_ADDR` | `--http-addr` | HTTP listen address (default `:8080`) |
| `KUBECONFIG` | `--kubeconfig` | Path to kubeconfig (when not using Argo Server) |

## Building from Source

```bash
# Clone the repository
git clone https://github.com/Joibel/mcp-for-argo-workflows.git
cd mcp-for-argo-workflows

# Build the binary
make build

# Run tests
make test

# Run linter
make lint
```

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.

### Development

```bash
# Install development tools
make tools

# Run all checks
make all
```

## License

Apache License 2.0 - see [LICENSE](LICENSE) for details.
