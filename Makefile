BASE_PATH := $(realpath -s $(dir $(abspath $(firstword $(MAKEFILE_LIST)))))

# Documentation
OUTLINE_CMD ?= outline
DOCUMENTATION_PATH ?= $(BASE_PATH)/_documentation
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

# Site
HUGO_SITE_PATH ?= $(BASE_PATH)/_site
HUGO_SITE_CONTENT_PATH ?= $(HUGO_SITE_PATH)/content
HUGO_SITE_TEMPLATE_PATH ?= $(HUGO_SITE_PATH)/themes/hugo-ascode-theme
HUGO_THEME_URL ?= https://github.com/mcuadros/hugo-ascode-theme
HUGO_PARAMS_VERSION ?= dev
export HUGO_PARAMS_VERSION
 
# Rules
.PHONY: documentation clean hugo-server

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

hugo-build: $(HUGO_SITE_PATH)
	hugo --minify --source $(HUGO_SITE_PATH) --config $(DOCUMENTATION_PATH)/config.toml

hugo-server: $(HUGO_SITE_PATH)
	hugo server --source $(HUGO_SITE_PATH) --config $(DOCUMENTATION_PATH)/config.toml

$(HUGO_SITE_PATH): $(HUGO_SITE_TEMPLATE_PATH)
	mkdir -p $@ \
	mkdir -p $(HUGO_SITE_CONTENT_PATH)
	mkdir -p $(HUGO_SITE_TEMPLATE_PATH)
	ln -s $(DOCUMENTATION_PATH) $(HUGO_SITE_CONTENT_PATH)/docs
	ln -s $(DOCUMENTATION_PATH)/_home.md $(HUGO_SITE_CONTENT_PATH)/_index.md

$(HUGO_SITE_TEMPLATE_PATH):
	git clone $(HUGO_THEME_URL) $(HUGO_SITE_TEMPLATE_PATH)

clean:
	rm -rf $(HUGO_SITE_PATH)