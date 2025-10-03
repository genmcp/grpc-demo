# gRPC to HTTP/JSON Proxy Demo

This project demonstrates how to create a RESTful HTTP/JSON proxy for an existing gRPC service using [gRPC-Gateway](https://github.com/grpc-ecosystem/grpc-gateway).

The key idea is to show that you can add a RESTful API front-end to a gRPC service **without modifying the original service's code**. The only requirement is having access to the service's `.proto` definition files.

## Project Structure

This repository is split into distinct parts to simulate a real-world scenario where the gRPC service might be maintained by a different team or exist as a legacy component.

-   `000_grpc_server/`: A standalone gRPC service for managing feature requests. This represents the existing service that we don't want to (or can't) modify.
-   `001_grpc_proxy/`: An HTTP/JSON proxy that translates RESTful API calls into gRPC requests and forwards them to the gRPC service.
-   `002_genmcp_grpc/`: Contains the configuration to expose the HTTP proxy endpoints as tools for AI models using `gen-mcp`.

## How it Works

1.  The **gRPC Server** (`000_grpc_server`) is a standard Go gRPC application. It defines its service and messages using a `features.proto` file and implements the service logic.
2.  The **gRPC Proxy** (`001_grpc_proxy`) uses the *same* `features.proto` file to generate not only the client stubs but also a reverse proxy server.
3.  We add annotations to a configuration file (`features.gprc.proxy.config.yaml`) to map HTTP methods and URL paths to the gRPC service's RPCs.
4.  `protoc` with the `protoc-gen-grpc-gateway` plugin reads the `.proto` file and the YAML configuration to generate a Go module (`pkg/features.pb.gw.go`) that handles the HTTP-to-gRPC translation.
5.  The proxy's `main.go` simply starts an HTTP server and registers this generated gateway handler.

## Getting Started

To run this demo, you will need to start the gRPC server, the HTTP proxy, and finally the `gen-mcp` server.

1.  **Start the gRPC Server:**
    -   Navigate to the `000_grpc_server` directory.
    -   Follow the instructions in `000_grpc_server/README.md`.

2.  **Start the HTTP Proxy:**
    -   Navigate to the `001_grpc_proxy` directory.
    -   Follow the instructions in `001_grpc_proxy/README.md`.

3.  **Expose the Service with gen-mcp:**
    -   Navigate to the `002_genmcp_grpc` directory.
    -   Follow the instructions in `002_genmcp_grpc/README.md`.

Once all three are running, you can make gRPC calls to port `50051`, RESTful API calls to port `9090`, and MCP tool calls to port `8080`.

TODO:
- Diagrams
- Remove TODOs from code
