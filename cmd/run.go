package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hashicorp/hcl2/hclwrite"
	"github.com/jessevdk/go-flags"
	"go.starlark.net/starlark"
)

// Command descriptions used in the flags.Parser.AddCommand.
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

// RunCmd implements the command `run`.
type RunCmd struct {
	commonCmd

	ToHCL          string `long:"to-hcl" description:"dumps resources to a hcl file"`
	PrintHCL       bool   `long:"print-hcl" description:"print resources to a hcl file"`
	PositionalArgs struct {
		File string `positional-arg-name:"file" description:"starlark source file"`
	} `positional-args:"true" required:"1"`
}

// Execute honors the flags.Commander interface.
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
	if c.ToHCL == "" && !c.PrintHCL {
		return nil
	}

	f := hclwrite.NewEmptyFile()
	c.runtime.Terraform.ToHCL(f.Body())

	if c.PrintHCL {
		os.Stdout.Write(f.Bytes())
	}

	if c.ToHCL == "" {
		return nil
	}

	return ioutil.WriteFile(c.ToHCL, f.Bytes(), 0644)
}

var _ flags.Commander = &RunCmd{}
