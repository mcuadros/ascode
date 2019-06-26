package main

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/configs/configschema"

	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/providers"
	"go.starlark.net/starlark"
)

type ProviderInstance struct {
	name     string
	provider *plugin.GRPCProvider

	dataSources map[string]providers.Schema
	nested      map[string]*configschema.NestedBlock
}

func NewProviderInstance(pm *PluginManager, name string) (*ProviderInstance, error) {
	cli := pm.Get(name, "")

	rpc, err := cli.Client()
	if err != nil {
		return nil, err
	}

	raw, err := rpc.Dispense(plugin.ProviderPluginName)
	if err != nil {
		return nil, err
	}

	// store the client so that the plugin can kill the child process
	provider := raw.(*plugin.GRPCProvider)
	response := provider.GetSchema()

	defer cli.Kill()
	return &ProviderInstance{
		name:        name,
		provider:    provider,
		dataSources: response.DataSources,
		nested:      computeNestedBlocks(response.DataSources),
	}, nil
}

func computeNestedBlocks(s map[string]providers.Schema) map[string]*configschema.NestedBlock {
	blks := make(map[string]*configschema.NestedBlock)
	for k, block := range s {
		for n, nested := range block.Block.BlockTypes {
			key := k + "_" + n
			doComputeNestedBlocks(key, nested, blks)
		}
	}

	fmt.Println(blks)
	return blks
}

func doComputeNestedBlocks(name string, b *configschema.NestedBlock, list map[string]*configschema.NestedBlock) {
	list[name] = b
	for k, block := range b.BlockTypes {
		key := name + "_" + k
		list[key] = block

		doComputeNestedBlocks(key, block, list)
	}
}

func (t *ProviderInstance) String() string {
	return fmt.Sprintf("provider(%q)", t.name)
}

func (t *ProviderInstance) Type() string {
	return "provider-instance"
}

func (t *ProviderInstance) Freeze()               {}
func (t *ProviderInstance) Truth() starlark.Bool  { return true }
func (t *ProviderInstance) Hash() (uint32, error) { return 1, nil }
func (t *ProviderInstance) Name() string          { return t.name }
func (t *ProviderInstance) Attr(name string) (starlark.Value, error) {
	name = t.name + "_" + name

	if schema, ok := t.dataSources[name]; ok {
		return NewResourceInstanceConstructor(name, schema.Block), nil
	}

	return starlark.None, nil
}

func (t *ProviderInstance) AttrNames() []string {
	names := make([]string, len(t.dataSources)+len(t.nested))

	var i int
	for k := range t.dataSources {
		parts := strings.SplitN(k, "_", 2)
		names[i] = parts[1]
		i++
	}
	for k := range t.nested {
		parts := strings.SplitN(k, "_", 2)
		names[i] = parts[1]
		i++
	}

	return names
}
