package provider

import (
	"fmt"

	"github.com/hashicorp/terraform/configs/configschema"
	"go.starlark.net/starlark"
)

type sString = starlark.String

type Computed struct {
	r    *Resource
	a    *configschema.Attribute
	name string
	sString
}

func NewComputed(r *Resource, a *configschema.Attribute, name string) *Computed {
	hash, _ := r.Hash()
	return &Computed{
		r:       r,
		a:       a,
		name:    name,
		sString: starlark.String(fmt.Sprintf("${%s.%s.%d.%s}", r.kind, r.typ, hash, name)),
	}
}

func (*Computed) Type() string {
	return "computed"
}
