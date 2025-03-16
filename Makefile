# Variables
APP_NAME = fic
CMD_PATH = main.go
BUILD_FLAGS = -tags osusergo,netgo -ldflags '-extldflags "-static"'

# Default target
all: build

# Build the binary
build:
	@echo "🔨 Building File-Integrity-Checker..."
	@go build $(BUILD_FLAGS) -o $(APP_NAME) $(CMD_PATH)
	@echo "✅ Build complete!"

# Run the application
run: build
	@./$(APP_NAME) --help

# Clean generated files
clean:
	@echo "🧹 Cleaning up..."
	@rm -f $(APP_NAME)
	@echo "✅ Clean complete!"

# Format the code
fmt:
	@echo "🧹 Formatting code..."
	@go fmt ./...

# Run tests
test:
	@echo "🧪 Running tests..."
	@go test ./...

# Install dependencies
deps:
	@echo "📦 Installing dependencies..."
	@go mod tidy

# Install the binary to /usr/local/bin (for system-wide use)
install: build
	@echo "🚀 Installing $(APP_NAME)..."
	@sudo cp $(APP_NAME) /usr/local/bin/
	@echo "✅ Installed successfully!"

# Uninstall the binary
uninstall:
	@echo "❌ Uninstalling $(APP_NAME)..."
	@sudo rm -f /usr/local/bin/$(APP_NAME)
	@echo "✅ Uninstalled successfully!"

# Help message
help:
	@echo "Available commands:"
	@echo "  make build      - Build the application"
	@echo "  make run        - Build and run the application"
	@echo "  make clean      - Clean build artifacts"
	@echo "  make fmt        - Format the code"
	@echo "  make test       - Run tests"
	@echo "  make deps       - Install dependencies"
	@echo "  make install    - Install the binary in /usr/local/bin"
	@echo "  make uninstall  - Remove the installed binary"

.PHONY: all build run clean fmt test deps install uninstall help
