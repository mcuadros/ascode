# Documentation
OUTLINE_CMD ?= outline
DOCUMENTATION_PATH ?= _documentation
DOCUMENTATION_RUNTIME_PATH ?= $(DOCUMENTATION_PATH)/reference
EXAMPLES_PATH ?= starlark/types/testdata/examples

RUNTIME_MODULES = \
	github.com/mcuadros/ascode/starlark/module/os \
	github.com/mcuadros/ascode/starlark/types \
	github.com/mcuadros/ascode/starlark/module/filepath \
	github.com/qri-io/starlib/encoding/base64 \
	github.com/qri-io/starlib/encoding/csv \
	github.com/qri-io/starlib/encoding/json \
	github.com/qri-io/starlib/encoding/yaml \
	github.com/qri-io/starlib/re \
	github.com/qri-io/starlib/http

# Build Info 
GO_LDFLAGS_CMD = go run _scripts/goldflags.go
GO_LDFLAGS_PACKAGE = cmd
GO_LDFLAGS_PACKAGES = \
 	starlarkVersion=go.starlark.net \
	terraformVersion=github.com/hashicorp/terraform


# Rules
.PHONY: $(RUNTIME_MODULES) $(COMMANDS) documentation

documentation: $(RUNTIME_MODULES)
$(RUNTIME_MODULES): $(DOCUMENTATION_RUNTIME_PATH)
	$(OUTLINE_CMD) package -t _scripts/template.md -d $(EXAMPLES_PATH) $@ \
		> $(DOCUMENTATION_RUNTIME_PATH)/`basename $@`.md

$(DOCUMENTATION_RUNTIME_PATH):
	mkdir -p $@

goldflags:
	@$(GO_LDFLAGS_CMD) $(GO_LDFLAGS_PACKAGE) . $(GO_LDFLAGS_PACKAGES)