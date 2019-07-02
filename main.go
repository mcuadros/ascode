package main

import (
	"io/ioutil"
	"log"
	"os"

	"go.starlark.net/repl"
	"go.starlark.net/resolve"
	"go.starlark.net/starlark"
)

func main() {
	log.SetOutput(ioutil.Discard)

	pm := &PluginManager{".providers"}
	resolve.AllowFloat = true

	provider := starlark.NewBuiltin("provider", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		name := args.Index(0).(starlark.String)
		version := args.Index(1).(starlark.String)

		return NewProviderInstance(pm, string(name), string(version))
	})

	thread := &starlark.Thread{Name: "thread", Load: repl.MakeLoad()}
	predeclared := starlark.StringDict{
		"provider": provider,
	}

	if _, err := starlark.ExecFile(thread, os.Args[1], nil, predeclared); err != nil {
		panic(err)
	}
}
