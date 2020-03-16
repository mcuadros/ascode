package types

import (
	"fmt"

	backend "github.com/hashicorp/terraform/backend/init"
	"go.starlark.net/starlark"
)

func init() {
	backend.Init(nil)
}

func BuiltinBackend() starlark.Value {
	return starlark.NewBuiltin("backend", func(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
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

		p, err := MakeBackend(name.GoString())
		if err != nil {
			return nil, err
		}

		return p, p.loadKeywordArgs(kwargs)
	})
}

type Backend struct {
	*Resource
}

func MakeBackend(name string) (*Backend, error) {
	fn := backend.Backend(name)
	if fn == nil {
		return nil, fmt.Errorf("unable to find backend %q", name)
	}

	return &Backend{
		Resource: MakeResource(name, "", BackendKind, fn().ConfigSchema(), nil),
	}, nil
}
