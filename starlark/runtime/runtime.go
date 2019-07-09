package runtime

import (
	"github.com/ascode-dev/ascode/starlark/module/filepath"
	"github.com/ascode-dev/ascode/starlark/module/os"
	"github.com/ascode-dev/ascode/starlark/types"
	"github.com/ascode-dev/ascode/terraform"
	"github.com/qri-io/starlib/encoding/base64"
	"github.com/qri-io/starlib/encoding/csv"
	"github.com/qri-io/starlib/encoding/json"
	"github.com/qri-io/starlib/http"
	"github.com/qri-io/starlib/re"
	"go.starlark.net/repl"
	"go.starlark.net/resolve"
	"go.starlark.net/starlark"
)

func init() {
	resolve.AllowFloat = true
}

type LoadModuleFunc func() (starlark.StringDict, error)

type Runtime struct {
	predeclared  starlark.StringDict
	modules      map[string]LoadModuleFunc
	fallbackLoad func(t *starlark.Thread, module string) (starlark.StringDict, error)
}

func NewRuntime(pm *terraform.PluginManager) *Runtime {
	return &Runtime{
		fallbackLoad: repl.MakeLoad(),
		modules: map[string]LoadModuleFunc{
			filepath.ModuleName: filepath.LoadModule,
			os.ModuleName:       os.LoadModule,

			"encoding/json":   json.LoadModule,
			"encoding/base64": base64.LoadModule,
			"encoding/csv":    csv.LoadModule,
			"re":              re.LoadModule,
			"http":            http.LoadModule,
		},
		predeclared: starlark.StringDict{
			"provider": types.BuiltinProvider(pm),
			"hcl":      types.BuiltinHCL(),
		},
	}
}

func (r *Runtime) ExecFile(filename string) (starlark.StringDict, error) {
	thread := &starlark.Thread{Name: "thread", Load: r.load}
	return starlark.ExecFile(thread, filename, nil, r.predeclared)
}

func (r *Runtime) load(t *starlark.Thread, module string) (starlark.StringDict, error) {
	if m, ok := r.modules[module]; ok {
		return m()
	}

	return nil, nil
}
