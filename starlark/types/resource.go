package types

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/hashicorp/terraform/configs/configschema"
	"github.com/oklog/ulid"
	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
)

// NameGenerator function used to generate Resource names, by default is based
// on a ULID generator.
var NameGenerator = func() string {
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)

	return fmt.Sprintf("id_%s", ulid.MustNew(ulid.Timestamp(t), entropy))
}

// Kind describes what kind of resource is represented by a Resource isntance.
type Kind string

const (
	ProviderKind   Kind = "provider"
	ResourceKind   Kind = "resource"
	DataSourceKind Kind = "data"
	NestedKind     Kind = "nested"
)

// Resource represents a resource as a starlark.Value, it can be of four kinds,
// provider, resource, data source or a nested resource.
type Resource struct {
	name   string
	typ    string
	kind   Kind
	block  *configschema.Block
	values *Values

	parent       *Resource
	dependenies  []*Resource
	provisioners []*Provisioner
}

// MakeResource returns a new resource of the given kind, type based on the
// given configschema.Block.
func MakeResource(name, typ string, k Kind, b *configschema.Block, parent *Resource) *Resource {
	return &Resource{
		name:   name,
		typ:    typ,
		kind:   k,
		block:  b,
		values: NewValues(),
		parent: parent,
	}
}

// LoadDict loads a dict in the resource.
func (r *Resource) LoadDict(d *starlark.Dict) error {
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
	return fmt.Sprintf("Resource<%s.%s>", r.kind, r.typ)
}

// Truth honors the starlark.Value interface.
func (r *Resource) Truth() starlark.Bool {
	return true // even when empty
}

// Freeze honors the starlark.Value interface.
func (r *Resource) Freeze() {}

// Name returns the resource name based on the hash.
func (r *Resource) Name() string {
	return r.name
}

// Hash honors the starlark.Value interface.
func (r *Resource) Hash() (uint32, error) {
	return r.values.Hash()
}

// Attr honors the starlark.HasAttrs interface.
func (r *Resource) Attr(name string) (starlark.Value, error) {
	switch name {
	case "depends_on":
		return starlark.NewBuiltin("depends_on", r.dependsOn), nil
	case "add_provisioner":
		return starlark.NewBuiltin("add_provisioner", r.addProvisioner), nil
	case "__dict__":
		return r.toDict(), nil
	}

	if a, ok := r.block.Attributes[name]; ok {
		return r.attrValue(name, a)
	}

	if b, ok := r.block.BlockTypes[name]; ok {
		return r.attrBlock(name, b)
	}

	return nil, nil
}

func (r *Resource) attrBlock(name string, b *configschema.NestedBlock) (starlark.Value, error) {
	v := r.values.Get(name)
	if v != nil {
		return v.Starlark(), nil
	}

	if b.MaxItems != 1 {
		return r.values.Set(name, MustValue(NewResourceCollection(name, NestedKind, &b.Block, r))).Starlark(), nil
	}

	return r.values.Set(name, MustValue(MakeResource("", name, NestedKind, &b.Block, r))).Starlark(), nil
}

func (r *Resource) attrValue(name string, attr *configschema.Attribute) (starlark.Value, error) {
	if attr.Computed {
		if !attr.Optional || !r.values.Has(name) {
			return NewComputed(r, attr.Type, name), nil
		}
	}

	if e := r.values.Get(name); e != nil {
		return e.Starlark(), nil
	}

	return starlark.None, nil
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

	if attr.Computed && !attr.Optional {
		return fmt.Errorf("%s: can't set computed %s attribute", r.typ, name)
	}

	if err := MustTypeFromCty(attr.Type).Validate(v); err != nil {
		return err
	}

	r.values.Set(name, MustValue(v))
	return nil
}

func (r *Resource) setFieldFromNestedBlock(name string, b *configschema.NestedBlock, v starlark.Value) error {
	attr, _ := r.Attr(name)

	switch resource := attr.(type) {
	case *Resource:
		if v.Type() != "dict" {
			return fmt.Errorf("expected dict, got %s", v.Type())
		}

		return resource.LoadDict(v.(*starlark.Dict))
	case *ResourceCollection:
		if v.Type() != "list" {
			return fmt.Errorf("expected list, got %s", v.Type())
		}

		return resource.LoadList(v.(*starlark.List))
	}

	return fmt.Errorf("unexpected value %s", v.Type())
}

func (r *Resource) toDict() *starlark.Dict {
	d := starlark.NewDict(r.values.Len())

	r.values.ForEach(func(e *NamedValue) error {
		if r, ok := e.Starlark().(*Resource); ok {
			d.SetKey(starlark.String(e.Name), r.toDict())
			return nil
		}

		if r, ok := e.Starlark().(*ResourceCollection); ok {
			d.SetKey(starlark.String(e.Name), r.toDict())
			return nil
		}

		d.SetKey(starlark.String(e.Name), e.Starlark())
		return nil
	})

	return d
}

func (r *Resource) dependsOn(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, _ []starlark.Tuple) (starlark.Value, error) {
	resources := make([]*Resource, len(args))
	for i, arg := range args {
		resource, ok := arg.(*Resource)
		if !ok || resource.kind != DataSourceKind && resource.kind != ResourceKind {
			return nil, fmt.Errorf("expected Resource<[data|resource].*>, got %s", arg.Type())
		}

		if r == resource {
			return nil, fmt.Errorf("can't depend on itself")
		}

		resources[i] = resource
	}

	r.dependenies = append(r.dependenies, resources...)
	return starlark.None, nil
}

func (r *Resource) addProvisioner(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, _ []starlark.Tuple) (starlark.Value, error) {
	provisioners := make([]*Provisioner, len(args))
	for i, arg := range args {
		provisioner, ok := arg.(*Provisioner)
		if !ok {
			return nil, fmt.Errorf("expected Provisioner<*>, got %s", arg.Type())
		}

		provisioners[i] = provisioner
	}

	r.provisioners = append(r.provisioners, provisioners...)
	return starlark.None, nil
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

	if x.values.Len() != y.values.Len() {
		return false, nil
	}

	for _, xval := range x.values.List() {
		yval := y.values.Get(xval.Name)
		if yval == nil {
			return false, nil
		}

		var eq bool
		var err error
		if xcol, ok := xval.Starlark().(*ResourceCollection); ok {
			ycol, ok := yval.Starlark().(*ResourceCollection)
			if !ok {
				return false, nil
			}

			eq, err = starlark.EqualDepth(xcol.List, ycol.List, depth-1)
		} else {
			eq, err = starlark.EqualDepth(xval.Starlark(), yval.Starlark(), depth-1)
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
