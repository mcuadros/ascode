package types

import (
	"fmt"

	"github.com/hashicorp/terraform/version"
	"github.com/mcuadros/ascode/terraform"
	"go.starlark.net/starlark"
)

type Terraform struct {
	b *Backend
	p *ProviderCollection
}

func MakeTerraform(pm *terraform.PluginManager) *Terraform {
	return &Terraform{
		p: NewProviderCollection(pm),
	}
}

// Attr honors the starlark.HasAttrs interface.
func (t *Terraform) Attr(name string) (starlark.Value, error) {
	switch name {
	case "version":
		return starlark.String(version.String()), nil
	case "provider":
		return t.p, nil
	case "backend":
		if t.b == nil {
			return starlark.None, nil
		}

		return t.b, nil
	}

	return starlark.None, nil
}

// SetField honors the starlark.HasSetField interface.
func (t *Terraform) SetField(name string, val starlark.Value) error {
	if name != "backend" {
		errmsg := fmt.Sprintf("terraform has no .%s field or method", name)
		return starlark.NoSuchAttrError(errmsg)
	}

	if b, ok := val.(*Backend); ok {
		t.b = b
		return nil
	}

	return fmt.Errorf("unexpected value %s at %s", val.Type(), name)
}

// AttrNames honors the starlark.HasAttrs interface.
func (t *Terraform) AttrNames() []string {
	return []string{"provider", "backend", "version"}
}

// Freeze honors the starlark.Value interface.
func (t *Terraform) Freeze() {} // immutable

// Hash honors the starlark.Value interface.
func (t *Terraform) Hash() (uint32, error) {
	return 0, fmt.Errorf("unhashable type: Terraform")
}

// String honors the starlark.Value interface.
func (t *Terraform) String() string {
	return "terraform"
}

// Truth honors the starlark.Value interface.
func (t *Terraform) Truth() starlark.Bool {
	return t.p.Len() != 0
}

// Type honors the starlark.Value interface.
func (t *Terraform) Type() string {
	return "Terraform"
}

var _ starlark.Value = &Terraform{}
