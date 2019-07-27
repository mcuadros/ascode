package types

import (
	"fmt"
	"strings"

	"github.com/zclconf/go-cty/cty"
	"go.starlark.net/starlark"
)

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
	case "dict", "Resource":
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

	if typ.IsMapType() {
		t.typ = "dict"
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
	case *starlark.Dict:
		if t.cty.IsMapType() {
			return nil
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
