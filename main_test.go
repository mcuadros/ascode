package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"

	"go.starlark.net/resolve"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarktest"
)

func init() {
	// The tests make extensive use of these not-yet-standard features.
	resolve.AllowLambda = true
	resolve.AllowNestedDef = true
	resolve.AllowFloat = true
	resolve.AllowSet = true
}

func TestProvider(t *testing.T) {
	test(t, "testdata/provider.star")
}

func TestNestedBlock(t *testing.T) {
	test(t, "testdata/nested.star")
}

func test(t *testing.T, filename string) {
	log.SetOutput(ioutil.Discard)
	thread := &starlark.Thread{Load: load}
	starlarktest.SetReporter(thread, t)

	provider := starlark.NewBuiltin("provider", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		name := args.Index(0).(starlark.String)
		version := args.Index(1).(starlark.String)
		return NewProviderInstance(&PluginManager{".providers"}, string(name), string(version))
	})

	predeclared := starlark.StringDict{
		"provider": provider,
	}

	if _, err := starlark.ExecFile(thread, filename, nil, predeclared); err != nil {
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

	return nil, fmt.Errorf("load not implemented")
}
