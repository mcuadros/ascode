package types

import (
	"fmt"
	"io/ioutil"
	"log"
	stdos "os"
	"path/filepath"
	"testing"

	"github.com/mcuadros/ascode/starlark/module/os"
	"github.com/mcuadros/ascode/starlark/test"
	"github.com/mcuadros/ascode/terraform"
	"go.starlark.net/resolve"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

var id int

func init() {
	NameGenerator = func() string {
		id++
		return fmt.Sprintf("id_%d", id)
	}

	resolve.AllowLambda = true
	resolve.AllowNestedDef = true
	resolve.AllowFloat = true
	resolve.AllowSet = true
	resolve.AllowGlobalReassign = true
}

func TestProvider(t *testing.T) {
	doTest(t, "testdata/provider.star")
}

func TestProvisioner(t *testing.T) {
	if stdos.Getenv("ALLOW_PROVISIONER_SKIP") != "" && !terraform.IsTerraformBinaryAvailable() {
		t.Skip("terraform binary now available in $PATH")
	}

	doTest(t, "testdata/provisioner.star")
}

func TestNestedBlock(t *testing.T) {
	doTest(t, "testdata/nested.star")
}

func TestResource(t *testing.T) {
	doTest(t, "testdata/resource.star")
}

func TestHCL(t *testing.T) {
	doTest(t, "testdata/hcl.star")
}

func TestHCLIntegration(t *testing.T) {
	doTest(t, "testdata/hcl_integration.star")
}

func doTest(t *testing.T, filename string) {
	doTestPrint(t, filename, nil)
}

func doTestPrint(t *testing.T, filename string, print func(*starlark.Thread, string)) {
	id = 0

	dir, _ := filepath.Split(filename)
	pm := &terraform.PluginManager{".providers"}

	log.SetOutput(ioutil.Discard)
	thread := &starlark.Thread{Load: load, Print: print}
	thread.SetLocal("base_path", dir)
	thread.SetLocal(PluginManagerLocal, pm)

	test.SetReporter(thread, t)

	predeclared := starlark.StringDict{}
	predeclared["tf"] = NewTerraform(pm)
	predeclared["provisioner"] = BuiltinProvisioner()
	predeclared["backend"] = BuiltinBackend()
	predeclared["hcl"] = BuiltinHCL()
	predeclared["validate"] = BuiltinValidate()
	predeclared["fn"] = BuiltinFunctionAttribute()
	predeclared["evaluate"] = BuiltinEvaluate(predeclared)
	predeclared["struct"] = starlark.NewBuiltin("struct", starlarkstruct.Make)
	predeclared["module"] = starlark.NewBuiltin("module", starlarkstruct.MakeModule)

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
		return test.LoadAssertModule()
	}

	if module == os.ModuleName {
		return os.LoadModule()
	}

	return nil, fmt.Errorf("load not implemented")
}
