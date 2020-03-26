package types

import (
	"fmt"

	"github.com/hashicorp/terraform/configs/configschema"
	"github.com/mcuadros/ascode/terraform"
	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
)

// ResourceCollection stores and instantiates resources for specific provider
// and resource.
//
//   outline: types
//     types:
//       ResourceCollection
//           ResourceCollection stores and instantiates resources for a specific
//           pair of provider and resource. The resources can be accessed by
//           indexing or using the built-in method of `dict`.
//
//         fields:
//           __provider__ Provider
//             Provider of this resource collection.
//           __kind__ string
//             Kind of the resource collection. Eg.: `data`
//           __type__ string
//             Type of the resource collection. Eg.: `aws_instance`
//
//         methods:
//           __call__(name="", values="") Resource
//             Returns a new resourced with the given name and values.
//
//             examples:
//               resource_collection_call.star
//                  Resource instantiation using values, dicts or kwargs.
//
//             params:
//               name string
//                 Local name of the resource, if `None` is provided it's
//                 autogenerated.
//               values dict
//                 List of arguments and nested blocks to be set in the new
//                 resource, these values can also be defined using `kwargs`.
//           search(key="id", value) list
//             Return all the Resources with the given value in the given key.
//
//             examples:
//               resource_collection_search.star
//
//             params:
//               name string
//                 Key to search. The default value is `id`.
//               value <any>
//                 Value to match in the given key.
//
type ResourceCollection struct {
	typ      string
	kind     Kind
	block    *configschema.Block
	provider *Provider
	parent   *Resource
	*starlark.List
}

var _ starlark.Value = &ResourceCollection{}
var _ starlark.HasAttrs = &ResourceCollection{}
var _ starlark.Callable = &ResourceCollection{}
var _ starlark.Comparable = &ResourceCollection{}

// NewResourceCollection returns a new ResourceCollection for the given values.
func NewResourceCollection(
	typ string, k Kind, block *configschema.Block, provider *Provider, parent *Resource,
) *ResourceCollection {
	return &ResourceCollection{
		typ:      typ,
		kind:     k,
		block:    block,
		provider: provider,
		parent:   parent,
		List:     starlark.NewList(nil),
	}
}

// LoadList loads a list of dicts on the collection. It clears the collection.
func (c *ResourceCollection) LoadList(l *starlark.List) error {
	if err := c.List.Clear(); err != nil {
		return err
	}

	for i := 0; i < l.Len(); i++ {
		dict, ok := l.Index(i).(*starlark.Dict)
		if !ok {
			return fmt.Errorf("%d: expected dict, got %s", i, l.Index(i).Type())
		}

		r := NewResource("", c.typ, c.kind, c.block, c.provider, c.parent)
		if dict != nil && dict.Len() != 0 {
			if err := r.loadDict(dict); err != nil {
				return err
			}
		}

		if err := c.List.Append(r); err != nil {
			return err
		}
	}

	return nil
}

// Path returns the path of the ResourceCollection.
func (c *ResourceCollection) Path() string {
	if c.parent != nil && c.parent.kind != ProviderKind {
		return fmt.Sprintf("%s.%s", c.parent.Path(), c.typ)
	}

	return fmt.Sprintf("%s.%s.%s", c.provider.typ, c.kind, c.typ)
}

// String honors the starlark.Value interface.
func (c *ResourceCollection) String() string {
	return fmt.Sprintf("ResourceCollection<%s>", c.Path())
}

// Type honors the starlark.Value interface.
func (c *ResourceCollection) Type() string {
	return "ResourceCollection"
}

// Truth honors the starlark.Value interface.
func (c *ResourceCollection) Truth() starlark.Bool {
	return true // even when empty
}

// Freeze honors the starlark.Value interface.
func (c *ResourceCollection) Freeze() {}

// Hash honors the starlark.Value interface.
func (c *ResourceCollection) Hash() (uint32, error) {
	return 0, fmt.Errorf("unhashable type: %s", c.Type())
}

// Name honors the starlark.Callable interface.
func (c *ResourceCollection) Name() string {
	return c.Type()
}

// CallInternal honors the starlark.Callable interface.
func (c *ResourceCollection) CallInternal(
	t *starlark.Thread, args starlark.Tuple, kwargs []starlark.Tuple,
) (starlark.Value, error) {

	r, err := MakeResource(c, t, nil, args, kwargs)
	if err != nil {
		return nil, err
	}

	return r, c.List.Append(r)
}

