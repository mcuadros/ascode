package provider

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl2/hclwrite"
	"github.com/hashicorp/terraform/plugin/discovery"

	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/providers"
	"go.starlark.net/starlark"
)

type ProviderInstance struct {
	name     string
	provider *plugin.GRPCProvider
	meta     discovery.PluginMeta

	dataSources *MapSchemaIntance
	resources   *MapSchemaIntance
}

func NewProviderInstance(pm *PluginManager, name, version string) (*ProviderInstance, error) {
	cli, meta := pm.Get(name, version)
	rpc, err := cli.Client()
	if err != nil {
		return nil, err
	}

	raw, err := rpc.Dispense(plugin.ProviderPluginName)
	if err != nil {
		return nil, err
	}

	provider := raw.(*plugin.GRPCProvider)
	response := provider.GetSchema()

	defer cli.Kill()
	return &ProviderInstance{
		name:        name,
		provider:    provider,
		meta:        meta,
		dataSources: NewMapSchemaInstance(name, response.DataSources),
		resources:   NewMapSchemaInstance(name, response.ResourceTypes),
	}, nil
}

func (t *ProviderInstance) String() string {
	return fmt.Sprintf("provider(%q)", t.name)
}

func (t *ProviderInstance) Type() string {
	return "provider-instance"
}

func (t *ProviderInstance) Freeze()               {}
func (t *ProviderInstance) Truth() starlark.Bool  { return true }
func (t *ProviderInstance) Hash() (uint32, error) { return 1, nil }
func (t *ProviderInstance) Name() string          { return t.name }
func (s *ProviderInstance) Attr(name string) (starlark.Value, error) {
	switch name {
	case "version":
		return starlark.String(s.meta.Version), nil
	case "data":
		return s.dataSources, nil
	case "resource":
		return s.resources, nil
	case "to_hcl":
		return BuiltinToHCL(s, hclwrite.NewEmptyFile()), nil
	}

	return starlark.None, nil
}

func (s *ProviderInstance) AttrNames() []string {
	return []string{"data", "resource"}
}

type MapSchemaIntance struct {
	prefix      string
	schemas     map[string]providers.Schema
	collections map[string]*ResourceCollection
}

func NewMapSchemaInstance(prefix string, schemas map[string]providers.Schema) *MapSchemaIntance {
	return &MapSchemaIntance{
		prefix:      prefix,
		schemas:     schemas,
		collections: make(map[string]*ResourceCollection),
	}
}

func (t *MapSchemaIntance) String() string {
	return fmt.Sprintf("schemas(%q)", t.prefix)
}

func (t *MapSchemaIntance) Type() string {
	return "schemas"
}

func (t *MapSchemaIntance) Freeze()               {}
func (t *MapSchemaIntance) Truth() starlark.Bool  { return true }
func (t *MapSchemaIntance) Hash() (uint32, error) { return 1, nil }
func (t *MapSchemaIntance) Name() string          { return t.prefix }

func (s *MapSchemaIntance) Attr(name string) (starlark.Value, error) {
	if name == "to_hcl" {
		return BuiltinToHCL(s, hclwrite.NewEmptyFile()), nil
	}

	name = s.prefix + "_" + name

	if c, ok := s.collections[name]; ok {
		return c, nil
	}

	if schema, ok := s.schemas[name]; ok {
		s.collections[name] = NewResourceCollection(name, false, schema.Block)
		return s.collections[name], nil
	}

	return starlark.None, nil
}

func (s *MapSchemaIntance) AttrNames() []string {
	names := make([]string, len(s.schemas))

	var i int
	for k := range s.schemas {
		parts := strings.SplitN(k, "_", 2)
		names[i] = parts[1]
		i++
	}

	return names
}
