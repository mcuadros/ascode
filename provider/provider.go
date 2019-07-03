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

type Provider struct {
	name     string
	provider *plugin.GRPCProvider
	meta     discovery.PluginMeta

	dataSources *MapSchema
	resources   *MapSchema
}

func MakeProvider(pm *PluginManager, name, version string) (*Provider, error) {
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
	return &Provider{
		name:        name,
		provider:    provider,
		meta:        meta,
		dataSources: NewMapSchema(name, DataResourceK, response.DataSources),
		resources:   NewMapSchema(name, ResourceK, response.ResourceTypes),
	}, nil
}

func (p *Provider) String() string {
	return fmt.Sprintf("provider(%q)", p.name)
}

func (p *Provider) Type() string {
	return "provider-instance"
}

func (p *Provider) Freeze()               {}
func (p *Provider) Truth() starlark.Bool  { return true }
func (p *Provider) Hash() (uint32, error) { return 1, nil }
func (p *Provider) Name() string          { return p.name }
func (p *Provider) Attr(name string) (starlark.Value, error) {
	switch name {
	case "version":
		return starlark.String(p.meta.Version), nil
	case "data":
		return p.dataSources, nil
	case "resource":
		return p.resources, nil
	case "to_hcl":
		return BuiltinToHCL(p, hclwrite.NewEmptyFile()), nil
	}

	return starlark.None, nil
}

func (p *Provider) AttrNames() []string {
	return []string{"data", "resource"}
}

type MapSchema struct {
	prefix      string
	kind        ResourceKind
	schemas     map[string]providers.Schema
	collections map[string]*ResourceCollection
}

func NewMapSchema(prefix string, k ResourceKind, schemas map[string]providers.Schema) *MapSchema {
	return &MapSchema{
		prefix:      prefix,
		kind:        k,
		schemas:     schemas,
		collections: make(map[string]*ResourceCollection),
	}
}

func (m *MapSchema) String() string {
	return fmt.Sprintf("schemas(%q)", m.prefix)
}

func (m *MapSchema) Type() string {
	return "schemas"
}

func (m *MapSchema) Freeze()               {}
func (m *MapSchema) Truth() starlark.Bool  { return true }
func (m *MapSchema) Hash() (uint32, error) { return 1, nil }
func (m *MapSchema) Name() string          { return m.prefix }

func (m *MapSchema) Attr(name string) (starlark.Value, error) {
	if name == "to_hcl" {
		return BuiltinToHCL(m, hclwrite.NewEmptyFile()), nil
	}

	name = m.prefix + "_" + name

	if c, ok := m.collections[name]; ok {
		return c, nil
	}

	if schema, ok := m.schemas[name]; ok {
		m.collections[name] = NewResourceCollection(name, m.kind, schema.Block)
		return m.collections[name], nil
	}

	return starlark.None, nil
}

func (s *MapSchema) AttrNames() []string {
	names := make([]string, len(s.schemas))

	var i int
	for k := range s.schemas {
		parts := strings.SplitN(k, "_", 2)
		names[i] = parts[1]
		i++
	}

	return names
}
