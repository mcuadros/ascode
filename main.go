package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/ascode-dev/ascode/starlark/runtime"
	"github.com/ascode-dev/ascode/starlark/types"
	"github.com/ascode-dev/ascode/terraform"

	"github.com/hashicorp/hcl2/hclwrite"
	"go.starlark.net/starlark"
)

func main() {
	log.SetOutput(ioutil.Discard)

	pm := &terraform.PluginManager{".providers"}
	runtime := runtime.NewRuntime(pm)

	out, err := runtime.ExecFile(os.Args[1])
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
