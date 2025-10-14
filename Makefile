.PHONY: all build install clean

# Go parameters
BINARY_NAME=airo
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin

# Installation paths
INSTALL_PATH=/usr/local/bin

# Version information
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BUILD_DATE ?= $(shell date -u '+%Y-%m-%d %H:%M:%S UTC')

# Linker flags to inject version info
LDFLAGS=-ldflags "-X 'main.version=$(VERSION)' -X 'main.commit=$(COMMIT)' -X 'main.buildDate=$(BUILD_DATE)'"

all: build install

build:
	@echo "Building Airo..."
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"
	@echo "Build Date: $(BUILD_DATE)"
	@go build $(LDFLAGS) -o $(GOBIN)/$(BINARY_NAME) ./src/main.go

install: build
	@echo "Installing Airo to $(INSTALL_PATH)..."
	@sudo cp $(GOBIN)/$(BINARY_NAME) $(INSTALL_PATH)
	@echo "Installation complete! Run 'airo --help' to get started"

clean:
	@echo "Cleaning..."
	@rm $(INSTALL_PATH)/$(BINARY_NAME)
	@rm -rf $(GOBIN)
	@go clean
	@echo "Cleaned!"
