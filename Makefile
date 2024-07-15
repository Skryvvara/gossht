PROJECT_NAME := gossht
CMD_DIR := ./cmd

GO := go
PLATFORMS := linux darwin windows freebsd
ARCHS := amd64 arm64

GOFLAGS := -ldflags "-s -w"

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

.PHONY: all clean $(foreach platform, $(PLATFORMS), $(foreach arch, $(ARCHS), build-$(platform)-$(arch)))

# Default target
all: $(foreach platform, $(PLATFORMS), $(foreach arch, $(ARCHS), build-$(platform)-$(arch)))

.PHONY: build
build: clean build-$(MACHINE)

# Build targets for each platform-architecture combination
define BUILD_template
build-$(1)-$(2):
	@echo "Building for $(1)-$(2)"
	GOOS=$(1) GOARCH=$(2) $(GO) build $(GOFLAGS) -o $(BIN_DIR)/$(PROJECT_NAME)-$(1)-$(2)$(if $(findstring windows,$(1)),.exe) $(CMD_DIR)
endef

$(foreach platform,$(PLATFORMS),$(foreach arch,$(ARCHS),$(eval $(call BUILD_template,$(platform),$(arch)))))

# Clean build artifacts
clean:
	rm -rf $(BIN_DIR)
	mkdir -p $(BIN_DIR)
