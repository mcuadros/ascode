package types

import (
	"fmt"

	"github.com/zclconf/go-cty/cty"
	"go.starlark.net/starlark"
)

type sString = starlark.String

type Computed struct {
	r    *Resource
	t    cty.Type
	name string
	path string

	sString
}

func NewComputed(r *Resource, t cty.Type, name string) *Computed {
	var parts []string
	var path string

	child := r

	for {
		if child.parent.kind == ProviderKind {
			if child.kind == ResourceKind {
				path = fmt.Sprintf("%s.%s", child.typ, child.Name())
			} else {
				path = fmt.Sprintf("%s.%s.%s", child.kind, child.typ, child.Name())
			}

			break
		}

		parts = append(parts, child.typ)
		child = child.parent
	}

	for i := len(parts) - 1; i >= 0; i-- {
		path += "." + parts[i]
	}

	// handling of MaxItems equals 1
	block, ok := r.parent.block.BlockTypes[r.typ]
	if ok && block.MaxItems == 1 {
		name = "0." + name
	}

	return NewComputedWithPath(r, t, name, path+"."+name)
}

func NewComputedWithPath(r *Resource, t cty.Type, name, path string) *Computed {
	return &Computed{
		r:       r,
		t:       t,
		name:    name,
		path:    path,
		sString: starlark.String(fmt.Sprintf("${%s}", path)),
	}
}

func (c *Computed) Type() string {
	return fmt.Sprintf("Computed<%s>", MustTypeFromCty(c.t).Starlark())
}

func (c *Computed) InnerType() *Type {
	t, _ := NewTypeFromCty(c.t)
	return t
}

func (c *Computed) Attr(name string) (starlark.Value, error) {
	switch name {
	case "__resource__":
		return c.r, nil
	case "__type__":
		return starlark.String(MustTypeFromCty(c.t).Starlark()), nil
	}

	if !c.t.IsObjectType() {
		return nil, nil
	}

	if !c.t.HasAttribute(name) {
		return nil, nil
	}

	path := fmt.Sprintf("%s.%s", c.path, name)
	return NewComputedWithPath(c.r, c.t.AttributeType(name), name, path), nil
}

func (c *Computed) AttrNames() []string {
	return []string{"__resource__", "__type__"}
}

func (c *Computed) doNested(name, path string, t cty.Type, index int) *Computed {
	return &Computed{
		r:    c.r,
		t:    t,
		name: c.name,
	}

}

func (c *Computed) Index(i int) starlark.Value {
	path := fmt.Sprintf("%s.%d", c.path, i)

	if c.t.IsSetType() {
		return NewComputedWithPath(c.r, *c.t.SetElementType(), c.name, path)
	}

	if c.t.IsListType() {
		return NewComputedWithPath(c.r, *c.t.ListElementType(), c.name, path)
	}

	return starlark.None
}

func (c *Computed) Len() int {
	if !c.t.IsSetType() && !c.t.IsListType() {
		return 0
	}

	return 1024
}

func BuiltinFunctionComputed() starlark.Value {
	return starlark.NewBuiltin("fn", func(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var function starlark.String
		var computed *Computed
		switch len(args) {
		case 2:
			var ok bool
			function, ok = args.Index(0).(starlark.String)
			if !ok {
				return nil, fmt.Errorf("expected string, got %s", args.Index(0).Type())
			}

			computed, ok = args.Index(1).(*Computed)
			if !ok {
				return nil, fmt.Errorf("expected Computed, got %s", args.Index(1).Type())
			}
		default:
			return nil, fmt.Errorf("unexpected positional arguments count")
		}

		path := fmt.Sprintf("%s(%s)", function.GoString(), computed.path)
		return NewComputedWithPath(computed.r, computed.t, computed.name, path), nil
	})
}
