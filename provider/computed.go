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
	return &Computed{
		r:       r,
		a:       a,
		name:    name,
		sString: starlark.String(fmt.Sprintf("${%s.%s.%s.%s}", r.kind, r.typ, r.name, name)),
	}
}

func (*Computed) Type() string {
	return "computed"
}
