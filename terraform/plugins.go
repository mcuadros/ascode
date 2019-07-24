package terraform

import (
	"os"
	"os/exec"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/terraform/command"
	tfplugin "github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/plugin/discovery"
	"github.com/mitchellh/cli"
)

// PluginManager is a wrapper arround the terraform tools to download and execute
// terraform plugins, like providers and provisioners.
type PluginManager struct {
	Path string
}

// Provider returns a client and the metadata for a given provider and version,
// first try to locate the provider in the local  path, if not found, it
// downloads it from terraform registry. If forceLocal just tries to find
// the binary in the local filesystem.
func (m *PluginManager) Provider(provider, version string, forceLocal bool) (*plugin.Client, discovery.PluginMeta, error) {
	meta, ok := m.getLocal("provider", provider, version)
	if !ok && !forceLocal {
		var err error
		meta, ok, err = m.getProviderRemote(provider, version)
		if err != nil {
			return nil, discovery.PluginMeta{}, err
		}

	}

	return client(meta), meta, nil
}

// Provisioner returns a client and the metadata for a given provisioner, it
// try to locate it at the local Path, if not try to execute it from the
// built-in plugins in the terraform binary.
func (m *PluginManager) Provisioner(provisioner string) (*plugin.Client, discovery.PluginMeta, error) {
	meta, ok := m.getLocal("provisioner", provisioner, "")
	if ok {
		return client(meta), meta, nil
	}

	// fallback to terraform internal provisioner.
	cmdLine, _ := command.BuildPluginCommandString("provisioner", provisioner)
	cmdArgv := strings.Split(cmdLine, command.TFSPACE)

	// override the internal to the terraform binary.
	cmdArgv[0] = "terraform"

	meta = discovery.PluginMeta{
		Name: provisioner,
		Path: strings.Join(cmdArgv, command.TFSPACE),
	}

	return client(meta), meta, nil
}

func client(m discovery.PluginMeta) *plugin.Client {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Level:  hclog.Error,
		Output: os.Stderr,
	})

	cmdArgv := strings.Split(m.Path, command.TFSPACE)

	return plugin.NewClient(&plugin.ClientConfig{
		Cmd:              exec.Command(cmdArgv[0], cmdArgv[1:]...),
		HandshakeConfig:  tfplugin.Handshake,
		VersionedPlugins: tfplugin.VersionedPlugins,
		Managed:          true,
		Logger:           logger,
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
		AutoMTLS:         true,
	})
}

const defaultVersionContraint = "> 0"

func (m *PluginManager) getProviderRemote(provider, v string) (discovery.PluginMeta, bool, error) {
	if v == "" {
		v = defaultVersionContraint
	}

	installer := &discovery.ProviderInstaller{
		Dir:                   m.Path,
		PluginProtocolVersion: discovery.PluginInstallProtocolVersion,
		Ui:                    cli.NewMockUi(),
	}

	meta, _, err := installer.Get(provider, discovery.ConstraintStr(v).MustParse())
	if err != nil {
		return discovery.PluginMeta{}, false, err
	}

	return meta, true, nil
}

func (m *PluginManager) getLocal(kind, provider, version string) (discovery.PluginMeta, bool) {
	set := discovery.FindPlugins(kind, []string{m.Path})
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
