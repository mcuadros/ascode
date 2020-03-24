package runtime

import (
	"fmt"
	osfilepath "path/filepath"

	"github.com/mcuadros/ascode/starlark/module/experimental/docker"
	"github.com/mcuadros/ascode/starlark/module/filepath"
	"github.com/mcuadros/ascode/starlark/module/os"
	"github.com/mcuadros/ascode/starlark/types"
	"github.com/mcuadros/ascode/terraform"
	"github.com/qri-io/starlib/encoding/base64"
	"github.com/qri-io/starlib/encoding/csv"
	"github.com/qri-io/starlib/encoding/json"
	"github.com/qri-io/starlib/encoding/yaml"
	"github.com/qri-io/starlib/http"
	"github.com/qri-io/starlib/re"
	"go.starlark.net/repl"
	"go.starlark.net/resolve"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

func init() {
	resolve.AllowRecursion = true
	resolve.AllowFloat = true
	resolve.AllowGlobalReassign = true
}

type LoadModuleFunc func() (starlark.StringDict, error)

type Runtime struct {
	Terraform   *types.Terraform
	predeclared starlark.StringDict
	modules     map[string]LoadModuleFunc
	moduleCache map[string]*moduleCache

	path string
}

func NewRuntime(pm *terraform.PluginManager) *Runtime {
	tf := types.NewTerraform(pm)

	return &Runtime{
		Terraform:   tf,
		moduleCache: make(map[string]*moduleCache),
		modules: map[string]LoadModuleFunc{
			filepath.ModuleName: filepath.LoadModule,
			os.ModuleName:       os.LoadModule,
			docker.ModuleName:   docker.LoadModule,

			"encoding/json":   json.LoadModule,
			"encoding/base64": base64.LoadModule,
			"encoding/csv":    csv.LoadModule,
			"encoding/yaml":   yaml.LoadModule,
			"re":              re.LoadModule,
			"http":            http.LoadModule,
		},
		predeclared: starlark.StringDict{
			"tf":          tf,
			"provisioner": types.BuiltinProvisioner(pm),
			"backend":     types.BuiltinBackend(pm),
			"hcl":         types.BuiltinHCL(),
			"fn":          types.BuiltinFunctionComputed(),
			"evaluate":    types.BuiltinEvaluate(),
			"struct":      starlark.NewBuiltin("struct", starlarkstruct.Make),
			"module":      starlark.NewBuiltin("module", starlarkstruct.MakeModule),
		},
	}
}

func (r *Runtime) ExecFile(filename string) (starlark.StringDict, error) {
	filename, _ = osfilepath.Abs(filename)
	r.path, _ = osfilepath.Split(filename)

	thread := &starlark.Thread{Name: "thread", Load: r.load}
	thread.SetLocal("base_path", r.path)

	return starlark.ExecFile(thread, filename, nil, r.predeclared)
}

func (r *Runtime) REPL() {
	thread := &starlark.Thread{Name: "thread", Load: r.load}
	repl.REPL(thread, r.predeclared)
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
