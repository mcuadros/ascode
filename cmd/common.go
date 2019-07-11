package cmd

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/ascode-dev/ascode/starlark/runtime"
	"github.com/ascode-dev/ascode/terraform"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

type commonCmd struct {
	PluginDir string `long:"plugin-dir" description:"directory containing plugin binaries" default:"$HOME/.terraform.d/plugins"`

	runtime *runtime.Runtime
}

func (c *commonCmd) init() {
	c.runtime = runtime.NewRuntime(&terraform.PluginManager{
		Path: os.ExpandEnv(c.PluginDir)},
	)
}
