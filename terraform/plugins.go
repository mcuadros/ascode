package terraform

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/terraform/addrs"
	"github.com/hashicorp/terraform/command"
	tfplugin "github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/plugin/discovery"
	"github.com/mitchellh/cli"
)

// PluginManager is a wrapper around the terraform tools to download and execute
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
		meta, ok, _ = m.getProviderRemoteDirectDownload(provider, version)
		if ok {
			return client(meta), meta, nil
		}

		var err error
		meta, _, err = m.getProviderRemote(provider, version)
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
	if !IsTerraformBinaryAvailable() {
		return nil, discovery.PluginMeta{}, ErrTerraformNotAvailable
	}

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

const releaseTemplateURL = "https://releases.hashicorp.com/terraform-provider-%s/%s/terraform-provider-%[1]s_%[2]s_%s_%s.zip"

func (m *PluginManager) getProviderRemoteDirectDownload(provider, v string) (discovery.PluginMeta, bool, error) {
	url := fmt.Sprintf(releaseTemplateURL, provider, v, runtime.GOOS, runtime.GOARCH)
	if err := m.downloadURL(url); err != nil {
		return discovery.PluginMeta{}, false, err
	}

	meta, ok := m.getLocal("provider", provider, v)
	return meta, ok, nil
}

func (m *PluginManager) downloadURL(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error downloading %s file: %w", url, err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid URL: %s", url)
	}

	defer resp.Body.Close()
	file, err := ioutil.TempFile("", "ascode")
	if err != nil {
		return err

	}

	if _, err := io.Copy(file, resp.Body); err != nil {
		return fmt.Errorf("error downloading %s file: %w", url, err)
	}

	file.Close()
	defer os.Remove(file.Name())

	archive, err := zip.OpenReader(file.Name())
	if err != nil {
		panic(err)
	}

	defer archive.Close()

	for _, f := range archive.File {
		file := filepath.Join(m.Path, f.Name)

		if !strings.HasPrefix(file, filepath.Clean(m.Path)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid path")
		}

		output, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		r, err := f.Open()
		if err != nil {
			return err
		}

		if _, err := io.Copy(output, r); err != nil {
			return err
		}

		output.Close()
		r.Close()
	}

	return nil
}

const defaultVersionContraint = "> 0"

func (m *PluginManager) getProviderRemote(provider, v string) (discovery.PluginMeta, bool, error) {
	if v == "" {
		v = defaultVersionContraint
	}

	m.getProviderRemoteDirectDownload(provider, v)
	installer := &discovery.ProviderInstaller{
		Dir:                   m.Path,
		PluginProtocolVersion: discovery.PluginInstallProtocolVersion,
		Ui:                    cli.NewMockUi(),
	}

	addr := addrs.NewLegacyProvider(provider)
	meta, _, err := installer.Get(addr, discovery.ConstraintStr(v).MustParse())
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

// ErrTerraformNotAvailable error used when `terraform` binary in not in the
// path and we try to use a provisioner.
var ErrTerraformNotAvailable = fmt.Errorf("provisioner error: executable file 'terraform' not found in $PATH")

// IsTerraformBinaryAvailable determines if Terraform binary is available in
// the path of the system. Terraform binary is a requirement for executing
// provisioner plugins, since they are built-in on the Terrafrom binary. :(
//
// https://github.com/hashicorp/terraform/issues/20896#issuecomment-479054649
func IsTerraformBinaryAvailable() bool {
	_, err := exec.LookPath("terraform")
	return err == nil
}
