package url

import (
	"path/filepath"
	"testing"

	"github.com/qri-io/starlib/testdata"
	"go.starlark.net/resolve"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarktest"
)

func TestFile(t *testing.T) {
	if filepath.Separator != '/' {
		// TODO(mcuadros): do proper testing on windows.
		t.Skip("skiping os test for Windows")
	}

	resolve.AllowFloat = true
	resolve.AllowGlobalReassign = true
	resolve.AllowLambda = true

	thread := &starlark.Thread{Load: testdata.NewLoader(LoadModule, ModuleName)}
	starlarktest.SetReporter(thread, t)

	// Execute test file
	_, err := starlark.ExecFile(thread, "testdata/test.star", nil, nil)
	if err != nil {
		if ee, ok := err.(*starlark.EvalError); ok {
			t.Error(ee.Backtrace())
		} else {
			t.Error(err)
		}
	}

}
