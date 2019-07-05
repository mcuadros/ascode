package provider

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
		if child.parent == nil {
			hash, _ := r.Hash()
			path = fmt.Sprintf("%s.%s.%d", child.kind, child.typ, hash)
			break
		}

		parts = append(parts, child.typ)
		child = child.parent
	}

	for i := len(parts) - 1; i >= 0; i-- {
		path += "." + parts[i]
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

func (*Computed) Type() string {
	return "computed"
}

func (c *Computed) InnerType() *Type {
	t, _ := NewTypeFromCty(c.t)
	return t
}

func (c *Computed) Attr(name string) (starlark.Value, error) {
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
	return nil
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
