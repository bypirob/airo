.PHONY: all build install clean

# Go parameters
BINARY_NAME=airo
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin

# Installation paths
INSTALL_PATH=/usr/local/bin

all: build install

build:
	@echo "Building Airo..."
	@go build -o $(GOBIN)/$(BINARY_NAME) ./src/main.go

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
