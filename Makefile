PROJECT_NAME := gossht
CMD_DIR := ./cmd

VERSION := v0.0.1-dev

GO := go
PLATFORMS := darwin freebsd linux windows
ARCHS := amd64 arm64

GOFLAGS := -ldflags "-s -w -X main.Version=$(VERSION)"

BIN_DIR := ./bin

ifeq ($(OS),Windows_NT)
    PLATFORM = windows
    ifeq ($(PROCESSOR_ARCHITECTURE),AMD64)
        MACHINE = $(PLATFORM)-amd64
    endif
    ifeq ($(PROCESSOR_ARCHITECTURE),ARM64)
        MACHINE = $(PLATFORM)-arm64
    endif
else
    UNAME_S := $(shell uname -s)
    ifeq ($(UNAME_S),Linux)
        PLATFORM = linux
    endif
    ifeq ($(UNAME_S),Darwin)
        PLATFORM = darwin
    endif
    ifeq ($(UNAME_S),FreeBSD)
        PLATFORM = freebsd
    endif
    UNAME_P := $(shell uname -p)
    ifeq ($(UNAME_P),x86_64)
        MACHINE += $(PLATFORM)-amd64
    endif
    ifeq ($(UNAME_P),arm)
        MACHINE = $(PLATFORM)-arm64
    endif
endif

# Default target
all: clean build

.PHONY: build
build: clean
	$(GO) build $(GOFLAGS) -o $(BIN_DIR)/$(PROJECT_NAME) $(CMD_DIR)

.PHONY: build-all
build-all: clean build-darwin build-freebsd build-linux build-windows

.PHONY: build-darwin
build-darwin: clean
	@echo "Building for darwin"
	GOOS=darwin GOARCH=amd64 $(GO) build $(GOFLAGS) -o $(BIN_DIR)/$(PROJECT_NAME)-darwin-amd64 $(CMD_DIR)
	GOOS=darwin GOARCH=arm64 $(GO) build $(GOFLAGS) -o $(BIN_DIR)/$(PROJECT_NAME)-darwin-arm64 $(CMD_DIR)

.PHONY: build-freebsd
build-freebsd: clean
	@echo "Building for freeBSD"
	GOOS=freebsd GOARCH=amd64 $(GO) build $(GOFLAGS) -o $(BIN_DIR)/$(PROJECT_NAME)-freebsd-amd64 $(CMD_DIR)
	GOOS=freebsd GOARCH=arm64 $(GO) build $(GOFLAGS) -o $(BIN_DIR)/$(PROJECT_NAME)-freebsd-arm64 $(CMD_DIR)

.PHONY: build-linux
build-linux: clean
	@echo "Building for linux"
	GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) -o $(BIN_DIR)/$(PROJECT_NAME)-linux-amd64 $(CMD_DIR)
	GOOS=linux GOARCH=arm64 $(GO) build $(GOFLAGS) -o $(BIN_DIR)/$(PROJECT_NAME)-linux-arm64 $(CMD_DIR)

.PHONY: build-windows
build-windows: clean
	@echo "Building for windows"
	GOOS=windows GOARCH=amd64 $(GO) build $(GOFLAGS) -o $(BIN_DIR)/$(PROJECT_NAME)-windows-amd64.exe $(CMD_DIR)
	GOOS=windows GOARCH=arm64 $(GO) build $(GOFLAGS) -o $(BIN_DIR)/$(PROJECT_NAME)-windows-arm64.exe $(CMD_DIR)

.PHONY: vendor
vendor:
	@go mod tidy
	@GOFLAGS="-mod=readonly" go mod vendor
	rm go.sum

# Clean build artifacts
clean:
	rm -rf $(BIN_DIR)
	mkdir -p $(BIN_DIR)
