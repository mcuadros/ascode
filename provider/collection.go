package provider

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform/configs/configschema"
	"go.starlark.net/starlark"
)

type ResourceCollection struct {
	typ    string
	nested bool
	block  *configschema.Block
	*starlark.List
}

func NewResourceCollection(typ string, nested bool, block *configschema.Block) *ResourceCollection {
	return &ResourceCollection{
		typ:    typ,
		nested: nested,
		block:  block,
		List:   starlark.NewList(nil),
	}
}

// String honors the starlark.Value interface.
func (c *ResourceCollection) String() string {
	return fmt.Sprintf("%s", c.typ)
}

// Type honors the starlark.Value interface.
func (c *ResourceCollection) Type() string {
	return fmt.Sprintf("%s_collection", c.typ)
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

// Attr honors the starlark.HasAttrs interface.
func (c *ResourceCollection) Attr(name string) (starlark.Value, error) {
	if name == "__json__" {
		return c.toJSON(), nil
	}

	return c.List.Attr(name)
}

// CallInternal honos the starlark.Callable interface.
func (c *ResourceCollection) CallInternal(thread *starlark.Thread, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var name starlark.String
	if len(args) != 0 {
		name = args.Index(0).(starlark.String)
	}

	resource, err := MakeResource(string(name), c.typ, c.block, kwargs)
	if err != nil {
		return nil, err
	}

	if err := c.List.Append(resource); err != nil {
		return nil, err
	}

	return resource, nil
}

func (c *ResourceCollection) MarshalJSON() ([]byte, error) {
	if c.nested {
		out := ValueToNative(c.List)
		return json.Marshal(out)
	}

	out := make(map[string]interface{}, c.List.Len())
	for i := 0; i < c.List.Len(); i++ {
		r := c.List.Index(i)
		out["foo"] = r
	}

	return json.Marshal(out)
}

func (r *ResourceCollection) toJSON() starlark.String {
	json, _ := json.MarshalIndent(r, "  ", "  ")
	return starlark.String(string(json))
}

func (c *ResourceCollection) toDict() *starlark.List {
	values := make([]starlark.Value, c.Len())
	for i := 0; i < c.Len(); i++ {
		values[i] = c.Index(i).(*Resource).toDict()
	}

	return starlark.NewList(values)
}
