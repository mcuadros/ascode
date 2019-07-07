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
	return "collection"
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
	var dict *starlark.Dict

	switch len(args) {
	case 0:
	case 1:
		var ok bool
		dict, ok = args.Index(0).(*starlark.Dict)
		if !ok {
			return nil, fmt.Errorf("resource: expected dict, go %s", args.Index(0).Type())
		}
	default:
		if c.kind != NestedKind {
			return nil, fmt.Errorf("resource: unexpected positional arguments count")
		}
	}

	resource := MakeResource(c.typ, c.kind, c.block, c.parent)
	if len(kwargs) != 0 {
		if err := resource.loadKeywordArgs(kwargs); err != nil {
			return nil, err
		}
	}

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

func (c *ResourceCollection) toDict() *starlark.List {
	values := make([]starlark.Value, c.Len())
	for i := 0; i < c.Len(); i++ {
		values[i] = c.Index(i).(*Resource).toDict()
	}

	return starlark.NewList(values)
}
