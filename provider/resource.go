package provider

import (
	"fmt"

	"github.com/hashicorp/terraform/configs/configschema"
	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
)

type fnSignature func(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error)

type ResourceKind string

const (
	ResourceK     ResourceKind = "resource"
	DataResourceK ResourceKind = "data"
	NestedK       ResourceKind = "nested"
)

type Resource struct {
	typ    string
	kind   ResourceKind
	block  *configschema.Block
	values map[string]*Value
}

func MakeResource(typ string, k ResourceKind, b *configschema.Block) *Resource {
	return &Resource{
		typ:    typ,
		kind:   k,
		block:  b,
		values: make(map[string]*Value),
	}
}

func (r *Resource) loadDict(d *starlark.Dict) error {
	for _, k := range d.Keys() {
		name := k.(starlark.String)
		value, _, _ := d.Get(k)
		if err := r.SetField(string(name), value); err != nil {
			return err
		}
	}

	return nil
}

func (r *Resource) loadKeywordArgs(kwargs []starlark.Tuple) error {
	for _, kwarg := range kwargs {
		name := kwarg.Index(0).(starlark.String)
		if err := r.SetField(string(name), kwarg.Index(1)); err != nil {
			return err
		}
	}

	return nil
}

// String honors the starlark.Value interface.
func (r *Resource) String() string {
	return fmt.Sprintf("%s(%q)", r.kind, r.typ)
}

// Type honors the starlark.Value interface.
func (r *Resource) Type() string {
	return string(r.kind)
}

// Truth honors the starlark.Value interface.
func (r *Resource) Truth() starlark.Bool {
	return true // even when empty
}

// Freeze honors the starlark.Value interface.
func (r *Resource) Freeze() {}

// Hash honors the starlark.Value interface.
func (r *Resource) Hash() (uint32, error) {
	// Same algorithm as Tuple.hash, but with different primes.
	var x, m uint32 = 8731, 9839
	for name, value := range r.values {
		namehash, _ := starlark.String(name).Hash()
		x = x ^ 3*namehash
		y, err := value.Hash()
		if err != nil {
			return 0, err
		}
		x = x ^ y*m
		m += 7349
	}

	return x, nil
}

// Attr honors the starlark.HasAttrs interface.
func (r *Resource) Attr(name string) (starlark.Value, error) {
	switch name {
	case "__dict__":
		return r.toDict(), nil
	}

	if a, ok := r.block.Attributes[name]; (ok && a.Computed) || name == "id" {
		return r.attrComputed(name, a)
	}

	if b, ok := r.block.BlockTypes[name]; ok {
		return r.attrBlock(name, b)
	}

	if v, ok := r.values[name]; ok {
		return v.Value(), nil
	}

	return nil, nil
}

func (r *Resource) attrComputed(name string, attr *configschema.Attribute) (starlark.Value, error) {
	return NewComputed(r, attr, name), nil
}

func (r *Resource) attrBlock(name string, b *configschema.NestedBlock) (starlark.Value, error) {
	if b.MaxItems != 1 {
		if _, ok := r.values[name]; !ok {
			r.values[name] = MustValue(NewResourceCollection(name, NestedK, &b.Block))
		}
	}

	if _, ok := r.values[name]; !ok {
		r.values[name] = MustValue(MakeResource(name, NestedK, &b.Block))
	}

	return r.values[name].Value(), nil
}

// AttrNames honors the starlark.HasAttrs interface.
func (r *Resource) AttrNames() []string {
	names := make([]string, len(r.block.Attributes)+len(r.block.BlockTypes))

	var i int
	for k := range r.block.Attributes {
		names[i] = k
		i++
	}

	for k := range r.block.BlockTypes {
		names[i] = k
		i++
	}

	return names
}

// SetField honors the starlark.HasSetField interface.
func (r *Resource) SetField(name string, v starlark.Value) error {
	if b, ok := r.block.BlockTypes[name]; ok {
		return r.setFieldFromNestedBlock(name, b, v)
	}

	attr, ok := r.block.Attributes[name]
	if !ok {
		errmsg := fmt.Sprintf("%s has no .%s field or method", r.typ, name)
		return starlark.NoSuchAttrError(errmsg)
	}

	if err := MustTypeFromCty(attr.Type).Validate(v); err != nil {
		return err
	}

	r.values[name] = MustValue(v)
	return nil
}

func (r *Resource) setFieldFromNestedBlock(name string, b *configschema.NestedBlock, v starlark.Value) error {
	switch v.Type() {
	case "dict":
		resource, _ := r.Attr(name)
		return resource.(*Resource).loadDict(v.(*starlark.Dict))
	}

	return fmt.Errorf("expected dict or list, got %s", v.Type())
}

func (r *Resource) localName() string {
	return fmt.Sprintf("")
}

func (r *Resource) toDict() *starlark.Dict {
	d := starlark.NewDict(len(r.values))
	for k, v := range r.values {
		if r, ok := v.Value().(*Resource); ok {
			d.SetKey(starlark.String(k), r.toDict())
			continue
		}

		if r, ok := v.Value().(*ResourceCollection); ok {
			d.SetKey(starlark.String(k), r.toDict())
			continue
		}

		d.SetKey(starlark.String(k), v.Value())
	}

	return d
}

// CompareSameType honors starlark.Comprable interface.
func (x *Resource) CompareSameType(op syntax.Token, y_ starlark.Value, depth int) (bool, error) {
	y := y_.(*Resource)
	switch op {
	case syntax.EQL:
		ok, err := x.doCompareSameType(y, depth)
		return ok, err
	case syntax.NEQ:
		ok, err := x.doCompareSameType(y, depth)
		return !ok, err
	default:
		return false, fmt.Errorf("%s %s %s not implemented", x.Type(), op, y.Type())
	}
}

func (x *Resource) doCompareSameType(y *Resource, depth int) (bool, error) {
	if x.typ != y.typ {
		return false, nil
	}

	if len(x.values) != len(y.values) {
		return false, nil
	}

	for key, xval := range x.values {
		fmt.Printf("%s = %T, %T\n", key, xval.Value(), y.values[key].Value())
		yval, found := y.values[key]
		if !found {
			return false, nil
		}

		var eq bool
		var err error
		if xcol, ok := xval.Value().(*ResourceCollection); ok {
			ycol, ok := yval.Value().(*ResourceCollection)
			if !ok {
				return false, nil
			}

			eq, err = starlark.EqualDepth(xcol.List, ycol.List, depth-1)
		} else {
			eq, err = starlark.EqualDepth(xval.Value(), yval.Value(), depth-1)
		}

		if err != nil {
			return false, err
		}

		if !eq {
			return false, nil
		}
	}

	return true, nil
}
