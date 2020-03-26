package types

import (
	"fmt"
	"strings"

	"github.com/mcuadros/ascode/terraform"

	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/plugin/discovery"
	"github.com/hashicorp/terraform/providers"
	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
)

const (
	// PluginManagerLocal is the key of the terraform.PluginManager in the thread.
	PluginManagerLocal = "plugin_manager"
)

// MakeProvider defines the Provider constructor.
func MakeProvider(
	t *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple,
) (starlark.Value, error) {

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

	pm := t.Local(PluginManagerLocal).(*terraform.PluginManager)
	p, err := NewProvider(pm, name.GoString(), version.GoString(), alias.GoString())
	if err != nil {
		return nil, err
	}

	return p, p.loadKeywordArgs(kwargs)
}

// Provider represents a provider as a starlark.Value.
//
//   outline: types
//     types:
//       Provider
//         A plugin for Terraform that makes a collection of related resources
//         available. A provider plugin is responsible for understanding API
//         interactions with some kind of service and exposing resources based
//         on that API.
//
//         examples:
//           provider.star
//           provider_resource.star
//             Resource instantiation from a Provider.
//
//         fields:
//           __version__ string
//             Provider version
//           __kind__ string
//             Kind of the provider. Fixed value `provider`
//           __type__ string
//             Type of the resource. Eg.: `aws_instance`
//           __name__ string
//             Local name of the provider, if none was provided to the constructor,
//             the name is auto-generated following the pattern `id_%s`.  At
//             Terraform is called [`alias`](https://www.terraform.io/docs/configuration/providers.html#alias-multiple-provider-instances)
//           __dict__ Dict
//             A dictionary containing all set arguments and blocks of the provider.
//           data MapSchema
//             Data sources defined by the provider.
//           resource MapSchema
//             Resources defined by the provider.
//           <argument> <scalar>
//             Arguments defined by the provider schema, thus can be of any
//             scalar type.
//           <block> Resource
//             Blocks defined by the provider schema, thus are nested resources,
//             containing other arguments and/or blocks.
//
type Provider struct {
	provider *plugin.GRPCProvider
	meta     discovery.PluginMeta

	dataSources *ResourceCollectionGroup
	resources   *ResourceCollectionGroup

	*Resource
}

var _ starlark.Value = &Provider{}
var _ starlark.HasAttrs = &Provider{}
var _ starlark.Comparable = &Provider{}

// NewProvider returns a new Provider instance from a given type, version and name.
func NewProvider(pm *terraform.PluginManager, typ, version, name string) (*Provider, error) {
	cli, meta, err := pm.Provider(typ, version, false)
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

	if name == "" {
		name = NameGenerator()
	}

	provider := raw.(*plugin.GRPCProvider)
	response := provider.GetSchema()

	defer cli.Kill()
	p := &Provider{
		provider: provider,
		meta:     meta,
	}

	p.Resource = NewResource(name, typ, ProviderKind, response.Provider.Block, p, nil)
	p.dataSources = NewResourceCollectionGroup(p, DataSourceKind, response.DataSources)
	p.resources = NewResourceCollectionGroup(p, ResourceKind, response.ResourceTypes)

	return p, nil
}

func (p *Provider) unpackArgs(args starlark.Tuple) (string, string, string, error) {
	var name, version, alias starlark.String
	switch len(args) {
	case 3:
		var ok bool
		alias, ok = args.Index(2).(starlark.String)
		if !ok {
			return "", "", "", fmt.Errorf("expected string, got %s", args.Index(2).Type())
		}
		fallthrough
	case 2:
		var ok bool
		version, ok = args.Index(1).(starlark.String)
		if !ok {
			return "", "", "", fmt.Errorf("expected string, got %s", args.Index(1).Type())
		}
		fallthrough
	case 1:
		var ok bool
		name, ok = args.Index(0).(starlark.String)
		if !ok {
			return "", "", "", fmt.Errorf("expected string, got %s", args.Index(0).Type())
		}
	default:
		return "", "", "", fmt.Errorf("unexpected positional arguments count")
	}

	return string(name), string(version), string(alias), nil
}

func (p *Provider) String() string {
	return fmt.Sprintf("Provider<%s>", p.typ)
}

