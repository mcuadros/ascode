package types

import (
	"fmt"
	"path/filepath"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

// BuiltinEvaluate returns a starlak.Builtin function to evalute Starlark files.
//
//   outline: types
//     functions:
//       evaluate(filename, predeclared=None) dict
//         Evaluates a Starlark file and returns it's global context. Kwargs may
//         be used to set predeclared.
//         params:
//           filename string
//             Name of the file to execute.
//           predeclared? dict
//             Defines the predeclared context for the execution. Execution does
//             not modify this dictionary
//
func BuiltinEvaluate() starlark.Value {
	return starlark.NewBuiltin("evaluate", func(t *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var raw starlark.String
		predeclared := starlark.StringDict{}

		switch len(args) {
		case 2:
			dict, ok := args.Index(1).(*starlark.Dict)
			if !ok {
				return nil, fmt.Errorf("expected dict, got %s", args.Index(1).Type())
			}

			for i, key := range dict.Keys() {
				if _, ok := key.(starlark.String); ok {
					continue
				}

				return nil, fmt.Errorf("expected string keys in dict, got %s at index %d", key.Type(), i)
			}

			kwargs = append(dict.Items(), kwargs...)
			fallthrough
		case 1:
			var ok bool
			raw, ok = args.Index(0).(starlark.String)
			if !ok {
				return nil, fmt.Errorf("expected string, got %s", args.Index(0).Type())
			}
		default:
			return nil, fmt.Errorf("unexpected positional arguments count")
		}

		for _, kwarg := range kwargs {
			predeclared[kwarg.Index(0).(starlark.String).GoString()] = kwarg.Index(1)
		}

		filename := raw.GoString()
		if base, ok := t.Local("base_path").(string); ok {
			filename = filepath.Join(base, filename)
		}

		_, file := filepath.Split(filename)
		name := file[:len(file)-len(filepath.Ext(file))]

		global, err := starlark.ExecFile(t, filename, nil, predeclared)
		return &starlarkstruct.Module{
			Name:    name,
			Members: global,
		}, err
	})
}
