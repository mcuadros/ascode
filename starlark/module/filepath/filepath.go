package filepath 

import (
	"path/filepath"
	"sync"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

const (
	// ModuleName defines the expected name for this Module when used
	// in starlark's load() function, eg: load('io/ioutil', 'json')
	ModuleName = "path/filepath"

	SeparatorVarName = "separator"
	AbsFuncName      = "abs"
	BaseFuncName     = "base"
	CleanFuncName    = "clean"
	DirFuncName      = "dir"
	ExtFuncName      = "ext"
	GlobFuncName     = "glob"
	IsAbsFuncName    = "is_abs"
	JoinFuncName     = "join"
	RelFuncName      = "rel"
)

var (
	once         sync.Once
	ioutilModule starlark.StringDict
)

// LoadModule loads the os module.
// It is concurrency-safe and idempotent.
//
//   outline: filepath
//     filepath implements utility routines for manipulating filename paths in a
//     way compatible with the target operating system-defined file path
//     path: path/filepath
func LoadModule() (starlark.StringDict, error) {
	once.Do(func() {
		ioutilModule = starlark.StringDict{
			"filepath": &starlarkstruct.Module{
				Name: "filepath",
				Members: starlark.StringDict{
					SeparatorVarName: starlark.String(string(filepath.Separator)),
					GlobFuncName:     starlark.NewBuiltin(GlobFuncName, Glob),
					AbsFuncName:      starlark.NewBuiltin(AbsFuncName, Abs),
					BaseFuncName:     starlark.NewBuiltin(BaseFuncName, Base),
					CleanFuncName:    starlark.NewBuiltin(CleanFuncName, Clean),
					DirFuncName:      starlark.NewBuiltin(DirFuncName, Dir),
					ExtFuncName:      starlark.NewBuiltin(ExtFuncName, Ext),
					IsAbsFuncName:    starlark.NewBuiltin(IsAbsFuncName, IsAbs),
					JoinFuncName:     starlark.NewBuiltin(JoinFuncName, Join),
					RelFuncName:      starlark.NewBuiltin(RelFuncName, Rel),
				},
			},
		}
	})

	return ioutilModule, nil
}

// Glob returns the names of all files matching pattern or nil if there is no
// matching file.
//
//   outline: filepath
//     functions:
//       glob(pattern) list
//         returns the names of all files matching pattern or None if there is
//         no matching file.
//         params:
//           pattern string
//             pattern ([syntax](https://golang.org/pkg/path/filepath/#Match))
func Glob(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var pattern string

	err := starlark.UnpackArgs(GlobFuncName, args, kwargs, "pattern", &pattern)
	if err != nil {
		return nil, err
	}

	list, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	values := make([]starlark.Value, len(list))
	for i, entry := range list {
		values[i] = starlark.String(entry)
	}

	return starlark.NewList(values), nil
}

// Abs returns an absolute representation of path. If the path is not absolute
// it will be joined with the current working directory to turn it into an/
// absolute path.
//
//   outline: filepath
//     functions:
//       abs(path) string
//         returns an absolute representation of path. If the path is not
//         absolute it will be joined with the current working directory to turn
//         it into an absolute path.
//         params:
//           path string
//             relative or absolute path
func Abs(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var path string

	err := starlark.UnpackArgs(AbsFuncName, args, kwargs, "path", &path)
	if err != nil {
		return nil, err
	}

	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	return starlark.String(abs), nil
}

// Base returns the last element of path. Trailing path separators are removed
// before extracting the last element. If the path is empty, Base returns ".".
// If the path consists entirely of separators, Base returns a single separator.
//
//   outline: filepath
//     functions:
//       base(path) string
//         returns the last element of path. Trailing path separators are
//         removed before extracting the last element. If the path is empty,
//         `base` returns ".". If the path consists entirely of separators,
//         `base` returns a single separator.
//         params:
//           path string
//             input path
func Base(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var path string

	err := starlark.UnpackArgs(BaseFuncName, args, kwargs, "path", &path)
	if err != nil {
		return nil, err
	}

	return starlark.String(filepath.Base(path)), nil
}

// Clean returns the shortest path name equivalent to path by purely lexical processing.
//
//   outline: filepath
//     functions:
//       clean(path) string
//         returns the shortest path name equivalent to path by purely lexical
//         processing.
//         params:
//           path string
//             input path
func Clean(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var path string

	err := starlark.UnpackArgs(CleanFuncName, args, kwargs, "path", &path)
	if err != nil {
		return nil, err
	}

	return starlark.String(filepath.Clean(path)), nil
}

// Dir returns all but the last element of path, typically the path's directory.
//
//   outline: filepath
//     functions:
//       dir(path) string
//         returns all but the last element of path, typically the path's
//         directory. After dropping the final element, `dir` calls `clean` on the
//         path and trailing slashes are removed. If the path is empty, `dir`
//         returns ".". If the path consists entirely of separators, `dir`
//         returns a single separator. The returned path does not end in a
//         separator unless it is the root directory.
//         params:
//           path string
//             input path
func Dir(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var path string

	err := starlark.UnpackArgs(DirFuncName, args, kwargs, "path", &path)
	if err != nil {
		return nil, err
	}

	return starlark.String(filepath.Dir(path)), nil
}

// Ext returns the file name extension used by path. The extension is the suffix
// beginning at the final dot in the final element of path; it is empty if there
// is no dot.
//
//   outline: filepath
//     functions:
//       ext(path) string
//         returns the file name extension used by path. The extension is the
//         suffix beginning at the final dot in the final element of path; it
//         is empty if there is no dot.
//         params:
//           path string
//             input path
func Ext(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var path string

	err := starlark.UnpackArgs(ExtFuncName, args, kwargs, "path", &path)
	if err != nil {
		return nil, err
	}

	return starlark.String(filepath.Ext(path)), nil
}

// IsAbs reports whether the path is absolute.
//
//   outline: filepath
//     functions:
//       is_abs(path) bool
//         reports whether the path is absolute.
//         params:
//           path string
//             input path
func IsAbs(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var path string

	err := starlark.UnpackArgs(IsAbsFuncName, args, kwargs, "path", &path)
	if err != nil {
		return nil, err
	}

	return starlark.Bool(filepath.IsAbs(path)), nil
}

// Join joins any number of path elements into a single path, adding a Separator
// if necessary.
//
//   outline: filepath
//     functions:
//       join(elements) string
//         joins any number of path elements into a single path, adding a
//         `filepath.separator` if necessary. Join calls Clean on the result;
//         in particular, all empty strings are ignored. On Windows, the result
//         is a UNC path if and only if the first path element is a UNC path.
//         params:
//           elements lists
//             list of path elements to be joined
func Join(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var elements *starlark.List

	err := starlark.UnpackArgs(JoinFuncName, args, kwargs, "elements", &elements)
	if err != nil {
		return nil, err
	}

	parts := make([]string, elements.Len())
	for i := 0; i < cap(parts); i++ {
		parts[i] = elements.Index(i).(starlark.String).GoString()
	}

	return starlark.String(filepath.Join(parts...)), nil
}

// Rel returns a relative path that is lexically equivalent to targpath when
// joined to basepath with an intervening separator.
//
//   outline: filepath
//     functions:
//       rel(basepath, targpath) string
//         returns a relative path that is lexically equivalent to targpath when
//         joined to basepath with an intervening separator. That is, `filepath.join(basepath, filepath.rel(basepath, targpath))`
//         is equivalent to targpath itself. On success, the returned path will
//         always be relative to basepath, even if basepath and targpath share
//         no elements. An error is returned if targpath can't be made relative
//         to basepath or if knowing the current working directory would be
//         necessary to compute it. Rel calls Clean on the result.
//         params:
//           basepath string
//             relative or absolute path
//           targpath string
//             relative or absolute path
func Rel(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var basepath, targpath string

	err := starlark.UnpackArgs(RelFuncName, args, kwargs, "basepath", &basepath, "targpath", &targpath)
	if err != nil {
		return nil, err
	}

	rel, err := filepath.Rel(basepath, targpath)
	if err != nil {
		return nil, err
	}

	return starlark.String(rel), nil
}
