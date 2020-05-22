package runtime

import (
	"fmt"
	osfilepath "path/filepath"

	"github.com/mcuadros/ascode/starlark/module/docker"
	"github.com/mcuadros/ascode/starlark/module/filepath"
	"github.com/mcuadros/ascode/starlark/module/os"
	"github.com/mcuadros/ascode/starlark/types"
	"github.com/mcuadros/ascode/terraform"
	"github.com/qri-io/starlib/encoding/base64"
	"github.com/qri-io/starlib/encoding/csv"
	"github.com/qri-io/starlib/encoding/json"
	"github.com/qri-io/starlib/encoding/yaml"
	"github.com/qri-io/starlib/http"
	"github.com/qri-io/starlib/math"
	"github.com/qri-io/starlib/re"
	"github.com/qri-io/starlib/time"
	"go.starlark.net/repl"
	"go.starlark.net/resolve"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

func init() {
	resolve.AllowRecursion = true
	resolve.AllowFloat = true
	resolve.AllowGlobalReassign = true
	resolve.AllowLambda = true
	resolve.AllowNestedDef = true
	resolve.AllowSet = true
}

// LoadModuleFunc is a concurrency-safe and idempotent function that returns
// the module when is called from the `load` funcion.
type LoadModuleFunc func() (starlark.StringDict, error)

// Runtime represents the AsCode runtime, it defines the available modules,
// the predeclared globals and handles how the `load` function behaves.
type Runtime struct {
	Terraform   *types.Terraform
	pm          *terraform.PluginManager
	predeclared starlark.StringDict
	modules     map[string]LoadModuleFunc
	moduleCache map[string]*moduleCache

	path string
}

// NewRuntime returns a new Runtime for the given terraform.PluginManager.
func NewRuntime(pm *terraform.PluginManager) *Runtime {
	tf := types.NewTerraform(pm)
	predeclared := starlark.StringDict{}
	predeclared["tf"] = tf
	predeclared["provisioner"] = types.BuiltinProvisioner()
	predeclared["backend"] = types.BuiltinBackend()
	predeclared["validate"] = types.BuiltinValidate()
	predeclared["hcl"] = types.BuiltinHCL()
	predeclared["fn"] = types.BuiltinFunctionAttribute()
	predeclared["ref"] = types.BuiltinRef()
	predeclared["evaluate"] = types.BuiltinEvaluate(predeclared)
	predeclared["struct"] = starlark.NewBuiltin("struct", starlarkstruct.Make)
	predeclared["module"] = starlark.NewBuiltin("module", starlarkstruct.MakeModule)

	return &Runtime{
		Terraform:   tf,
		pm:          pm,
		moduleCache: make(map[string]*moduleCache),
		modules: map[string]LoadModuleFunc{
			filepath.ModuleName: filepath.LoadModule,
			os.ModuleName:       os.LoadModule,
			docker.ModuleName:   docker.LoadModule,

			"encoding/json":   json.LoadModule,
			"encoding/base64": base64.LoadModule,
			"encoding/csv":    csv.LoadModule,
			"encoding/yaml":   yaml.LoadModule,
			"math":            math.LoadModule,
			"re":              re.LoadModule,
			"time":            time.LoadModule,
			"http":            http.LoadModule,
		},
		predeclared: predeclared,
	}
}

// ExecFile parses, resolves, and executes a Starlark file.
func (r *Runtime) ExecFile(filename string) (starlark.StringDict, error) {
	fullpath, _ := osfilepath.Abs(filename)
	r.path, _ = osfilepath.Split(fullpath)

	thread := &starlark.Thread{Name: "thread", Load: r.load}
	r.setLocals(thread)

	return starlark.ExecFile(thread, filename, nil, r.predeclared)
}

// REPL executes a read, eval, print loop.
func (r *Runtime) REPL() {
	thread := &starlark.Thread{Name: "thread", Load: r.load}
	r.setLocals(thread)

	repl.REPL(thread, r.predeclared)
}

func (r *Runtime) setLocals(t *starlark.Thread) {
	t.SetLocal("base_path", r.path)
	t.SetLocal(types.PluginManagerLocal, r.pm)
}

func (r *Runtime) load(t *starlark.Thread, module string) (starlark.StringDict, error) {
	if m, ok := r.modules[module]; ok {
		return m()
	}

	filename := osfilepath.Join(r.path, module)
	return r.loadFile(t, filename)
}

type moduleCache struct {
	globals starlark.StringDict
	err     error
}

func (r *Runtime) loadFile(thread *starlark.Thread, module string) (starlark.StringDict, error) {
	e, ok := r.moduleCache[module]
	if e == nil {
		if ok {
			// request for package whose loading is in progress
			return nil, fmt.Errorf("cycle in load graph")
		}

		// Add a placeholder to indicate "load in progress".
		r.moduleCache[module] = nil

		thread := &starlark.Thread{Name: "exec " + module, Load: thread.Load}
		globals, err := starlark.ExecFile(thread, module, nil, r.predeclared)

		e = &moduleCache{globals, err}
		r.moduleCache[module] = e
	}

	return e.globals, e.err
}
