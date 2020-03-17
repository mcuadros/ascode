package types

import (
	"fmt"
	"strings"

	"github.com/mcuadros/ascode/terraform"

	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/plugin/discovery"
	"github.com/hashicorp/terraform/providers"
	"go.starlark.net/starlark"
)

func BuiltinProvider(pm *terraform.PluginManager) starlark.Value {
	return starlark.NewBuiltin("provider", func(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var name, version, alias starlark.String
		switch len(args) {
		case 3:
			var ok bool
			alias, ok = args.Index(2).(starlark.String)
			if !ok {
				return nil, fmt.Errorf("expected string, got %s", args.Index(2).Type())
			}
			fallthrough
		case 2:
			var ok bool
			version, ok = args.Index(1).(starlark.String)
			if !ok {
				return nil, fmt.Errorf("expected string, got %s", args.Index(1).Type())
			}
			fallthrough
		case 1:
			var ok bool
			name, ok = args.Index(0).(starlark.String)
			if !ok {
				return nil, fmt.Errorf("expected string, got %s", args.Index(0).Type())
			}
		default:
			return nil, fmt.Errorf("unexpected positional arguments count")
		}

		p, err := MakeProvider(pm, name.GoString(), version.GoString(), alias.GoString())
		if err != nil {
			return nil, err
		}

		return p, p.loadKeywordArgs(kwargs)
	})
}

// Provider represents a provider as a starlark.Value.
type Provider struct {
	name, alias string
	provider    *plugin.GRPCProvider
	meta        discovery.PluginMeta

	dataSources *MapSchema
	resources   *MapSchema

	*Resource
}

// MakeProvider returns a new Provider instance from a given name version and alias.
func MakeProvider(pm *terraform.PluginManager, name, version, alias string) (*Provider, error) {
	cli, meta, err := pm.Provider(name, version, false)
	if err != nil {
		return nil, err
	}

	rpc, err := cli.Client()
	if err != nil {
		return nil, err
	}

	raw, err := rpc.Dispense(plugin.ProviderPluginName)
	if err != nil {
		return nil, err
	}

	if alias == "" {
		alias = NameGenerator()
	}

	provider := raw.(*plugin.GRPCProvider)
	response := provider.GetSchema()

	defer cli.Kill()
	p := &Provider{
		name:     name,
		alias:    alias,
		provider: provider,
		meta:     meta,

		Resource: MakeResource(alias, name, ProviderKind, response.Provider.Block, nil),
	}

	p.dataSources = NewMapSchema(p, name, DataSourceKind, response.DataSources)
	p.resources = NewMapSchema(p, name, ResourceKind, response.ResourceTypes)

	return p, nil
}

func (p *Provider) String() string {
	return fmt.Sprintf("provider(%q)", p.name)
}

// Type honors the starlark.Value interface. It shadows p.Resource.Type.
func (p *Provider) Type() string {
	return fmt.Sprintf("Provider<%s>", p.name)
}

// Attr honors the starlark.HasAttrs interface.
func (p *Provider) Attr(name string) (starlark.Value, error) {
	switch name {
	case "alias":
		return starlark.String(p.alias), nil
	case "version":
		return starlark.String(p.meta.Version), nil
	case "data":
		return p.dataSources, nil
	case "resource":
		return p.resources, nil
	}

	return p.Resource.Attr(name)
}

// AttrNames honors the starlark.HasAttrs interface.
func (p *Provider) AttrNames() []string {
	return append(p.Resource.AttrNames(), "data", "resource", "version")
}

type MapSchema struct {
	p *Provider

	prefix      string
	kind        Kind
	schemas     map[string]providers.Schema
	collections map[string]*ResourceCollection
}

func NewMapSchema(p *Provider, prefix string, k Kind, schemas map[string]providers.Schema) *MapSchema {
	return &MapSchema{
		p:           p,
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
	name = m.prefix + "_" + name

	if c, ok := m.collections[name]; ok {
		return c, nil
	}

	if schema, ok := m.schemas[name]; ok {
		m.collections[name] = NewResourceCollection(name, m.kind, schema.Block, m.p.Resource)
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
