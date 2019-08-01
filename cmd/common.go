package cmd

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/mcuadros/ascode/starlark/runtime"
	"github.com/mcuadros/ascode/terraform"
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
