.PHONY: build run test clean install lint deps build-all release

APP_NAME=etcdtui
BUILD_DIR=bin
MAIN_PATH=./cmd/etcdtui

# Version info (can be overridden)
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build flags
LDFLAGS=-ldflags "-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildTime=$(BUILD_TIME)"

# Platforms for cross-compilation
PLATFORMS=darwin/amd64 darwin/arm64 linux/amd64 linux/arm64

build:
	@echo "Building $(APP_NAME) $(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(APP_NAME)"

run:
	@echo "Running $(APP_NAME)..."
	TERM=xterm-256color go run $(MAIN_PATH)

test:
	@echo "Running tests..."
	go test -v ./...

clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	rm -rf dist
	@echo "Clean complete"

install:
	@echo "Installing $(APP_NAME)..."
	go install $(LDFLAGS) $(MAIN_PATH)
	@echo "Install complete"

lint:
	@echo "Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install it from https://golangci-lint.run/usage/install/"; \
	fi

deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy
	@echo "Dependencies updated"

# Build for all platforms
build-all:
	@echo "Building $(APP_NAME) $(VERSION) for all platforms..."
	@mkdir -p $(BUILD_DIR)
	@for platform in $(PLATFORMS); do \
		GOOS=$${platform%/*} GOARCH=$${platform#*/} \
		go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-$${platform%/*}-$${platform#*/} $(MAIN_PATH) && \
		echo "Built: $(BUILD_DIR)/$(APP_NAME)-$${platform%/*}-$${platform#*/}"; \
	done
	@echo "All builds complete"

# Create release archives
release: build-all
	@echo "Creating release archives..."
	@mkdir -p $(BUILD_DIR)/release
	@for platform in $(PLATFORMS); do \
		os=$${platform%/*}; \
		arch=$${platform#*/}; \
		name=$(APP_NAME)-$(VERSION)-$${os}-$${arch}; \
		mkdir -p $(BUILD_DIR)/release/$$name; \
		cp $(BUILD_DIR)/$(APP_NAME)-$${os}-$${arch} $(BUILD_DIR)/release/$$name/$(APP_NAME); \
		cp README.md LICENSE $(BUILD_DIR)/release/$$name/ 2>/dev/null || true; \
		cd $(BUILD_DIR)/release && tar -czf $$name.tar.gz $$name && rm -rf $$name; \
		cd ../..; \
		echo "Created: $(BUILD_DIR)/release/$$name.tar.gz"; \
	done
	@cd $(BUILD_DIR)/release && sha256sum *.tar.gz > checksums.txt 2>/dev/null || shasum -a 256 *.tar.gz > checksums.txt
	@echo "Release archives created in $(BUILD_DIR)/release/"

help:
	@echo "Available targets:"
	@echo "  build      - Build the application"
	@echo "  build-all  - Build for all platforms (darwin/linux, amd64/arm64)"
	@echo "  run        - Run the application without building"
	@echo "  test       - Run tests"
	@echo "  clean      - Remove build artifacts"
	@echo "  install    - Install the application to GOPATH/bin"
	@echo "  lint       - Run golangci-lint"
	@echo "  deps       - Download and tidy dependencies"
	@echo "  release    - Create release archives for all platforms"
	@echo "  help       - Show this help message"
	@echo ""
	@echo "Version info:"
	@echo "  VERSION=$(VERSION)"
	@echo "  COMMIT=$(COMMIT)"
