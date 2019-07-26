package types

import (
	"testing"

	"go.starlark.net/starlark"

	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
)

func TestMustValue(t *testing.T) {
	testCases := []struct {
		v      starlark.Value
		cty    cty.Type
		value  cty.Value
		native interface{}
	}{
		{
			starlark.String("foo"),
			cty.String,
			cty.StringVal("foo"),
			"foo",
		},
		{
			starlark.MakeInt(42),
			cty.Number,
			cty.NumberIntVal(42),
			int64(42),
		},
		{
			starlark.Float(42),
			cty.Number,
			cty.NumberFloatVal(42),
			42.,
		},
		{
			starlark.Bool(true),
			cty.Bool,
			cty.True,
			true,
		},
		{
			starlark.NewList([]starlark.Value{starlark.String("foo")}),
			cty.List(cty.NilType),
			cty.ListVal([]cty.Value{cty.StringVal("foo")}),
			[]interface{}{"foo"},
		},
	}

	for _, tc := range testCases {
		value := MustValue(tc.v)
		assert.Equal(t, value.Type().Cty(), tc.cty)
		assert.Equal(t, value.Starlark(), tc.v)
		assert.Equal(t, value.Cty(), tc.value)
		assert.Equal(t, value.Interface(), tc.native)
	}
}

func TestValuesSet(t *testing.T) {
	var values Values
	val := values.Set("foo", MustValue(starlark.MakeInt(42)))

	assert.Equal(t, val.Name, "foo")
	assert.Equal(t, val.Interface(), int64(42))

	val = values.Set("foo", MustValue(starlark.MakeInt(84)))
	assert.Equal(t, val.Interface(), int64(84))
}

func TestValuesGet(t *testing.T) {
	var values Values
	values.Set("foo", MustValue(starlark.MakeInt(42)))
	values.Set("foo", MustValue(starlark.MakeInt(42*2)))

	val := values.Get("foo")
	assert.Equal(t, val.Interface(), int64(42*2))

	val.Value = MustValue(starlark.MakeInt(42 * 3))

	val = values.Get("foo")
	assert.Equal(t, val.Interface(), int64(42*3))

}

func TestValuesHash(t *testing.T) {
	var a Values
	a.Set("foo", MustValue(starlark.MakeInt(42)))
	a.Set("bar", MustValue(starlark.MakeInt(42*32)))

	hashA, err := a.Hash()
	assert.NoError(t, err)
	assert.Equal(t, hashA, uint32(0xfede4ab3))

	var b Values
	b.Set("bar", MustValue(starlark.MakeInt(42*32)))
	b.Set("foo", MustValue(starlark.MakeInt(42)))

	hashB, err := b.Hash()
	assert.NoError(t, err)
	assert.Equal(t, hashA, hashB)
}

func TestValuesToStringDict(t *testing.T) {
	var a Values
	a.Set("foo", MustValue(starlark.MakeInt(42)))
	a.Set("bar", MustValue(starlark.MakeInt(42*32)))

	dict := make(starlark.StringDict, 0)
	a.ToStringDict(dict)

	assert.Len(t, dict, 2)
}

func TestValuesForEach(t *testing.T) {
	var a Values
	a.Set("foo", MustValue(starlark.MakeInt(42)))
	a.Set("bar", MustValue(starlark.MakeInt(42*32)))

	var result []string
	a.ForEach(func(v *NamedValue) error {
		result = append(result, v.Name)
		return nil
	})

	assert.Equal(t, result, []string{"bar", "foo"})
}
