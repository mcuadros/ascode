package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hashicorp/hcl/v2/hclwrite"
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
	PrintHCL       bool   `long:"print-hcl" description:"prints resources to a hcl file"`
	NoValidate     bool   `long:"no-validate" description:"skips the validation of the resources"`
	PositionalArgs struct {
		File string `positional-arg-name:"file" description:"starlark source file"`
	} `positional-args:"true" required:"1"`
}

// Execute honors the flags.Commander interface.
func (c *RunCmd) Execute(args []string) error {
	c.init()

	_, err := c.runtime.ExecFile(c.PositionalArgs.File)
	if err != nil {
		if err, ok := err.(*starlark.EvalError); ok {
			fmt.Println(err.Backtrace())
			os.Exit(1)
			return nil
		}

		return err
	}

	c.validate()
	return c.dumpToHCL()
}

func (c *RunCmd) validate() {
	if c.NoValidate {
		return
	}

	errs := c.runtime.Terraform.Validate()
	for _, err := range errs {
		fmt.Fprintln(os.Stderr, err)
	}

	if len(errs) != 0 {
		os.Exit(1)
	}
}

func (c *RunCmd) dumpToHCL() error {
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
