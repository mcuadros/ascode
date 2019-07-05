package terraform

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

func (m *PluginManager) Get(provider, version string) (*plugin.Client, discovery.PluginMeta) {
	meta, ok := m.getLocal(provider, version)
	if !ok {
		meta, ok = m.getRemote(provider, version)
	}

	return client(meta), meta
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

func (m *PluginManager) getRemote(provider, v string) (discovery.PluginMeta, bool) {
	installer := &discovery.ProviderInstaller{
		Dir:                   m.Path,
		PluginProtocolVersion: discovery.PluginInstallProtocolVersion,
		Ui:                    cli.NewMockUi(),
	}

	pm, _, err := installer.Get(provider, discovery.ConstraintStr(v).MustParse())
	if err != nil {
		panic(err)
	}

	return pm, true
}

func (m *PluginManager) getLocal(provider, version string) (discovery.PluginMeta, bool) {
	set := discovery.FindPlugins("provider", []string{m.Path})
	set = set.WithName(provider)
	if len(set) == 0 {
		return discovery.PluginMeta{}, false
	}

	if version != "" {
		set = set.WithVersion(discovery.VersionStr(version).MustParse())
	}

	if len(set) == 0 {
		return discovery.PluginMeta{}, false
	}

	return set.Newest(), true
}
