package types

import (
	"testing"

	"go.starlark.net/starlark"

	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
)

func TestNewTypeFromStarlark(t *testing.T) {
	testCases := []struct {
		typ string
		cty cty.Type
	}{
		{"bool", cty.Bool},
		{"int", cty.Number},
		{"float", cty.Number},
		{"string", cty.String},
	}

	for _, tc := range testCases {
		typ, err := NewTypeFromStarlark(tc.typ)
		assert.NoError(t, err)
		assert.Equal(t, typ.Cty(), tc.cty)
	}
}
func TestNewTypeFromStarlark_NonScalar(t *testing.T) {
	typ := MustTypeFromStarlark("list")
	assert.True(t, typ.Cty().IsListType())

	typ = MustTypeFromStarlark("dict")
	assert.True(t, typ.Cty().IsMapType())

	typ = MustTypeFromStarlark("ResourceCollection<bar>")
	assert.True(t, typ.Cty().IsListType())

	typ = MustTypeFromStarlark("Resource<foo>")
	assert.True(t, typ.Cty().IsMapType())
}

func TestNewTypeFromCty(t *testing.T) {
	testCases := []struct {
		typ string
		cty cty.Type
	}{
		{"string", cty.String},
		{"int", cty.Number},
		{"bool", cty.Bool},
		{"list", cty.List(cty.String)},
		{"dict", cty.Map(cty.String)},
		{"set", cty.Set(cty.String)},
		{"tuple", cty.Tuple([]cty.Type{})},
	}

	for _, tc := range testCases {
		typ, err := NewTypeFromCty(tc.cty)
		assert.NoError(t, err)
		assert.Equal(t, typ.Starlark(), tc.typ)
	}
}

func TestTypeValidate(t *testing.T) {
	testCases := []struct {
		t   string
		v   starlark.Value
		err bool
	}{
		{"string", starlark.String("foo"), false},
		{"int", starlark.String("foo"), true},
		{"int", starlark.MakeInt(42), false},
		{"int", starlark.MakeInt64(42), false},
		{"string", starlark.MakeInt(42), true},
		{"int", starlark.Float(42.), false},
	}

	for i, tc := range testCases {
		typ := MustTypeFromStarlark(tc.t)
		err := typ.Validate(tc.v)
		if tc.err {
			assert.Error(t, err, i)
		} else {
			assert.NoError(t, err, i)
		}
	}
}

func TestTypeValidate_Dict(t *testing.T) {
	typ := MustTypeFromCty(cty.Map(cty.String))
	dict := starlark.NewDict(1)
	dict.SetKey(starlark.String("foo"), starlark.MakeInt(42))

	err := typ.Validate(dict)
	assert.NoError(t, err)
}

func TestTypeValidate_List(t *testing.T) {
	typ := MustTypeFromCty(cty.List(cty.String))
	err := typ.Validate(starlark.NewList([]starlark.Value{
		starlark.String("foo"),
		starlark.String("bar"),
	}))

	assert.NoError(t, err)
}

func TestTypeValidate_ListError(t *testing.T) {
	typ := MustTypeFromCty(cty.List(cty.Number))
	err := typ.Validate(starlark.NewList([]starlark.Value{
		starlark.MakeInt(42),
		starlark.String("bar"),
	}))

	assert.Errorf(t, err, "index 1: expected int, got string")
}
