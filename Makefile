.PHONY: build run test clean install lint deps

APP_NAME=etcdtui
BUILD_DIR=bin
MAIN_PATH=./cmd/etcdtui

build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)
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
	@echo "Clean complete"

install:
	@echo "Installing $(APP_NAME)..."
	go install $(MAIN_PATH)
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

help:
	@echo "Available targets:"
	@echo "  build    - Build the application"
	@echo "  run      - Run the application without building"
	@echo "  test     - Run tests"
	@echo "  clean    - Remove build artifacts"
	@echo "  install  - Install the application to GOPATH/bin"
	@echo "  lint     - Run golangci-lint"
	@echo "  deps     - Download and tidy dependencies"
	@echo "  help     - Show this help message"