func (c *ResourceCollection) toDict() *starlark.List {
	values := make([]starlark.Value, c.Len())
	for i := 0; i < c.Len(); i++ {
		values[i] = c.Index(i).(*Resource).toDict()
	}

	return starlark.NewList(values)
}

// Attr honors the starlark.HasAttrs interface.
func (c *ResourceCollection) Attr(name string) (starlark.Value, error) {
	switch name {
	case "search":
		return starlark.NewBuiltin("search", c.search), nil
	case "__provider__":
		if c.kind.IsProviderRelated() {
			if c.provider == nil {
				return starlark.None, nil
			}

			return c.provider, nil
		}
	case "__kind__":
		return starlark.String(c.kind), nil
	case "__type__":
		return starlark.String(c.typ), nil
	}

	return c.List.Attr(name)
}

// AttrNames honors the starlark.HasAttrs interface.
func (c *ResourceCollection) AttrNames() []string {
	return append(c.List.AttrNames(),
		"search", "__provider__", "__kind__", "__type__")
}

func (c *ResourceCollection) search(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, _ []starlark.Tuple) (starlark.Value, error) {
	var key string
	var value starlark.Value

	switch len(args) {
	case 1:
		key = "id"
		value = args.Index(0)
	case 2:
		string, ok := args.Index(0).(starlark.String)
		if !ok {
			return nil, fmt.Errorf("resource: expected string, got %s", args.Index(0).Type())
		}

		key = string.GoString()
		value = args.Index(1)
	default:
		return nil, fmt.Errorf("search: unexpected positional arguments count")
	}

	list := starlark.NewList(nil)
	for i := 0; i < c.Len(); i++ {
		r := c.Index(i).(*Resource)
		v, ok := getValue(r, key).(starlark.Comparable)
		if !ok || v.Type() != value.Type() {
			continue
		}

		match, err := v.CompareSameType(syntax.EQL, value, 2)
		if err != nil {
			return starlark.None, err
		}

		if match {
			list.Append(r)
		}
	}

	return list, nil
}

func getValue(r *Resource, key string) starlark.Value {
	if key == "id" {
		return starlark.String(r.name)
	}

	if !r.values.Has(key) {
		return starlark.None
	}

	return r.values.Get(key).Starlark()
}

// ProviderCollection represents a nested Dict of providers, indexed by
// provider type and provider name.
//
//   outline: types
//     types:
//       ProviderCollection
//         ProviderCollection holds the providers in a nested dictionary,
//         indexed by provider type and provider name. The values can be
//         accessed by indexing or using the built-in method of `dict`.
//
//         examples:
//           provider_collection.star
//
type ProviderCollection struct {
	pm *terraform.PluginManager
	*Dict
}

var _ starlark.Value = &ProviderCollection{}
var _ starlark.HasAttrs = &ProviderCollection{}
var _ starlark.Callable = &ProviderCollection{}
var _ starlark.Comparable = &ProviderCollection{}

// NewProviderCollection returns a new ProviderCollection.
func NewProviderCollection(pm *terraform.PluginManager) *ProviderCollection {
	return &ProviderCollection{
		pm:   pm,
		Dict: NewDict(),
	}
}

// Type honors the starlark.Value interface.
func (c *ProviderCollection) Type() string {
	return "ProviderCollection"
}

// Truth honors the starlark.Value interface.
func (c *ProviderCollection) Truth() starlark.Bool {
	return true // even when empty
}

// Freeze honors the starlark.Value interface.
func (c *ProviderCollection) Freeze() {}

// Hash honors the starlark.Value interface.
func (c *ProviderCollection) Hash() (uint32, error) {
	return 0, fmt.Errorf("unhashable type: %s", c.Type())
}

// Name honors the starlark.Callable interface.
func (c *ProviderCollection) Name() string {
	return c.Type()
}

// CallInternal honors the starlark.Callable interface.
func (c *ProviderCollection) CallInternal(
	t *starlark.Thread, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {

	v, err := MakeProvider(t, nil, args, kwargs)
	if err != nil {
		return nil, err
	}

	p := v.(*Provider)
	n := starlark.String(p.typ)
	a := starlark.String(p.name)

	if _, ok, _ := c.Get(n); !ok {
		c.SetKey(n, NewDict())
	}

	providers, _, _ := c.Get(n)
	if _, ok, _ := providers.(*Dict).Get(a); ok {
		return nil, fmt.Errorf("already exists a provider %q with the alias %q", p.typ, p.name)

	}

	if err := providers.(*Dict).SetKey(a, p); err != nil {
		return nil, err
	}

	return v, nil
}
