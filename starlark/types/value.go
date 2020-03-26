package types

import (
	"fmt"
	"sort"

	"github.com/hashicorp/terraform/configs/configschema"
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
	case "dict":
		dict := v.v.(*starlark.Dict)
		values := make(map[string]cty.Value)
		for _, t := range dict.Items() {
			key := fmt.Sprintf("%s", MustValue(t.Index(0)).Interface())
			values[key] = MustValue(t.Index(1)).Cty()
		}

		return cty.MapVal(values)
	case "Attribute":
		return cty.StringVal(v.v.(*Attribute).GoString())
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

	case *starlark.Dict:
		values := make(map[string]interface{})
		for _, t := range cast.Items() {
			key := fmt.Sprintf("%s", MustValue(t.Index(0)).Interface())
			values[key] = MustValue(t.Index(1)).Interface()
		}

		return values
	default:
		return v
	}
}

// Hash honors the starlark.Value interface.
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
type Values struct {
	names  sort.StringSlice
	values map[string]*NamedValue
}

// NewValues return a new instance of Values
func NewValues() *Values {
	return &Values{values: make(map[string]*NamedValue)}
}

// Set sets a name and a value and returns it as a NamedValue.
func (a *Values) Set(name string, v *Value) *NamedValue {
	if e, ok := a.values[name]; ok {
		e.Value = v
		return e
	}

	e := &NamedValue{Name: name, Value: v}
	a.values[name] = e
	a.names = append(a.names, name)
	return e
}

// Has returns true if Values contains a NamedValue with this name.
func (a Values) Has(name string) bool {
	_, ok := a.values[name]
	return ok
}

// Get returns the NamedValue with the given name, if any.
func (a Values) Get(name string) *NamedValue {
	if e, ok := a.values[name]; ok {
		return e
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
	sort.Sort(a.names) // we sort the list before hash it.
	for _, name := range a.names {
		d[name] = a.values[name].Starlark()
	}
}

// ForEach call cb for each value on Values, it stop the iteration an error
// is returned.
func (a Values) ForEach(cb func(*NamedValue) error) error {
	sort.Sort(a.names) // we sort the list before hash it.

	for _, name := range a.names {
		if err := cb(a.values[name]); err != nil {
			return err
		}
	}

	return nil
}

// List return a list of NamedValues sorted by name.
func (a Values) List() []*NamedValue {
	sort.Sort(a.names) // we sort the list before hash it.

	list := make([]*NamedValue, len(a.names))
	for i, name := range a.names {
		list[i] = a.values[name]
	}

	return list
}

// Len return the length.
func (a Values) Len() int {
	return len(a.values)
}

// Cty returns the cty.Value based on a given schema.
func (a Values) Cty(schema *configschema.Block) cty.Value {
	values := make(map[string]cty.Value)
	for key, value := range schema.Attributes {
		v := value.EmptyValue()
		if a.Has(key) {
			v = a.Get(key).Cty()
		}

		values[key] = v
	}

	return cty.ObjectVal(values)
}

// Dict is a starlark.Dict HCLCompatible.
type Dict struct {
	*starlark.Dict
}

// NewDict returns a new empty Dict.
func NewDict() *Dict {
	return &Dict{starlark.NewDict(0)}
}
