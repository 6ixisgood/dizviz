.PHONY: all proto agent controlplane clean deps dev help

# Load user-specific configuration if it exists
-include Makefile.local


export PATH := $(PATH):$(HOME)/go/bin

# Default target
all: agent controlplane

# Generate protobuf code
proto:
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/agent.proto

# Build binaries
agent:
	@go build -o bin/agent cmd/agent/main.go

controlplane:
	@go build -o bin/controlplane cmd/controlplane/*.go

# Development mode (no build needed, uses stub display)
dev:
	@go run cmd/agent/main.go $(ARGS)

# Clean
clean:
	@rm -rf bin/ proto/*.pb.go

# Install dependencies
deps:
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Help
help:
	@echo "DizViz Makefile"
	@echo ""
	@echo "Targets:"
	@echo "  all          - Build agent and control plane (default)"
	@echo "  agent        - Build agent binary"
	@echo "  controlplane - Build control plane server binary"
	@echo "  proto        - Generate protobuf code"
	@echo "  dev          - Run agent without building (use ARGS for custom flags)"
	@echo "  clean        - Remove build artifacts"
	@echo "  deps         - Install Go dependencies"
	@echo ""
	@echo "Configuration:"
	@echo "  Create Makefile.local to customize CGO flags and paths"
	@echo "  See Makefile.local.example for template"
	@echo ""
	@echo "Examples:"
	@echo "  make dev"
	@echo "  make dev ARGS='--config=custom.yaml --name=test'"
