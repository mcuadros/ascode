package types

import (
	"fmt"

	"github.com/hashicorp/terraform/configs/configschema"
	"go.starlark.net/starlark"
)

type ResourceCollection struct {
	typ    string
	kind   Kind
	block  *configschema.Block
	parent *Resource
	*starlark.List
}

func NewResourceCollection(typ string, k Kind, block *configschema.Block, parent *Resource) *ResourceCollection {
	return &ResourceCollection{
		typ:    typ,
		kind:   k,
		block:  block,
		parent: parent,
		List:   starlark.NewList(nil),
	}
}

// String honors the starlark.Value interface.
func (c *ResourceCollection) String() string {
	return fmt.Sprintf("%s", c.typ)
}

// Type honors the starlark.Value interface.
func (c *ResourceCollection) Type() string {
	return fmt.Sprintf("ResourceCollection<%s.%s>", c.kind, c.typ)
}

// Truth honors the starlark.Value interface.
func (c *ResourceCollection) Truth() starlark.Bool {
	return true // even when empty
}

// Freeze honors the starlark.Value interface.
func (c *ResourceCollection) Freeze() {}

// Hash honors the starlark.Value interface.
func (c *ResourceCollection) Hash() (uint32, error) { return 42, nil }

// Name honors the starlark.Callable interface.
func (c *ResourceCollection) Name() string {
	return c.typ
}

// CallInternal honos the starlark.Callable interface.
func (c *ResourceCollection) CallInternal(thread *starlark.Thread, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	name, dict, err := c.unpackArgs(args, kwargs)
	if err != nil {
		return nil, err
	}

	if (c.kind == ResourceKind || c.kind == DataSourceKind) && name == "" {
		name = NameGenerator()
	}

	resource := MakeResource(name, c.typ, c.kind, c.block, c.parent)
	if dict != nil && dict.Len() != 0 {
		if err := resource.loadDict(dict); err != nil {
			return nil, err
		}
	}

	if err := c.List.Append(resource); err != nil {
		return nil, err
	}

	return resource, nil
}

func (c *ResourceCollection) unpackArgsWithKwargs(args starlark.Tuple, kwargs []starlark.Tuple) (string, *starlark.Dict, error) {
	dict := starlark.NewDict(len(kwargs))
	var name starlark.String

	for _, kwarg := range kwargs {
		dict.SetKey(kwarg.Index(0), kwarg.Index(1))
	}

	if len(args) == 1 {
		var ok bool
		name, ok = args.Index(0).(starlark.String)
		if !ok {
			return "", nil, fmt.Errorf("resource: expected string, got %s", args.Index(0).Type())
		}
	}

	if len(args) > 1 {
		return "", nil, fmt.Errorf("resource: unexpected positional args mixed with kwargs")
	}

	return string(name), dict, nil
}

func (c *ResourceCollection) unpackArgs(args starlark.Tuple, kwargs []starlark.Tuple) (string, *starlark.Dict, error) {
	var dict *starlark.Dict
	var name starlark.String

	if len(args) == 0 && len(kwargs) == 0 {
		return "", nil, nil
	}

	if len(kwargs) != 0 {
		return c.unpackArgsWithKwargs(args, kwargs)
	}

	switch len(args) {
	case 0:
	case 1:
		switch v := args.Index(0).(type) {
		case starlark.String:
			return string(v), nil, nil
		case *starlark.Dict:
			return "", v, nil
		default:
			return "", nil, fmt.Errorf("resource: expected string or dict, got %s", args.Index(0).Type())
		}
	case 2:
		var ok bool
		name, ok = args.Index(0).(starlark.String)
		if !ok {
			return "", nil, fmt.Errorf("resource: expected string, got %s", args.Index(0).Type())
		}

		dict, ok = args.Index(1).(*starlark.Dict)
		if !ok {
			return "", nil, fmt.Errorf("resource: expected dict, got %s", args.Index(1).Type())
		}
	default:
		if c.kind != NestedKind {
			return "", nil, fmt.Errorf("resource: unexpected positional arguments count")
		}
	}

	return string(name), dict, nil
}

func (c *ResourceCollection) toDict() *starlark.List {
	values := make([]starlark.Value, c.Len())
	for i := 0; i < c.Len(); i++ {
		values[i] = c.Index(i).(*Resource).toDict()
	}

	return starlark.NewList(values)
}
