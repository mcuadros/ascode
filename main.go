package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/mcuadros/ascode/cmd"
)

var version string
var build string

func main() {
	parser := flags.NewNamedParser("ascode", flags.Default)
	parser.LongDescription = "AsCode - The real infrastructure as code."
	parser.AddCommand("run", cmd.RunCmdShortDescription, cmd.RunCmdLongDescription, &cmd.RunCmd{})
	parser.AddCommand("repl", cmd.REPLCmdShortDescription, cmd.REPLCmdLongDescription, &cmd.REPLCmd{})

	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			fmt.Printf("Build information\n  commit: %s\n  date:%s\n", version, build)
			os.Exit(0)
		}

		os.Exit(1)
	}
}
