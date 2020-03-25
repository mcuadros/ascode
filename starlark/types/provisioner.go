package types

import (
	"fmt"

	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/plugin/discovery"
	"github.com/mcuadros/ascode/terraform"
	"go.starlark.net/starlark"
)

// BuiltinProvisioner returns a starlak.Builtin function capable of instantiate
// new Provisioner instances.
//
//   outline: types
//     functions:
//       provisioner(type) Provisioner
//         Instantiates a new Provisioner
//
//         params:
//           type string
//             Provisioner type.
//
func BuiltinProvisioner() starlark.Value {
	return starlark.NewBuiltin("provisioner", MakeProvisioner)
}

// MakeProvisioner defines the Provisioner constructor.
func MakeProvisioner(
	t *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple,
) (starlark.Value, error) {

	pm := t.Local(PluginManagerLocal).(*terraform.PluginManager)

	var name starlark.String
	switch len(args) {
	case 1:
		var ok bool
		name, ok = args.Index(0).(starlark.String)
		if !ok {
			return nil, fmt.Errorf("expected string, got %s", args.Index(0).Type())
		}
	default:
		return nil, fmt.Errorf("unexpected positional arguments count")
	}

	p, err := NewProvisioner(pm, name.GoString())
	if err != nil {
		return nil, err
	}

	return p, p.loadKeywordArgs(kwargs)
}

// Provisioner represents a Terraform provider of a specif type.
//
//   outline: types
//     types:
//       Provisioner
//         Provisioner represents a Terraform provider of a specif type. As
//         written in the terraform documentation: "*Provisioners are a Last Resort*"
//
//         fields:
//           __kind__ string
//             Kind of the provisioner. Fixed value `provisioner`
//           __type__ string
//             Type of the resource. Eg.: `aws_instance
//           <argument> <scalar>
//             Arguments defined by the provisioner schema, thus can be of any
//             scalar type.
//           <block> Resource
//             Blocks defined by the provisioner schema, thus are nested resources,
//             containing other arguments and/or blocks.
//
type Provisioner struct {
	provisioner *plugin.GRPCProvisioner
	meta        discovery.PluginMeta
	*Resource
}

// NewProvisioner returns a new Provisioner for the given type.
func NewProvisioner(pm *terraform.PluginManager, typ string) (*Provisioner, error) {
	cli, meta, err := pm.Provisioner(typ)
	if err != nil {
		return nil, err
	}

	rpc, err := cli.Client()
	if err != nil {
		return nil, err
	}

	raw, err := rpc.Dispense(plugin.ProvisionerPluginName)
	if err != nil {
		return nil, err
	}

	provisioner := raw.(*plugin.GRPCProvisioner)
	response := provisioner.GetSchema()

	defer cli.Kill()
	return &Provisioner{
		provisioner: provisioner,
		meta:        meta,

		Resource: NewResource(NameGenerator(), typ, ProvisionerKind, response.Provisioner, nil, nil),
	}, nil
}

// Type honors the starlark.Value interface. It shadows p.Resource.Type.
func (p *Provisioner) Type() string {
	return "Provisioner"
}

func (p *Provisioner) String() string {
	return fmt.Sprintf("Provisioner<%s>", p.typ)
}
