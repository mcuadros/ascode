package types

import (
	"fmt"
	"sort"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

// ValidationError is an error returned by Validabler.Validate.
type ValidationError struct {
	// Msg reason of the error
	Msg string
	// CallStack of the instantiation of the value being validated.
	CallStack starlark.CallStack
}

// NewValidationError returns a new ValidationError.
func NewValidationError(cs starlark.CallStack, format string, args ...interface{}) *ValidationError {
	return &ValidationError{
		Msg:       fmt.Sprintf(format, args...),
		CallStack: cs,
	}
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.CallStack.At(1).Pos, e.Msg)
}

// Value returns the error as a starlark.Value.
func (e *ValidationError) Value() starlark.Value {
	values := []starlark.Tuple{
		{starlark.String("msg"), starlark.String(e.Msg)},
		{starlark.String("pos"), starlark.String(e.CallStack.At(1).Pos.String())},
	}

	return starlarkstruct.FromKeywords(starlarkstruct.Default, values)
}

// ValidationErrors represents a list of ValidationErrors.
type ValidationErrors []*ValidationError

// Value returns the errors as a starlark.Value.
func (e ValidationErrors) Value() starlark.Value {
	values := make([]starlark.Value, len(e))
	for i, err := range e {
		values[i] = err.Value()
	}

	return starlark.NewList(values)
}

// Validabler defines if the resource is validable.
type Validabler interface {
	Validate() ValidationErrors
}

// BuiltinValidate returns a starlak.Builtin function to validate objects
// implementing the Validabler interface.
//
//   outline: types
//     functions:
//       validate(resource) list
//         Returns a list with validating errors if any. A validating error is
//         a struct with two fields: `msg` and `pos`
//         params:
//           resource <resource>
//             resource to be validated.
//
func BuiltinValidate() starlark.Value {
	return starlark.NewBuiltin("validate", func(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, _ []starlark.Tuple) (starlark.Value, error) {
		if args.Len() != 1 {
			return nil, fmt.Errorf("exactly one argument is required")
		}

		value := args.Index(0)
		v, ok := value.(Validabler)
		if !ok {
			return nil, fmt.Errorf("value type %s doesn't support validation", value.Type())
		}

		errors := v.Validate()
		return errors.Value(), nil
	})
}

// Validate honors the Validabler interface.
func (t *Terraform) Validate() (errs ValidationErrors) {
	if t.b != nil {
		errs = append(errs, t.b.Validate()...)
	}

	errs = append(errs, t.p.Validate()...)
	return
}

// Validate honors the Validabler interface.
func (d *Dict) Validate() (errs ValidationErrors) {
	for _, v := range d.Keys() {
		p, _, _ := d.Get(v)
		t, ok := p.(Validabler)
		if !ok {
			continue
		}

		errs = append(errs, t.Validate()...)
	}

	return
}

// Validate honors the Validabler interface.
func (p *Provider) Validate() (errs ValidationErrors) {
	errs = append(errs, p.Resource.Validate()...)
	errs = append(errs, p.dataSources.Validate()...)
	errs = append(errs, p.resources.Validate()...)

	return
}

// Validate honors the Validabler interface.
func (g *ResourceCollectionGroup) Validate() (errs ValidationErrors) {
	names := make(sort.StringSlice, len(g.collections))
	var i int
	for name := range g.collections {
		names[i] = name
		i++
	}

	sort.Sort(names)
	for _, name := range names {
		errs = append(errs, g.collections[name].Validate()...)
	}

	return
}

// Validate honors the Validabler interface.
func (c *ResourceCollection) Validate() (errs ValidationErrors) {
	if c.nestedblock != nil {
		l := c.Len()
		max, min := c.nestedblock.MaxItems, c.nestedblock.MinItems
		if max != 0 && l > max {
			errs = append(errs, NewValidationError(c.parent.CallStack(),
				"%s: max. length is %d, current len %d", c, max, l,
			))
		}

		if l < min {
			errs = append(errs, NewValidationError(c.parent.CallStack(),
				"%s: min. length is %d, current len %d", c, min, l,
			))
		}
	}

	for i := 0; i < c.Len(); i++ {
		errs = append(errs, c.Index(i).(*Resource).Validate()...)
	}

	return
}

// Validate honors the Validabler interface.
func (r *Resource) Validate() ValidationErrors {
	return append(
		r.doValidateAttributes(),
		r.doValidateBlocks()...,
	)
}

func (r *Resource) doValidateAttributes() (errs ValidationErrors) {
	for k, attr := range r.block.Attributes {
		if attr.Optional {
			continue
		}

		v := r.values.Get(k)
		if attr.Required {
			fails := v == nil
			if !fails {
				if l, ok := v.Starlark().(*starlark.List); ok && l.Len() == 0 {
					fails = true
				}
			}

			if fails {
				errs = append(errs, NewValidationError(r.CallStack(), "%s: attr %q is required", r, k))
			}
		}
	}

	return
}

func (r *Resource) doValidateBlocks() (errs ValidationErrors) {
	for k, block := range r.block.BlockTypes {
		v := r.values.Get(k)
		if block.MinItems > 0 && v == nil {
			errs = append(errs, NewValidationError(r.CallStack(), "%s: attr %q is required", r, k))
			continue
		}

		if v == nil {
			continue
		}

		errs = append(errs, v.Starlark().(Validabler).Validate()...)
	}

	return
}
