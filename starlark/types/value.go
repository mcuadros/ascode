package types

import (
	"fmt"
	"sort"

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

// Starlark returns the starlark.Value.
func (v *Value) Starlark() starlark.Value {
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

// NamedValue represents a Value with a given name.
type NamedValue struct {
	Name string
	*Value
}

// Values is a list of NamedValues.
type Values []*NamedValue

// Set sets a name and a value and returns it as a NamedValue.
func (a *Values) Set(name string, v *Value) *NamedValue {
	e := a.Get(name)
	if e != nil {
		e.Value = v
		return e
	}

	e = &NamedValue{Name: name, Value: v}
	*a = append(*a, e)
	return e
}

// Has returns true if Values contains a NamedValue with this name.
func (a Values) Has(name string) bool {
	return a.Get(name) != nil
}

// Get returns the NamedValue with the given name, if any.
func (a Values) Get(name string) *NamedValue {
	for _, e := range a {
		if e.Name == name {
			return e
		}
	}

	return nil
}

// Hash honors the starlark.Value interface.
func (a Values) Hash() (uint32, error) {
	// Same algorithm as Tuple.hash, but with different primes.
	var x, m uint32 = 9199, 7207

	err := a.ForEach(func(v *NamedValue) error {
		namehash, _ := starlark.String(v.Name).Hash()
		x = x ^ 3*namehash
		y, err := v.Hash()
		if err != nil {
			return err
		}

		x = x ^ y*m
		m += 6203
		return nil
	})

	if err != nil {
		return 0, err
	}

	return x, nil
}

// ToStringDict adds a name/value entry to d for each field of the struct.
func (a Values) ToStringDict(d starlark.StringDict) {
	for _, e := range a {
		d[e.Name] = e.Starlark()
	}
}

// ForEach call cb for each value on Values, it stop the iteration an error
// is returned.
func (a Values) ForEach(cb func(*NamedValue) error) error {
	sort.Sort(a) // we sort the list before hash it.

	for _, v := range a {
		if err := cb(v); err != nil {
			return err
		}
	}

	return nil
}

func (a Values) Len() int           { return len(a) }
func (a Values) Less(i, j int) bool { return a[i].Name < a[j].Name }
func (a Values) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
