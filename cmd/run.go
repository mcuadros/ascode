package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hashicorp/hcl2/hclwrite"
	"github.com/mcuadros/ascode/starlark/types"
	"go.starlark.net/starlark"
)

const (
	RunCmdShortDescription = "Run parses, resolves, and executes a Starlark file."
	RunCmdLongDescription  = RunCmdShortDescription + "\n\n" +
		"When a provider is instantiated is automatically installed, at the \n" +
		"default location (~/.terraform.d/plugins), this can be overrided \n" +
		"using the flag `--plugin-dir=<PATH>`. \n\n" +
		"The Starlark file can be \"transpiled\" to a HCL file using the flag \n" +
		"`--to-hcl=<FILE>`. This file can be used directly with Terraform init \n" +
		"and plan commands.\n"
)

type RunCmd struct {
	commonCmd

	ToHCL          string `long:"to-hcl" description:"dumps context resources to a hcl file"`
	PositionalArgs struct {
		File string `positional-arg-name:"file" description:"starlark source file"`
	} `positional-args:"true" required:"1"`
}

func (c *RunCmd) Execute(args []string) error {
	c.init()

	out, err := c.runtime.ExecFile(c.PositionalArgs.File)
	if err != nil {
		if err, ok := err.(*starlark.EvalError); ok {
			fmt.Println(err.Backtrace())
			os.Exit(1)
			return nil
		}

		return err
	}

	return c.dumpToHCL(out)
}

func (c *RunCmd) dumpToHCL(ctx starlark.StringDict) error {
	if c.ToHCL == "" {
		return nil
	}

	f := hclwrite.NewEmptyFile()
	for _, v := range ctx {
		// TODO(mcuadros): replace this logic with a global object terraform
		switch o := v.(type) {
		case *types.Provider:
			o.ToHCL(f.Body())
		case *types.Backend:
			o.ToHCL(f.Body())
		}
	}

	return ioutil.WriteFile(c.ToHCL, f.Bytes(), 0644)
}
