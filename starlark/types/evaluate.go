package types

import (
	"fmt"
	"path/filepath"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

func BuiltinEvaluate() starlark.Value {
	return starlark.NewBuiltin("evaluate", func(t *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var raw starlark.String
		switch len(args) {
		case 1:
			var ok bool
			raw, ok = args.Index(0).(starlark.String)
			if !ok {
				return nil, fmt.Errorf("expected string, got %s", args.Index(0).Type())
			}
		default:
			return nil, fmt.Errorf("unexpected positional arguments count")
		}

		dict := starlark.StringDict{}
		for _, kwarg := range kwargs {
			dict[kwarg.Index(0).(starlark.String).GoString()] = kwarg.Index(1)
		}

		filename := raw.GoString()
		if base, ok := t.Local("base_path").(string); ok {
			filename = filepath.Join(base, filename)
		}

		_, file := filepath.Split(filename)
		name := file[:len(file)-len(filepath.Ext(file))]

		global, err := starlark.ExecFile(t, filename, nil, dict)
		return &starlarkstruct.Module{
			Name:    name,
			Members: global,
		}, err
	})
}
