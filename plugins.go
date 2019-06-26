package main

import (
	"os"
	"os/exec"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	tfplugin "github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/plugin/discovery"
	"github.com/mitchellh/cli"
)

type PluginManager struct {
	Path string
}

func (m *PluginManager) Get(provider, version string) *plugin.Client {
	meta, ok := m.getLocal(provider, version)
	if !ok {
		meta, ok = m.getRemote(provider, version)
	}

	return client(meta)
}

func client(m discovery.PluginMeta) *plugin.Client {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Level:  hclog.Error,
		Output: os.Stderr,
	})

	return plugin.NewClient(&plugin.ClientConfig{
		Cmd:              exec.Command(m.Path),
		HandshakeConfig:  tfplugin.Handshake,
		VersionedPlugins: tfplugin.VersionedPlugins,
		Managed:          true,
		Logger:           logger,
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
		AutoMTLS:         true,
	})
}

func (m *PluginManager) getRemote(provider, version string) (discovery.PluginMeta, bool) {
	installer := &discovery.ProviderInstaller{
		Dir:                   m.Path,
		PluginProtocolVersion: discovery.PluginInstallProtocolVersion,
		Ui:                    cli.NewMockUi(),
	}

	pm, _, err := installer.Get(provider, discovery.Constraints{})
	if err != nil {
		panic(err)
	}

	return pm, true
}

func (m *PluginManager) getLocal(provider, version string) (discovery.PluginMeta, bool) {
	set := discovery.FindPlugins("provider", []string{m.Path})
	if len(set) == 0 {
		return discovery.PluginMeta{}, false
	}

	set = set.WithName(provider)
	if version != "" {
		//		set = set.WithVersion(version)
	}

	return set.Newest(), true
}
