# Docker Compose Setup

This directory contains a Docker Compose configuration for running MCP for Argo Workflows as a remote HTTP/SSE server.

## Overview

Docker Compose provides a simple way to run the MCP server locally in HTTP/SSE mode for testing or development purposes.

## Prerequisites

1. [Docker](https://docs.docker.com/get-docker/) installed
2. [Docker Compose](https://docs.docker.com/compose/install/) installed
3. Access to an Argo Workflows installation (via kubeconfig or Argo Server)

## Quick Start

### Using Direct Kubernetes API

If you have a kubeconfig file:

```bash
# Copy your kubeconfig
cp ~/.kube/config ./kubeconfig

# Start the server
docker compose up -d

# View logs
docker compose logs -f
```

### Using Argo Server

Edit the `.env` file or set environment variables:

```bash
# Create .env file
cat > .env << EOF
ARGO_SERVER=argo-server.argo:2746
ARGO_TOKEN=Bearer eyJhbGciOiJSUzI1NiIs...
ARGO_NAMESPACE=argo
EOF

# Start the server
docker compose --profile argo-server up -d
```

## Configuration

### Environment Variables

Create a `.env` file or set these variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `MCP_HTTP_ADDR` | HTTP listen address | `:8080` |
| `ARGO_NAMESPACE` | Default namespace | `argo` |
| `ARGO_SERVER` | Argo Server address | - |
| `ARGO_TOKEN` | Auth token for Argo Server | - |
| `ARGO_SECURE` | Use TLS for Argo Server | `true` |
| `ARGO_INSECURE_SKIP_VERIFY` | Skip TLS verification | `false` |

### Mounting kubeconfig

The docker-compose.yaml mounts `./kubeconfig` as the kubeconfig file. Ensure this file exists:

```bash
cp ~/.kube/config ./kubeconfig
```

## Usage

### Start the Server

```bash
docker compose up -d
```

### Check Status

```bash
docker compose ps
docker compose logs -f
```

### Stop the Server

```bash
docker compose down
```

### Connect MCP Clients

The server listens on `http://localhost:8080` by default. Configure your MCP client to connect to this address.

## Example: Connecting with Claude Code

When using Claude Code with a remote MCP server running in Docker:

1. Start the Docker Compose setup
2. Configure Claude Code to use the HTTP endpoint

```json
{
  "mcpServers": {
    "argo-workflows": {
      "url": "http://localhost:8080"
    }
  }
}
```

## Troubleshooting

### Container not starting

```bash
# Check container status
docker compose ps

# View logs
docker compose logs mcp-for-argo-workflows
```

### Connection refused

Ensure the port 8080 is not in use:

```bash
lsof -i :8080
```

### kubeconfig errors

Verify your kubeconfig is valid:

```bash
KUBECONFIG=./kubeconfig kubectl cluster-info
```

### Argo Server connection issues

Test connectivity to Argo Server:

```bash
docker compose exec mcp-for-argo-workflows wget -qO- http://argo-server.argo:2746/api/v1/info
```
