# Documentation
OUTLINE_CMD ?= outline
DOCUMENTATION_PATH ?= _documentation
DOCUMENTATION_REFERENCE_PATH ?= $(DOCUMENTATION_PATH)/reference
DOCUMENTATION_REFERENCE_TEMPLATE ?= $(DOCUMENTATION_REFERENCE_PATH)/reference.md.tmpl
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
.PHONY: documentation

documentation: $(RUNTIME_MODULES)
$(RUNTIME_MODULES): $(DOCUMENTATION_RUNTIME_PATH)
	$(OUTLINE_CMD) package \
		-t $(DOCUMENTATION_REFERENCE_TEMPLATE) \
		-d $(EXAMPLES_PATH) \
		$@ \
		> $(DOCUMENTATION_REFERENCE_PATH)/`basename $@`.md

$(DOCUMENTATION_REFERENCE_PATH):
	mkdir -p $@

goldflags:
	@$(GO_LDFLAGS_CMD) $(GO_LDFLAGS_PACKAGE) . $(GO_LDFLAGS_PACKAGES)