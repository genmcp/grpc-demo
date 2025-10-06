# gRPC to HTTP Proxy (`001_grpc_proxy`)

This directory contains an HTTP/JSON reverse proxy for the gRPC Feature Request service.

The purpose of this example is to demonstrate how to create an HTTP gateway for an existing gRPC server (`000_grpc_server`) using **gRPC-Gateway**.

The key takeaway is that we can achieve this without modifying the original server's code; we only need its `.proto` definition file.

## Prerequisites

-   The gRPC server from the `../000_grpc_server` directory **must be running**.
-   Go (1.24+)
-   Protocol Buffer Compiler (`protoc`), version 3+
-   Go plugins for the protocol compiler:
    ```bash
    # Ensure dependencies are available for the tools
    go mod tidy
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
    ```
-   Ensure your `$GOPATH/bin` directory is in your `PATH`.
-   `curl` for testing.

## Getting Started

### 1. Generate Go Code from Protobuf

First, copy the service definition from the gRPC server.
```bash
cp ../000_grpc_server/features.proto .
```

Next, run `protoc` to generate the necessary Go code and the OpenAPI specification.

```bash
# 1. Generate base gRPC client stubs (features.pb.go, features_grpc.pb.go)
protoc -I . \
       --go_out=./pkg \
       --go_opt=paths=source_relative \
       --go-grpc_out=./pkg \
       --go-grpc_opt=paths=source_relative \
       ./features.proto

# 2. Generate the gRPC-Gateway proxy (features.pb.gw.go)
# This uses features.gprc.proxy.config.yaml to map HTTP routes to gRPC RPCs.
protoc -I . \
       --grpc-gateway_out ./pkg \
       --grpc-gateway_opt paths=source_relative \
       --grpc-gateway_opt grpc_api_configuration=./features.gprc.proxy.config.yaml \
       ./features.proto

# 3. Generate the OpenAPI v2 (Swagger) specification (features.swagger.json)
# This uses the same proxy config plus features.openapi.config.yaml for metadata.
protoc -I . \
       --openapiv2_out ./ \
       --openapiv2_opt grpc_api_configuration=./features.gprc.proxy.config.yaml \
       --openapiv2_opt openapi_configuration=./features.openapi.config.yaml \
       features.proto
```

These commands create `pkg/features.pb.go`, `pkg/features_grpc.pb.go`, `pkg/features.pb.gw.go`, and `features.swagger.json`.

### 2. Tidy Dependencies

Download the necessary Go modules defined in `go.mod`.

```bash
go mod tidy
```

### 3. Run the HTTP Gateway

Ensure the gRPC server from `../000_grpc_server` is running. Then, start the proxy. It will listen for HTTP requests on port `9090` and forward them to the gRPC server on `localhost:50051`.

```bash
go run main.go
```
You should see the output: `HTTP server listening on :9090`

## Testing the Service (with `curl`)

With the proxy running, you can now interact with the gRPC service using standard HTTP requests.

**List all features:**
```bash
curl -s localhost:9090/features
```

**Get a specific feature:**
```bash
curl -s localhost:9090/features/1
```

**Get the top feature:**
```bash
curl -s localhost:9090/top_feature
```

**Add a new feature:**
```bash
curl -s -X POST -d '{"title": "HTTP Gateway Support", "description": "Access gRPC via REST.", "details": "Use grpc-gateway."}' localhost:9090/features
```

**Vote for a feature:**
```bash
curl -s -X POST -d '{"id": 1}' localhost:9090/features/vote
```

**Mark a feature as complete:**
```bash
curl -s -X POST -d '{"id": 2}' localhost:9090/features/complete
```

**Delete a feature:**
```bash
curl -s -X DELETE localhost:9090/features/3
```
