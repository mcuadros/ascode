package types

import (
	"fmt"
	"strings"

	"github.com/zclconf/go-cty/cty"
	"go.starlark.net/starlark"
)

// Value is helper to manipulate and transform starlark.Value to go types and
// cty.Value.
type Value struct {
	t Type
	v starlark.Value
}

// MustValue returns a Value from a starlark.Value, it panics if error.
func MustValue(v starlark.Value) *Value {
	value, err := NewValue(v)
	if err != nil {
		panic(err)
	}

	return value
}

// NewValue returns a Value from a starlark.Value.
func NewValue(v starlark.Value) (*Value, error) {
	t, err := NewTypeFromStarlark(v.Type())
	if err != nil {
		return nil, err
	}

	return &Value{t: *t, v: v}, nil
}

// Value returns the starlark.Value.
func (v *Value) Value() starlark.Value {
	return v.v
}

// Type returns the Type of the value.
func (v *Value) Type() *Type {
	return &v.t
}

// Cty returns the cty.Value.
func (v *Value) Cty() cty.Value {
	switch v.t.Starlark() {
	case "string":
		return cty.StringVal(v.Interface().(string))
	case "int":
		return cty.NumberIntVal(v.Interface().(int64))
	case "float":
		return cty.NumberFloatVal(v.Interface().(float64))
	case "bool":
		return cty.BoolVal(v.Interface().(bool))
	case "list":
		list := v.v.(*starlark.List)
		if list.Len() == 0 {
			return cty.ListValEmpty(v.t.Cty())
		}

		values := make([]cty.Value, list.Len())
		for i := 0; i < list.Len(); i++ {
			values[i] = MustValue(list.Index(i)).Cty()
		}

		return cty.ListVal(values)
	case "Computed":
		return cty.StringVal(v.v.(*Computed).GoString())
	default:
		return cty.StringVal(fmt.Sprintf("unhandled: %s", v.t.typ))
	}
}

// Interface returns the value as a Go value.
func (v *Value) Interface() interface{} {
	switch cast := v.v.(type) {
	case starlark.Bool:
		return bool(cast)
	case starlark.String:
		return string(cast)
	case starlark.Int:
		i, _ := cast.Int64()
		return i
	case starlark.Float:
		return float64(cast)
	case *ResourceCollection:
		return MustValue(cast.List).Interface()
	case *starlark.List:
		out := make([]interface{}, cast.Len())
		for i := 0; i < cast.Len(); i++ {
			out[i] = MustValue(cast.Index(i)).Interface()
		}

		return out
	default:
		return v
	}
}

func (v *Value) Hash() (uint32, error) {
	switch value := v.v.(type) {
	case *starlark.List:
		// Use same algorithm as Python.
		var x, mult uint32 = 0x345678, 1000003
		for i := 0; i < value.Len(); i++ {
			y, err := value.Index(i).Hash()
			if err != nil {
				return 0, err
			}
			x = x ^ y*mult
			mult += 82520 + uint32(value.Len()+value.Len())
		}
		return x, nil
	default:
		return value.Hash()
	}
}

// Type is a helper to manipulate and transform starlark.Type and cty.Type
type Type struct {
	typ string
	cty cty.Type
}

// MustTypeFromStarlark returns a Type from a given starlark type string.
// Panics if error.
func MustTypeFromStarlark(typ string) *Type {
	t, err := NewTypeFromStarlark(typ)
	if err != nil {
		panic(err)
	}

	return t
}

// NewTypeFromStarlark returns a Type from a given starlark type string.
func NewTypeFromStarlark(typ string) (*Type, error) {
	t := &Type{}
	t.typ = typ

	complex := strings.SplitN(typ, "<", 2)
	if len(complex) == 2 {
		typ = complex[0]
	}

	switch typ {
	case "bool":
		t.cty = cty.Bool
	case "int", "float":
		t.cty = cty.Number
	case "string":
		t.cty = cty.String
	case "list", "ResourceCollection":
		t.cty = cty.List(cty.NilType)
	case "Resource":
		t.cty = cty.Map(cty.NilType)
	case "Computed":
		t.cty = cty.String
	default:
		return nil, fmt.Errorf("unexpected %q type", typ)
	}

	return t, nil
}

// MustTypeFromCty returns a Type froma given cty.Type. Panics if error.
func MustTypeFromCty(typ cty.Type) *Type {
	t, err := NewTypeFromCty(typ)
	if err != nil {
		panic(err)
	}

	return t
}

// NewTypeFromCty returns a Type froma given cty.Type.
func NewTypeFromCty(typ cty.Type) (*Type, error) {
	t := &Type{}
	t.cty = typ

	switch typ {
	case cty.String:
		t.typ = "string"
	case cty.Number:
		t.typ = "int"
	case cty.Bool:
		t.typ = "bool"
	}

	if typ.IsListType() {
		t.typ = "list"
	}

	if typ.IsSetType() {
		t.typ = "set"
	}

	if typ.IsTupleType() {
		t.typ = "tuple"
	}

	return t, nil
}

// Starlark returns the type as starlark type string.
func (t *Type) Starlark() string {
	return t.typ
}

// Cty returns the type as cty.Type.
func (t *Type) Cty() cty.Type {
	return t.cty
}

// Validate validates a value againts the type.
func (t *Type) Validate(v starlark.Value) error {
	switch v.(type) {
	case starlark.String:
		if t.cty == cty.String {
			return nil
		}
	case starlark.Int, starlark.Float:
		if t.cty == cty.Number {
			return nil
		}
	case starlark.Bool:
		if t.cty == cty.Bool {
			return nil
		}
	case *Computed:
		if t.cty == v.(*Computed).t {
			return nil
		}

		vt := v.(*Computed).InnerType().Starlark()
		return fmt.Errorf("expected %s, got %s", t.typ, vt)
	case *starlark.List:
		if t.cty.IsListType() || t.cty.IsSetType() {
			return t.validateListType(v.(*starlark.List), t.cty.ElementType())
		}
	}

	return fmt.Errorf("expected %s, got %s", t.typ, v.Type())
}

func (t *Type) validateListType(l *starlark.List, expected cty.Type) error {
	for i := 0; i < l.Len(); i++ {
		if err := MustTypeFromCty(expected).Validate(l.Index(i)); err != nil {
			return fmt.Errorf("index %d: %s", i, err)
		}
	}

	return nil
}
