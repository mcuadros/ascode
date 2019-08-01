# Package configuration
PROJECT = ascode
COMMANDS = .
DEPENDENCIES =

# Documentation
OUTLINE_CMD ?= outline
DOCUMENTATION_PATH ?= _documentation
DOCUMENTATION_RUNTIME_PATH ?= $(DOCUMENTATION_PATH)/runtime

RUNTIME_MODULES = \
	github.com/mcuadros/ascode/starlark/module/os \
	github.com/mcuadros/ascode/starlark/module/filepath \
	github.com/qri-io/starlib/encoding/base64 \
	github.com/qri-io/starlib/encoding/csv \
	github.com/qri-io/starlib/encoding/json \
	github.com/qri-io/starlib/encoding/yaml \
	github.com/qri-io/starlib/re \
	github.com/qri-io/starlib/http

# Build information
BUILD ?= $(shell date +"%m-%d-%Y_%H_%M_%S")
COMMIT ?= $(shell git rev-parse --short HEAD)
GIT_DIRTY = $(shell test -n "`git status --porcelain`" && echo "-dirty" || true)
DEV_PREFIX := dev
VERSION ?= $(DEV_PREFIX)-$(COMMIT)$(GIT_DIRTY)

# Travis CI
ifneq ($(TRAVIS_TAG), )
	VERSION := $(TRAVIS_TAG)
endif

# Packages content
PKG_OS ?= darwin linux
PKG_ARCH = amd64

# Golang config
LD_FLAGS ?= -X main.version=$(VERSION) -X main.build=$(BUILD) -X main.commit=$(COMMIT)
GO_CMD = go
GO_GET = $(GO_CMD) get -v -t
GO_BUILD = $(GO_CMD) build -ldflags "$(LD_FLAGS)"
GO_TEST = $(GO_CMD) test -v

# Environment
BUILD_PATH := build
BIN_PATH := $(BUILD_PATH)/bin

# Rules
.PHONY: $(RUNTIME_MODULES) $(COMMANDS) documentation

documentation: $(RUNTIME_MODULES)
$(RUNTIME_MODULES): $(DOCUMENTATION_RUNTIME_PATH)
	@$(OUTLINE_CMD) package -t _scripts/template.md $@ > $(DOCUMENTATION_RUNTIME_PATH)/`basename $@`.md

$(DOCUMENTATION_RUNTIME_PATH):
	mkdir -p $@

build: $(COMMANDS)
$(COMMANDS):
	@if [ "$@" == "." ]; then \
		BIN=`basename $(CURDIR)` ; \
	else \
		BIN=`basename $@` ; \
	fi && \
	for os in $(PKG_OS); do \
		NBIN="$${BIN}" ; \
		if [ "$${os}" == windows ]; then \
			NBIN="$${NBIN}.exe"; \
		fi && \
		for arch in $(PKG_ARCH); do \
			mkdir -p $(BUILD_PATH)/$(PROJECT)_$${os}_$${arch} && \
			$(GO_BUILD_ENV) GOOS=$${os} GOARCH=$${arch} \
				$(GO_BUILD) -o "$(BUILD_PATH)/$(PROJECT)_$${os}_$${arch}/$${NBIN}" ./$@; \
		done; \
	done

packages: build
	@cd $(BUILD_PATH); \
	for os in $(PKG_OS); do \
		for arch in $(PKG_ARCH); do \
			TAR_VERSION=`echo $(VERSION) | tr "/" "-"`; \
			tar -cvzf $(PROJECT)_$${TAR_VERSION}_$${os}_$${arch}.tar.gz $(PROJECT)_$${os}_$${arch}/; \
		done; \
	done