package types

import (
	"fmt"

	"github.com/zclconf/go-cty/cty"
	"go.starlark.net/starlark"
)

// sTring alias required to avoid name collision with the method String.
type sString = starlark.String

// Attribute is a reference to an argument of a Resource. Used mainly
// for Computed arguments of Resources.
//
//   outline: types
//     types:
//       Attribute
//         Attribute is a reference to an argument of a Resource. Used mainly
//         for Computed arguments of Resources.
//
//         Attribute behaves as the type of the argument represented, this means
//         that them can be assigned to other resource arguments of the same
//         type. And, if the type is a list are indexable.
//
//         examples:
//           attribute.star
//
//         fields:
//           __resource__ Resource
//             Resource of the attribute.
//           __type__ string
//             Type of the attribute. Eg.: `string`
type Attribute struct {
	r    *Resource
	t    cty.Type
	name string
	path string

	sString
}

var _ starlark.Value = &Attribute{}
var _ starlark.HasAttrs = &Attribute{}
var _ starlark.Indexable = &Attribute{}
var _ starlark.Comparable = &Attribute{}

// NewAttribute returns a new Attribute for a given value or block of a Resource.
// The path is calculated traversing the parents of the given Resource.
func NewAttribute(r *Resource, t cty.Type, name string) *Attribute {
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

	return NewAttributeWithPath(r, t, name, path+"."+name)
}

// NewAttributeWithPath returns a new Attribute for a given value or block of a Resource.
func NewAttributeWithPath(r *Resource, t cty.Type, name, path string) *Attribute {
	return &Attribute{
		r:       r,
		t:       t,
		name:    name,
		path:    path,
		sString: starlark.String(fmt.Sprintf("${%s}", path)),
	}
}

// Type honors the starlark.Value interface.
func (c *Attribute) Type() string {
	return fmt.Sprintf("Attribute<%s>", MustTypeFromCty(c.t).Starlark())
}

// InnerType returns the inner Type represented by this Attribute.
func (c *Attribute) InnerType() *Type {
	t, _ := NewTypeFromCty(c.t)
	return t
}

// Attr honors the starlark.HasAttrs interface.
func (c *Attribute) Attr(name string) (starlark.Value, error) {
	switch name {
	case "__resource__":
		return c.r, nil
	case "__type__":
		return starlark.String(MustTypeFromCty(c.t).Starlark()), nil
	}

	if !c.t.IsObjectType() {
		return nil, fmt.Errorf("%s it's not a object", c.Type())
	}

	if !c.t.HasAttribute(name) {
		errmsg := fmt.Sprintf("%s has no .%s field", c.Type(), name)
		return nil, starlark.NoSuchAttrError(errmsg)
	}

	path := fmt.Sprintf("%s.%s", c.path, name)
	return NewAttributeWithPath(c.r, c.t.AttributeType(name), name, path), nil
}

// AttrNames honors the starlark.HasAttrs interface.
func (c *Attribute) AttrNames() []string {
	return []string{"__resource__", "__type__"}
}

func (c *Attribute) doNested(name, path string, t cty.Type, index int) *Attribute {
	return &Attribute{
		r:    c.r,
		t:    t,
		name: c.name,
	}

}

// Index honors the starlark.Indexable interface.
func (c *Attribute) Index(i int) starlark.Value {
	path := fmt.Sprintf("%s.%d", c.path, i)

	if c.t.IsSetType() {
		return NewAttributeWithPath(c.r, *c.t.SetElementType(), c.name, path)
	}

	if c.t.IsListType() {
		return NewAttributeWithPath(c.r, *c.t.ListElementType(), c.name, path)
	}

	return starlark.None
}

// Len honors the starlark.Indexable interface.
func (c *Attribute) Len() int {
	if !c.t.IsSetType() && !c.t.IsListType() {
		return 0
	}

	return 1024
}

// BuiltinFunctionAttribute returns a built-in function that wraps Attributes
// in HCL functions.
//
//   outline: types
//     functions:
//       fn(name, target) Attribute
//         Fn wraps an Attribute in a HCL function. Since the Attributes value
//         are only available in the `apply` phase of Terraform, the only method
//         to manipulate this values is using the Terraform
//         [HCL functions](https://www.terraform.io/docs/configuration/functions.html).
//
//
//         params:
//           name string
//             Name of the HCL function to be applied. Eg.: `base64encode`
//           target Attribute
//             Target Attribute of the HCL function.
//
func BuiltinFunctionAttribute() starlark.Value {
	// TODO(mcuadros): implement multiple arguments support.
	return starlark.NewBuiltin("fn", func(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var function starlark.String
		var computed *Attribute
		switch len(args) {
		case 2:
			var ok bool
			function, ok = args.Index(0).(starlark.String)
			if !ok {
				return nil, fmt.Errorf("expected string, got %s", args.Index(0).Type())
			}

			computed, ok = args.Index(1).(*Attribute)
			if !ok {
				return nil, fmt.Errorf("expected Attribute, got %s", args.Index(1).Type())
			}
		default:
			return nil, fmt.Errorf("unexpected positional arguments count")
		}

		path := fmt.Sprintf("%s(%s)", function.GoString(), computed.path)
		return NewAttributeWithPath(computed.r, computed.t, computed.name, path), nil
	})
}
