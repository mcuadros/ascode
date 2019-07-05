package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/ascode-dev/ascode/starlark/types"
	"github.com/ascode-dev/ascode/terraform"

	"github.com/hashicorp/hcl2/hclwrite"
	"go.starlark.net/repl"
	"go.starlark.net/resolve"
	"go.starlark.net/starlark"
)

func main() {
	log.SetOutput(ioutil.Discard)

	pm := &terraform.PluginManager{".providers"}
	resolve.AllowFloat = true

	provider := starlark.NewBuiltin("provider", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		name := args.Index(0).(starlark.String)
		version := args.Index(1).(starlark.String)

		return types.MakeProvider(pm, string(name), string(version))
	})

	thread := &starlark.Thread{Name: "thread", Load: repl.MakeLoad()}
	predeclared := starlark.StringDict{
		"provider": provider,
	}

	out, err := starlark.ExecFile(thread, os.Args[1], nil, predeclared)
	if err != nil {
		if err, ok := err.(*starlark.EvalError); ok {
			log.Fatal(err.Backtrace())
		}
		log.Fatal(err)
	}

	for _, v := range out {
		p, ok := v.(*types.Provider)
		if !ok {
			continue
		}

		f := hclwrite.NewEmptyFile()
		p.ToHCL(f.Body())

		fmt.Println(string(f.Bytes()))
	}
}
