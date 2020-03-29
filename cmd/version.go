package cmd

import (
	"fmt"
	"runtime"

	"github.com/jessevdk/go-flags"
)

// Command descriptions used in the flags.Parser.AddCommand.
const (
	VersionCmdShortDescription = "Version prints information about this binary."
	VersionCmdLongDescription  = VersionCmdShortDescription + "\n\n" +
		"Includes build information about AsCode like version and build\n" +
		"date but also versions from the Go runtime and other dependencies."
)

var (
	version          string
	commit           string
	build            string
	terraformVersion string
	starlarkVersion  string
	starlibVersion   string
)

// VersionCmd implements the command `version`.
type VersionCmd struct{}

// Execute honors the flags.Commander interface.
func (c *VersionCmd) Execute(args []string) error {
	fmt.Printf("Go Version: %s\n", runtime.Version())
	fmt.Printf("AsCode Version: %s\n", version)
	fmt.Printf("AsCode Commit: %s\n", commit)
	fmt.Printf("AsCode Build Date: %s\n", build)
	fmt.Printf("Terraform Version: %s\n", terraformVersion)
	fmt.Printf("Starlark Version: %s\n", starlarkVersion)
	fmt.Printf("Starlib Version: %s\n", starlibVersion)

	return nil
}

var _ flags.Commander = &VersionCmd{}
