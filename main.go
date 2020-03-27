package main

import (
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/mcuadros/ascode/cmd"
)

var version string
var build string

func main() {
	parser := flags.NewNamedParser("ascode", flags.Default)
	parser.LongDescription = "AsCode - Terraform Alternative Syntax."
	parser.AddCommand("run", cmd.RunCmdShortDescription, cmd.RunCmdLongDescription, &cmd.RunCmd{})
	parser.AddCommand("repl", cmd.REPLCmdShortDescription, cmd.REPLCmdLongDescription, &cmd.REPLCmd{})
	parser.AddCommand("version", cmd.VersionCmdShortDescription, cmd.VersionCmdLongDescription, &cmd.VersionCmd{})

	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		}

		os.Exit(1)
	}
}
