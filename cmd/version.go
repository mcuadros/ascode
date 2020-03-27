package cmd

import (
	"fmt"
	"runtime"
)

const (
	VersionCmdShortDescription = "Version prints information about this binary."
	VersionCmdLongDescription  = VersionCmdShortDescription + "\n\n" +
		"Includes build information about AsCode like version and build\n" +
		"date but also versions from the Go runtime and other dependencies."
)

var version string
var commit string
var build string
var terraformVersion string
var starlarkVersion string

type VersionCmd struct{}

func (c *VersionCmd) Execute(args []string) error {
	fmt.Printf("Go Version: %s\n", runtime.Version())
	fmt.Printf("AsCode Version: %s\n", version)
	fmt.Printf("AsCode Commit: %s\n", commit)
	fmt.Printf("AsCode Build Date: %s\n", build)
	fmt.Printf("Terraform Version: %s\n", terraformVersion)
	fmt.Printf("Starlark Version: %s\n", starlarkVersion)

	return nil
}
