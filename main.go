package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/hashicorp/hcl2/hclwrite"
	"github.com/mcuadros/terra/provider"
	"go.starlark.net/repl"
	"go.starlark.net/resolve"
	"go.starlark.net/starlark"
)

func main() {
	log.SetOutput(ioutil.Discard)

	pm := &provider.PluginManager{".providers"}
	resolve.AllowFloat = true

	pro := starlark.NewBuiltin("provider", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		name := args.Index(0).(starlark.String)
		version := args.Index(1).(starlark.String)

		return provider.MakeProvider(pm, string(name), string(version))
	})

	thread := &starlark.Thread{Name: "thread", Load: repl.MakeLoad()}
	predeclared := starlark.StringDict{
		"provider": pro,
	}

	out, err := starlark.ExecFile(thread, os.Args[1], nil, predeclared)
	if err != nil {
		if err, ok := err.(*starlark.EvalError); ok {
			log.Fatal(err.Backtrace())
		}
		log.Fatal(err)
	}

	for _, v := range out {
		p, ok := v.(*provider.Provider)
		if !ok {
			continue
		}

		f := hclwrite.NewEmptyFile()
		p.ToHCL(f.Body())

		fmt.Println(string(f.Bytes()))
	}
}
