OUTLINE_CMD ?= outline
DOCUMENTATION_PATH ?= _documentation
DOCUMENTATION_RUNTIME_PATH ?= $(DOCUMENTATION_PATH)/runtime

RUNTIME_MODULES = \
	github.com/ascode-dev/ascode/starlark/module/os \
	github.com/ascode-dev/ascode/starlark/module/filepath \
	github.com/qri-io/starlib/encoding/base64 \
	github.com/qri-io/starlib/encoding/csv \
	github.com/qri-io/starlib/encoding/json \
	github.com/qri-io/starlib/encoding/yaml \
	github.com/qri-io/starlib/re \
	github.com/qri-io/starlib/http


# Rules
.PHONY: $(RUNTIME_MODULES) documentation

documentation: $(RUNTIME_MODULES)

$(RUNTIME_MODULES): $(DOCUMENTATION_RUNTIME_PATH)
	@$(OUTLINE_CMD) package -t _scripts/template.md $@ > $(DOCUMENTATION_RUNTIME_PATH)/`basename $@`.md

$(DOCUMENTATION_RUNTIME_PATH):
	mkdir -p $@