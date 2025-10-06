# gRPC to HTTP/JSON Proxy Demo

This project demonstrates how to create a RESTful HTTP/JSON proxy for an existing gRPC service using [gRPC-Gateway](https://github.com/grpc-ecosystem/grpc-gateway) and subsequently expose it as a set of tools for AI models using [gen-mcp](https://github.com/genmcp/gen-mcp).

The key objective is to show that you can add a RESTful API front-end to a gRPC service **without modifying the original service's code**. The only requirement is access to the service's `.proto` definition files.

## Technology Stack

*   **Go**: The language used for both the gRPC server and the HTTP proxy.
*   **gRPC**: The primary RPC framework for the core service.
*   **Protocol Buffers**: The interface definition language for the gRPC service.
*   **gRPC-Gateway**: The tool used to generate the reverse proxy that translates HTTP/JSON to gRPC.
*   **gen-mcp**: A tool for exposing APIs to AI models using the Model-Context Protocol.

## Architecture

The demo consists of three independent services that run simultaneously:

1.  **gRPC Server**: The original, "untouchable" service.
2.  **HTTP Proxy**: The gRPC-Gateway server that exposes a RESTful API.
3.  **MCP Server**: The `gen-mcp` server that exposes the RESTful API as tools for AI.

The flow of a request from an AI model is as follows:

```mermaid
sequenceDiagram
    participant AI Model
    participant MCP Server (gen-mcp)
    participant HTTP Proxy (gRPC-Gateway)
    participant gRPC Server (Core Service)

    AI Model->>MCP Server (gen-mcp): Tool Call (e.g., get_top_feature)
    MCP Server (gen-mcp)->>HTTP Proxy (gRPC-Gateway): HTTP GET /top_feature
    HTTP Proxy (gRPC-Gateway)->>gRPC Server (Core Service): gRPC Call GetTopFeature()
    gRPC Server (Core Service)-->>HTTP Proxy (gRPC-Gateway): gRPC Response
    HTTP Proxy (gRPC-Gateway)-->>MCP Server (gen-mcp): HTTP Response
    MCP Server (gen-mcp)-->>AI Model: Tool Output
```

## Project Structure

This repository is split into distinct parts to simulate a real-world scenario where the gRPC service might be maintained by a different team or exist as a legacy component.

-   `000_grpc_server/`: A standalone gRPC service for managing feature requests. This represents the existing service that we don't want to (or can't) modify.
-   `001_grpc_proxy/`: An HTTP/JSON proxy that translates RESTful API calls into gRPC requests and forwards them to the gRPC service.
-   `002_genmcp_grpc/`: Contains the configuration to expose the HTTP proxy endpoints as tools for AI models using `gen-mcp`.

## Getting Started

To run the full demo, you will need to start the gRPC server, the HTTP proxy, and finally the `gen-mcp` server in separate terminal sessions.

1.  **Start the gRPC Server:**
    *   Navigate to the `000_grpc_server` directory.
    *   Follow the instructions in `000_grpc_server/README.md`.

2.  **Start the HTTP Proxy:**
    *   Navigate to the `001_grpc_proxy` directory.
    *   Follow the instructions in `001_grpc_proxy/README.md`.

3.  **Expose the Service with gen-mcp:**
    *   Navigate to the `002_genmcp_grpc` directory.
    *   Follow the instructions in `002_genmcp_grpc/README.md`.

Once all three services are running, you can make gRPC calls to port `50051`, RESTful API calls to port `9090`, and MCP tool calls to port `8080` at http://localhost:8080/mcp.
