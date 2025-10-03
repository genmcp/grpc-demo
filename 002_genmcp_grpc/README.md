# Exposing the gRPC Service with gen-mcp

This directory contains the configuration to expose the gRPC service (via its HTTP proxy) as a set of tools that AI models can consume using [gen-mcp](https://github.com/genmcp/gen-mcp).

The `mcpfile.yaml` in this directory defines how HTTP endpoints from the proxy are mapped to MCP tools.

## Prerequisites

- The gRPC server from `../000_grpc_server` **must be running**.
- The HTTP proxy from `../001_grpc_proxy` **must be running**.
- The `gen-mcp` CLI must be installed.

## Getting Started

### 1. (Optional) Generate the `mcpfile.yaml`

This directory already includes a pre-configured `mcpfile.yaml`. However, if you wanted to generate it yourself from the proxy's OpenAPI specification, you would run the following command from this directory:

```bash
# The -H flag is required because the swagger file does not contain a server URL.
genmcp convert ../001_grpc_proxy/features.swagger.json \
  -H localhost:9090 \
  -o mcpfile.yaml
```

This command reads the swagger file from the proxy directory, sets the base URL for the API calls, and generates the `mcpfile.yaml` in the current directory. You can then customize the generated file to improve tool descriptions, add examples, or change which endpoints are exposed.

### 2. Run the MCP Server

With the gRPC server and the HTTP proxy running, start the `gen-mcp` server from this directory:

```bash
genmcp run -f mcpfile.yaml
```

The MCP server will start on port `8080` (as configured in `mcpfile.yaml`). Now, an AI assistant or any MCP-compatible client can connect to `http://localhost:8080/mcp`. When a tool is invoked, `gen-mcp` will make an HTTP call to the proxy, which in turn calls the original gRPC service.

## 3. Testing the MCP Service

You can test the MCP server using any MCP client.

Using https://github.com/modelcontextprotocol/inspector
```shell
nvm use 22
npx @modelcontextprotocol/inspector
# localhost:6274 will be opened on your browser 
# In the UI, connect to the MCP server:
# - Transport Type: Streamable HTTP
# - MCP URL: http://localhost:8080/mcp
```
