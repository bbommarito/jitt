.PHONY: build test clean install lint fmt help

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the binary
	go build -v -o jitt ./cmd/jitt

install: ## Install the binary to GOPATH/bin
	go install ./cmd/jitt

test: ## Run tests
	go test -v -race ./...

test-coverage: ## Run tests with coverage
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint: ## Run golangci-lint
	golangci-lint run

fmt: ## Format code
	go fmt ./...
	goimports -w .

clean: ## Clean build artifacts
	rm -f jitt
	rm -rf dist/
	rm -f coverage.out coverage.html

build-all: ## Build for all platforms
	mkdir -p dist
	# Linux
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o dist/jitt-linux-amd64 ./cmd/jitt
	GOOS=linux GOARCH=arm64 go build -ldflags "-s -w" -o dist/jitt-linux-arm64 ./cmd/jitt
	# macOS
	GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o dist/jitt-darwin-amd64 ./cmd/jitt
	GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" -o dist/jitt-darwin-arm64 ./cmd/jitt
	# Windows
	GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o dist/jitt-windows-amd64.exe ./cmd/jitt

ginkgo: ## Run tests with Ginkgo (requires ginkgo CLI)
	ginkgo -r -v

dev-setup: ## Set up development environment
	@echo "Installing development dependencies..."
	go install github.com/onsi/ginkgo/v2/ginkgo@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	@echo "Development environment ready!"
