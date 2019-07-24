package terraform

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/plugin/discovery"
	"github.com/stretchr/testify/assert"
)

func TestPluginManager_Provider(t *testing.T) {
	path, err := ioutil.TempDir("", "provider")
	assert.NoError(t, err)

	pm := &PluginManager{Path: path}
	cli, meta, err := pm.Provider("github", "2.1.0", false)
	assert.NoError(t, err)
	assert.NotNil(t, cli)
	assert.Equal(t, meta.Version, discovery.VersionStr("2.1.0"))

	cli, meta, err = pm.Provider("github", "2.1.0", true)
	assert.NoError(t, err)
	assert.NotNil(t, cli)
	assert.Equal(t, meta.Version, discovery.VersionStr("2.1.0"))

	cli, meta, err = pm.Provider("github", "2.1.0", false)
	assert.NoError(t, err)
	assert.NotNil(t, cli)
	assert.Equal(t, meta.Version, discovery.VersionStr("2.1.0"))
}

func TestPluginManager_ProviderDefault(t *testing.T) {
	path, err := ioutil.TempDir("", "provider")
	assert.NoError(t, err)

	pm := &PluginManager{Path: path}
	cli, meta, err := pm.Provider("github", "", false)
	assert.NoError(t, err)
	assert.NotNil(t, cli)
	assert.NotEqual(t, meta.Version, discovery.VersionStr("2.1.0"))

	cli, meta, err = pm.Provider("github", "", true)
	assert.NoError(t, err)
	assert.NotNil(t, cli)
	assert.NotEqual(t, meta.Version, discovery.VersionStr("2.1.0"))

	fmt.Println(meta.Path)
	assert.Equal(t, strings.Index(meta.Path, path), 0)
}

func TestPluginManager_ProvisionerDefault(t *testing.T) {
	path, err := ioutil.TempDir("", "provisioner")
	assert.NoError(t, err)

	pm := &PluginManager{Path: path}
	cli, meta, err := pm.Provisioner("file")
	assert.NoError(t, err)

	assert.NotNil(t, cli)
	assert.Equal(t, strings.Index(meta.Path, "terraform-TFSPACE-internal-plugin-"), 0)

}
