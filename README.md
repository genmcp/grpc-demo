# Feature Request gRPC Service

This directory contains a gRPC version of the Feature Request API (https://github.com/genmcp/gen-mcp/tree/main/examples/http-conversion), with feature parity to the HTTP server.

## Prerequisites

- Go (1.19+)
- Protocol Buffer Compiler (`protoc`), version 3+
- Go plugins for the protocol compiler:
  ```bash
  go install google.golang.org/protobuf/cmd/protoc-gen-go
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
  go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
  go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
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
go run cmd/grpc/main.go
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

### 5. Create the HTTP Gateway

**Generate the grpc-gateway stubs:**
```shell
protoc -I . \
--grpc-gateway_out ./pkg \
--grpc-gateway_opt paths=source_relative \
--grpc-gateway_opt grpc_api_configuration=./features.config.yaml \
features.proto
```

**Generate the OpenAPI (Swagger) specification:**
```shell
protoc -I . \
--openapiv2_out ./ \
--openapiv2_opt grpc_api_configuration=./features.config.yaml \
features.proto
```

**Run the HTTP Gateway:**
```bash
go run cmd/proxy/main.go
```

## 6. Test the Service (using `curl`)

**List all features:**
```bash
curl -s localhost:8081/features
```

**Get a specific feature:**
```bash
curl -s localhost:8081/features/1
```

**Get the top feature:**
```bash
curl -s localhost:8081/top_feature
```

**Add a new feature:**
```bash
curl -s -X POST -d '{"title": "HTTP Gateway Support", "description": "Access gRPC via REST.", "details": "Use grpc-gateway."}' localhost:8081/features
```

**Vote for a feature:**
```bash
curl -s -X POST -d '{"feature_id": 1}' localhost:8081/features/vote
```

**Mark a feature as complete:**
```bash
curl -s -X POST -d '{"feature_id": 2}' localhost:8081/features/complete
```

**Delete a feature:**
```bash
curl -s -X DELETE localhost:8081/features/3
```


TODO:
- Actually, let's move the GRPC server to another repository - or another directory in this repo.
  - Our point is that you don't have to change your proto files or your GRPC server code to get HTTP endpoints.
