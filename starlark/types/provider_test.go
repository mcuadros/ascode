package types

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"

	"github.com/ascode-dev/ascode/starlark/module/os"
	"github.com/ascode-dev/ascode/terraform"
	"go.starlark.net/resolve"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarktest"
)

var id int

func init() {
	NameGenerator = func() string {
		id++
		return fmt.Sprintf("id_%d", id)
	}

	// The tests make extensive use of these not-yet-standard features.
	resolve.AllowLambda = true
	resolve.AllowNestedDef = true
	resolve.AllowFloat = true
	resolve.AllowSet = true
}

func TestProvider(t *testing.T) {
	test(t, "testdata/provider.star")
}

func TestProvisioner(t *testing.T) {
	test(t, "testdata/provisioner.star")
}

func TestNestedBlock(t *testing.T) {
	test(t, "testdata/nested.star")
}

func TestResource(t *testing.T) {
	test(t, "testdata/resource.star")
}

func TestHCL(t *testing.T) {
	test(t, "testdata/hcl.star")
}

func test(t *testing.T, filename string) {
	id = 0

	log.SetOutput(ioutil.Discard)
	thread := &starlark.Thread{Load: load}
	starlarktest.SetReporter(thread, t)

	pm := &terraform.PluginManager{".providers"}

	predeclared := starlark.StringDict{
		"provider":    BuiltinProvider(pm),
		"provisioner": BuiltinProvisioner(pm),
		"hcl":         BuiltinHCL(),
	}

	_, err := starlark.ExecFile(thread, filename, nil, predeclared)
	if err != nil {
		if err, ok := err.(*starlark.EvalError); ok {
			t.Fatal(err.Backtrace())
		}
		t.Fatal(err)
	}
}

// load implements the 'load' operation as used in the evaluator tests.
func load(thread *starlark.Thread, module string) (starlark.StringDict, error) {
	if module == "assert.star" {
		return starlarktest.LoadAssertModule()
	}

	if module == os.ModuleName {
		return os.LoadModule()
	}

	return nil, fmt.Errorf("load not implemented")
}
