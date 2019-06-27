package main

import (
	"fmt"

	"github.com/hashicorp/terraform/configs/configschema"
	"github.com/zclconf/go-cty/cty"
	"go.starlark.net/starlark"
)

type fnSignature func(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error)

type ResourceInstance struct {
	name   string
	typ    string
	block  *configschema.Block
	values map[string]starlark.Value
}

func NewResourceInstanceConstructor(typ string, b *configschema.Block, ptrs *PointerList) starlark.Value {
	return starlark.NewBuiltin(
		fmt.Sprintf("_%s", typ),
		func(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			name := args.Index(0).(starlark.String)

			resource, err := MakeResourceInstance(string(name), typ, b, kwargs)
			if err != nil {
				return nil, err
			}

			if ptrs != nil {
				if err := ptrs.Append(resource); err != nil {
					return nil, err
				}
			}

			return resource, nil
		},
	)
}

func MakeResourceInstance(name, typ string, b *configschema.Block, kwargs []starlark.Tuple) (*ResourceInstance, error) {
	r := &ResourceInstance{
		name:   name,
		typ:    typ,
		block:  b,
		values: make(map[string]starlark.Value),
	}

	return r, r.loadKeywordArgs(kwargs)
}

func (r *ResourceInstance) loadDict(d *starlark.Dict) error {
	for _, k := range d.Keys() {
		name := k.(starlark.String)
		value, _, _ := d.Get(k)
		if err := r.SetField(string(name), value); err != nil {
			return err
		}
	}

	return nil
}

func (r *ResourceInstance) loadKeywordArgs(kwargs []starlark.Tuple) error {
	for _, kwarg := range kwargs {
		name := kwarg.Index(0).(starlark.String)
		if err := r.SetField(string(name), kwarg.Index(1)); err != nil {
			return err
		}
	}

	return nil
}

// String honors the starlark.Value interface.
func (r *ResourceInstance) String() string {
	return fmt.Sprintf("%s(%q)", r.typ, r.name)
}

// Type honors the starlark.Value interface.
func (r *ResourceInstance) Type() string {
	return r.typ
}

// Truth honors the starlark.Value interface.
func (r *ResourceInstance) Truth() starlark.Bool {
	return true // even when empty
}

// Freeze honors the starlark.Value interface.
func (r *ResourceInstance) Freeze() {}

// Hash honors the starlark.Value interface.
func (r *ResourceInstance) Hash() (uint32, error) {
	return 42, nil
}

// Attr honors the starlark.HasAttrs interface.
func (r *ResourceInstance) Attr(name string) (starlark.Value, error) {
	if name == "__dict__" {
		return r.toDict(), nil
	}

	if b, ok := r.block.BlockTypes[name]; ok {
		if b.MaxItems != 1 {
			if _, ok := r.values[name]; !ok {
				r.values[name] = NewPointerList()
			}

			return NewResourceInstanceConstructor(name, &b.Block, r.values[name].(*PointerList)), nil
		}

		if _, ok := r.values[name]; !ok {
			r.values[name], _ = MakeResourceInstance("", name, &b.Block, nil)
		}
	}

	if v, ok := r.values[name]; ok {
		return v, nil
	}

	return nil, nil
}

func (r *ResourceInstance) getNestedBlockAttr(name string, b *configschema.NestedBlock) (starlark.Value, error) {
	if v, ok := r.values[name]; ok {
		return v, nil
	}

	var err error
	r.values[name], err = MakeResourceInstance("", name, &b.Block, nil)
	return r.values[name], err
}

// AttrNames honors the starlark.HasAttrs interface.
func (r *ResourceInstance) AttrNames() []string {
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
func (r *ResourceInstance) SetField(name string, v starlark.Value) error {
	if b, ok := r.block.BlockTypes[name]; ok {
		return r.setFieldFromNestedBlock(name, b, v)
	}

	attr, ok := r.block.Attributes[name]
	if !ok {
		errmsg := fmt.Sprintf("%s has no .%s field or method", r.typ, name)
		return starlark.NoSuchAttrError(errmsg)
	}

	if err := ValidateType(v, attr.Type); err != nil {
		return err
	}

	r.values[name] = v
	return nil
}

func (r *ResourceInstance) setFieldFromNestedBlock(name string, b *configschema.NestedBlock, v starlark.Value) error {
	switch v.Type() {
	case "dict":
		resource, _ := r.Attr(name)
		return resource.(*ResourceInstance).loadDict(v.(*starlark.Dict))
	}

	return fmt.Errorf("expected dict or list, got %s", v.Type())
}

func (r *ResourceInstance) toDict() *starlark.Dict {
	d := starlark.NewDict(len(r.values))
	for k, v := range r.values {
		if r, ok := v.(*ResourceInstance); ok {
			d.SetKey(starlark.String(k), r.toDict())
			continue
		}

		if r, ok := v.(*PointerList); ok {
			d.SetKey(starlark.String(k), r.toDict())
			continue
		}

		d.SetKey(starlark.String(k), v)
	}

	return d
}

type PointerList struct {
	*starlark.List
}

func NewPointerList() *PointerList {
	return &PointerList{List: starlark.NewList(nil)}
}

func (l *PointerList) toDict() *starlark.List {
	values := make([]starlark.Value, l.Len())
	for i := 0; i < l.Len(); i++ {
		values[i] = l.Index(i).(*ResourceInstance).toDict()
	}

	return starlark.NewList(values)
}

/*
NoneType                     # the type of None
bool                         # True or False
int                          # a signed integer of arbitrary magnitude
float                        # an IEEE 754 double-precision floating point number
string                       # a byte string
list                         # a modifiable sequence of values
tuple                        # an unmodifiable sequence of values
dict                         # a mapping from values to values
set                          # a set of values
function                     # a function implemented in Starlark
builtin_function_or_method
*/

func FromStarlarkType(typ string) cty.Type {
	switch typ {
	case "bool":
		return cty.Bool
	case "int":
	case "float":
		return cty.Number
	case "string":
		return cty.String
	}

	return cty.NilType
}

func ValidateListType(l *starlark.List, expected cty.Type) error {
	for i := 0; i < l.Len(); i++ {
		if err := ValidateType(l.Index(i), expected); err != nil {
			return fmt.Errorf("index %d: %s", i, err)
		}
	}

	return nil
}

func ValidateType(v starlark.Value, expected cty.Type) error {
	switch v.(type) {
	case starlark.String:
		if expected == cty.String {
			return nil
		}
	case starlark.Int:
		if expected == cty.Number {
			return nil
		}
	case starlark.Bool:
		if expected == cty.Bool {
			return nil
		}
	case *starlark.List:
		if expected.IsListType() || expected.IsSetType() {
			return ValidateListType(v.(*starlark.List), expected.ElementType())
		}
	}

	return fmt.Errorf("expected %s, got %s", ToStarlarkType(expected), v.Type())
}

func ToStarlarkType(t cty.Type) string {
	switch t {
	case cty.String:
		return "string"
	case cty.Number:
		return "int"
	case cty.Bool:
		return "bool"
	}

	if t.IsListType() {
		return "list"
	}

	if t.IsSetType() {
		return "set"
	}

	return "(unknown)"
}
