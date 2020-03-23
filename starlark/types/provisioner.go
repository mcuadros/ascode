package types

import (
	"fmt"

	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/plugin/discovery"
	"github.com/mcuadros/ascode/terraform"
	"go.starlark.net/starlark"
)

func BuiltinProvisioner(pm *terraform.PluginManager) starlark.Value {
	return starlark.NewBuiltin("provisioner", func(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
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

		p, err := MakeProvisioner(pm, name.GoString())
		if err != nil {
			return nil, err
		}

		return p, p.loadKeywordArgs(kwargs)
	})
}

type Provisioner struct {
	provisioner *plugin.GRPCProvisioner
	meta        discovery.PluginMeta
	*Resource
}

func MakeProvisioner(pm *terraform.PluginManager, typ string) (*Provisioner, error) {
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

		Resource: MakeResource(NameGenerator(), typ, ProvisionerKind, response.Provisioner, nil, nil),
	}, nil
}

// Type honors the starlark.Value interface. It shadows p.Resource.Type.
func (p *Provisioner) Type() string {
	return "Provisioner"
}

func (p *Provisioner) String() string {
	return fmt.Sprintf("Provisioner<%s>", p.typ)
}
