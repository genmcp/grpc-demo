# Feature Request gRPC Service

This directory contains a standalone gRPC service for managing feature requests.

This server is intended to represent an existing, "untouchable" gRPC service. The code in this directory implements the core business logic and exposes it via gRPC. The next directory, `001_grpc_proxy`, demonstrates how to add an HTTP/JSON interface to this service without modifying any of the code here.

The service is taken from https://github.com/genmcp/gen-mcp/tree/main/examples/http-conversion, with feature parity to the HTTP server defined there.

## Prerequisites

- Go (1.19+)
- Protocol Buffer Compiler (`protoc`), version 3+
- Go plugins for the protocol compiler:
  ```bash
  # need to run this to resolve dependencies
  go mod tidy
  go install google.golang.org/protobuf/cmd/protoc-gen-go
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
  ```
- Ensure your `$GOPATH/bin` directory is in your `PATH`.
- `grpcurl` for testing (optional but recommended).

## Getting Started

### 1. Generate Go Code from Protobuf

From this directory, run the following command to generate the gRPC server and client stubs from the `.proto` file.

```bash
protoc --go_out=./pkg \
       --go_opt=paths=source_relative \
       --go-grpc_out=./pkg \
       --go-grpc_opt=paths=source_relative \
       features.proto
```

This command reads `features.proto` and creates `pkg/features.pb.go` and `pkg/features_grpc.pb.go`.

### 2. Tidy Dependencies

Download the necessary Go modules defined in `go.mod`.

```bash
go mod tidy
```

### 3. Run the gRPC Server

Start the server. It will listen for gRPC connections on port `50051`.

```bash
go run main.go
```
You should see the output: `gRPC server listening on :50051`

## Testing the Service (with `grpcurl`)

With the server running, you can test its RPCs using a tool like `grpcurl`.

**List all features:**
```bash
grpcurl -plaintext localhost:50051 features.FeatureService/ListFeatures
```

**Get a specific feature:**
```bash
grpcurl -plaintext -d '{"id": 1}' localhost:50051 features.FeatureService/GetFeature
```

**Get the top feature:**
```bash
grpcurl -plaintext localhost:50051 features.FeatureService/GetTopFeature
```

**Add a new feature:**
```bash
grpcurl -plaintext -d '{"title": "gRPC Support", "description": "Add gRPC endpoints.", "details": "Use protobuf."}' localhost:50051 features.FeatureService/AddFeature
```

**Vote for a feature:**
```bash
grpcurl -plaintext -d '{"feature_id": 1}' localhost:50051 features.FeatureService/VoteFeature
```

**Mark a feature as complete:**
```bash
grpcurl -plaintext -d '{"feature_id": 2}' localhost:50051 features.FeatureService/CompleteFeature
```

**Delete a feature:**
```bash
grpcurl -plaintext -d '{"id": 3}' localhost:50051 features.FeatureService/DeleteFeature
```

## 4. Next Steps

Now, proceed to the `001_grpc_proxy` directory to set up the HTTP gateway.
