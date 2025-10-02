# Feature Request gRPC Service

This directory contains a gRPC version of the Feature Request API (https://github.com/genmcp/gen-mcp/tree/main/examples/http-conversion), with feature parity to the HTTP server.

## Prerequisites

- Go (1.19+)
- Protocol Buffer Compiler (`protoc`), version 3+
- Go plugins for the protocol compiler:
  ```bash
  go install google.golang.org/protobuf/cmd/protoc-gen-go
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
  ```
- Ensure your `$GOPATH/bin` directory is in your `PATH`.
- `grpcurl` for testing (optional).

## Getting Started

### 1. Generate Go code from Protobuf

From this directory run the following command to generate the gRPC server and client stubs:

```bash
protoc --go_out=./pkg \ 
       --go_opt=paths=source_relative \
       --go-grpc_out=./pkg \ 
       --go-grpc_opt=paths=source_relative \
       features.proto
```

This will create `pkg/features.pb.go` and `pkg/features_grpc.pb.go`.

### 2. Tidy Dependencies

Download the necessary Go modules:

```bash
go mod tidy
```

### 3. Run the gRPC Server

```bash
go run cmd/main.go
```

The gRPC server will start and listen on port `50051`.

### 4. Test the Service (using `grpcurl`)

You can test the running server using a tool like `grpcurl`.

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
grpcurl -plaintext -d '{"title": "gRPC Support", "description": "Add gRPC endpoints for all services.", "details": "Use protobuf to define the service."}' localhost:50051 features.FeatureService/AddFeature
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