// Type honors the starlark.Value interface. It shadows p.Resource.Type.
func (p *Provider) Type() string {
	return "Provider"
}

// Attr honors the starlark.Attr interface.
func (p *Provider) Attr(name string) (starlark.Value, error) {
	switch name {
	case "__version__":
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
	return append(p.Resource.AttrNames(), "data", "resource", "__version__")
}

// CompareSameType honors starlark.Comparable interface.
func (p *Provider) CompareSameType(op syntax.Token, yv starlark.Value, depth int) (bool, error) {
	y := yv.(*Provider)
	switch op {
	case syntax.EQL:
		return p == y, nil
	case syntax.NEQ:
		return p != y, nil
	default:
		return false, fmt.Errorf("%s %s %s not implemented", p.Type(), op, y.Type())
	}
}

// ResourceCollectionGroup represents a group by kind (resource or data resource)
// of ResourceCollections for a given provider.
//
//   outline: types
//     types:
//       ResourceCollectionGroup
//         ResourceCollectionGroup represents a group by kind (resource or data
//         resource) of ResourceCollections for a given provider.
//
//         fields:
//           __provider__ Provider
//             Provider of this group.
//           __kind__ string
//             Kind of the resources (`data` or `resource`).
//           <resource-name> ResourceCollection
//             Returns a ResourceCollection if the resource name is valid for
//             the schema of the provider. The resource name should be provided
//             without the provider prefix, `aws_instance` becomes
//             just an `instance`.
//
type ResourceCollectionGroup struct {
	provider    *Provider
	kind        Kind
	schemas     map[string]providers.Schema
	collections map[string]*ResourceCollection
}

var _ starlark.Value = &ResourceCollectionGroup{}
var _ starlark.HasAttrs = &ResourceCollectionGroup{}

// NewResourceCollectionGroup returns a new ResourceCollectionGroup for a given
// provider and kind based on the given schema.
func NewResourceCollectionGroup(p *Provider, k Kind, schema map[string]providers.Schema) *ResourceCollectionGroup {
	return &ResourceCollectionGroup{
		provider:    p,
		kind:        k,
		schemas:     schema,
		collections: make(map[string]*ResourceCollection),
	}
}

// Path returns the path of the ResourceCollectionGroup.
func (g *ResourceCollectionGroup) Path() string {
	return fmt.Sprintf("%s.%s", g.provider.typ, g.kind)
}

// String honors the starlark.String interface.
func (g *ResourceCollectionGroup) String() string {
	return fmt.Sprintf("ResourceCollectionGroup<%s>", g.Path())
}

// Type honors the starlark.Value interface.
func (*ResourceCollectionGroup) Type() string {
	return "ResourceCollectionGroup"
}

// Freeze honors the starlark.Value interface.
func (*ResourceCollectionGroup) Freeze() {}

// Truth honors the starlark.Value interface. True even if empty.
func (*ResourceCollectionGroup) Truth() starlark.Bool { return true }

// Hash honors the starlark.Value interface.
func (g *ResourceCollectionGroup) Hash() (uint32, error) {
	return 0, fmt.Errorf("unhashable type: %s", g.Type())
}

// Attr honors the starlark.HasAttrs interface.
func (g *ResourceCollectionGroup) Attr(name string) (starlark.Value, error) {
	switch name {
	case "__provider__":
		return g.provider, nil
	case "__kind__":
		return starlark.String(g.kind), nil
	}

	name = g.provider.typ + "_" + name
	if c, ok := g.collections[name]; ok {
		return c, nil
	}

	if schema, ok := g.schemas[name]; ok {
		g.collections[name] = NewResourceCollection(name, g.kind, schema.Block, g.provider, g.provider.Resource)
		return g.collections[name], nil
	}

	return starlark.None, nil
}

// AttrNames honors the starlark.HasAttrs interface.
func (g *ResourceCollectionGroup) AttrNames() []string {
	names := make([]string, len(g.schemas))

	var i int
	for k := range g.schemas {
		parts := strings.SplitN(k, "_", 2)
		names[i] = parts[1]
		i++
	}

	return append(names, "__kind__", "__provider__")
}
