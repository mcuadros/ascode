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

	thread := &starlark.Thread{Name: "thread", Load: repl.MakeLoad()}
	predeclared := starlark.StringDict{
		"provider": types.BuiltinProvider(pm),
	}

	out, err := starlark.ExecFile(thread, os.Args[1], nil, predeclared)
	if err != nil {
		fmt.Println(err)
		if err, ok := err.(*starlark.EvalError); ok {
			fmt.Println(err.Backtrace())
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
