package doc

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"text/template"

	"github.com/hashicorp/terraform/plugin"
	"github.com/mcuadros/ascode/terraform"
	"github.com/stretchr/testify/assert"
)

func TestDocumentation(t *testing.T) {
	f, err := os.Open("fixtures/ignition_file.md")
	assert.NoError(t, err)

	res := NewResourceDocumentation("d", "ignition_file")
	res.Decode(f)

	assert.Len(t, res.Attribs, 7)
	assert.Len(t, res.Blocks["content"], 2)
	assert.Len(t, res.Blocks["source"], 3)
}

func TestDo(t *testing.T) {
	t.Skip()

	pm := &terraform.PluginManager{".providers"}
	cli, meta, err := pm.Provider("ignition", "1.1.0", false)
	if err != nil {
		panic(err)
	}

	fmt.Println(meta)

	rpc, err := cli.Client()
	if err != nil {
		panic(err)
	}

	raw, err := rpc.Dispense(plugin.ProviderPluginName)
	if err != nil {
		panic(err)
	}

	provider := raw.(*plugin.GRPCProvider)
	response := provider.GetSchema()

	rd, err := NewDocumentation("ignition")
	assert.NoError(t, err)

	doc := rd.Do(meta.Name, response)

	str, err := ioutil.ReadFile("../../_scripts/template.md")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	tplFuncMap := make(template.FuncMap)
	tplFuncMap["split"] = strings.Split

	temp, err := template.New("foo").Funcs(tplFuncMap).Parse(string(str))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := temp.Execute(os.Stdout, []interface{}{doc}); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
