package os

import (
	"io/ioutil"
	"os"
	"sync"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

const (
	// ModuleName defines the expected name for this Module when used
	// in starlark's load() function, eg: load('io/ioutil', 'json')
	ModuleName = "os"

	GetwdFuncName     = "getwd"
	ChdirFuncName     = "chdir"
	GetenvFuncName    = "getenv"
	SetenvFuncName    = "setenv"
	WriteFileFuncName = "write_file"
	ReadFileFuncName  = "read_file"
	MkdirFuncName     = "mkdir"
	MkdirAllFuncName  = "mkdir_all"
	RemoveFuncName    = "remove"
	RemoveAllFuncName = "remove_all"
	RenameFuncName    = "rename"
)

var (
	once         sync.Once
	ioutilModule starlark.StringDict
)

// LoadModule loads the os module.
// It is concurrency-safe and idempotent.
//
//   outline: os
//     os provides a platform-independent interface to operating system functionality.
//     path: os
func LoadModule() (starlark.StringDict, error) {
	once.Do(func() {
		ioutilModule = starlark.StringDict{
			"os": &starlarkstruct.Module{
				Name: "os",
				Members: starlark.StringDict{
					ChdirFuncName:     starlark.NewBuiltin(ChdirFuncName, Chdir),
					GetwdFuncName:     starlark.NewBuiltin(GetwdFuncName, Getwd),
					SetenvFuncName:    starlark.NewBuiltin(SetenvFuncName, Setenv),
					GetenvFuncName:    starlark.NewBuiltin(GetenvFuncName, Getenv),
					WriteFileFuncName: starlark.NewBuiltin(WriteFileFuncName, WriteFile),
					ReadFileFuncName:  starlark.NewBuiltin(ReadFileFuncName, ReadFile),
					MkdirFuncName:     starlark.NewBuiltin(MkdirFuncName, Mkdir),
					MkdirAllFuncName:  starlark.NewBuiltin(MkdirAllFuncName, MkdirAll),
					RemoveFuncName:    starlark.NewBuiltin(MkdirFuncName, Remove),
					RemoveAllFuncName: starlark.NewBuiltin(MkdirFuncName, RemoveAll),
					RenameFuncName:    starlark.NewBuiltin(RenameFuncName, Rename),
				},
			},
		}
	})

	return ioutilModule, nil
}

// Chdir changes the current working directory to the named directory.
//
//   outline: os
//     functions:
//       chdir(dir)
//         changes the current working directory to the named directory.
//         params:
//           dir string
//             target dir
func Chdir(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var dir string

	err := starlark.UnpackArgs(ChdirFuncName, args, kwargs, "dir", &dir)
	if err != nil {
		return nil, err
	}

	return starlark.None, os.Chdir(dir)
}

// Getwd returns a rooted path name corresponding to the current directory.
//
//   outline: os
//     functions:
//       getwd() dir
//         returns a rooted path name corresponding to the current directory.
func Getwd(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	dir, err := os.Getwd()
	return starlark.String(dir), err
}

// Setenv sets the value of the environment variable named by the key. It returns an error, if any.
//
//   outline: os
//     functions:
//       setenv(key, value) dir
//         sets the value of the environment variable named by the key.
//         params:
//           key string
//             name of the environment variable
//           value string
//             value of the environment variable
func Setenv(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		key   string
		value string
	)

	err := starlark.UnpackArgs(SetenvFuncName, args, kwargs, "key", &key, "value", &value)
	if err != nil {
		return nil, err
	}

	return starlark.None, os.Setenv(key, value)
}

// Getenv retrieves the value of the environment variable named by the key.
//
//   outline: os
//     functions:
//       getenv(key) dir
//         retrieves the value of the environment variable named by the key.
//         params:
//           key string
//             name of the environment variable
func Getenv(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		key string
		def string
	)

	err := starlark.UnpackArgs(GetenvFuncName, args, kwargs, "key", &key, "default?", &def)
	if err != nil {
		return nil, err
	}

	value := os.Getenv(key)
	if value == "" {
		value = def
	}

	return starlark.String(value), nil
}

// WriteFile writes data to a file named by filename. If the file does not
// exist, WriteFile creates it with permissions perm; otherwise WriteFile
// truncates it before writing.
//
//   outline: os
//     functions:
//       write_file(filename, data, perms=0o644)
//         retrieves the value of the environment variable named by the key.
//         params:
//           filename string
//             name of the file to be written
//           data string
//              content to be witten to the file
//           perms int
//              optional, permission of the file
func WriteFile(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		filename string
		content  string
		perms    = 0644
	)

	err := starlark.UnpackArgs(WriteFileFuncName, args, kwargs, "filename", &filename, "content", &content, "perms?", &perms)
	if err != nil {
		return nil, err
	}

	return starlark.None, ioutil.WriteFile(filename, []byte(content), os.FileMode(perms))
}

// ReadFile reads the file named by filename and returns the contents.
//
//   outline: os
//     functions:
//       read_file(filename) string
//         reads the file named by filename and returns the contents.
//         params:
//           filename string
//             name of the file to be written
//           data string
//              content to be witten to the file
//           perms int
//              optional, permission of the file
func ReadFile(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var filename string

	err := starlark.UnpackArgs(ReadFileFuncName, args, kwargs, "filename", &filename)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return starlark.String(string(data)), nil
}

// Mkdir creates a new directory with the specified name and permission bits (before umask).
//
//   outline: os
//     functions:
//       mkdir(name, perms=0o777)
//         creates a new directory with the specified name and permission bits (before umask).
//         params:
//           name string
//             name of the folder to be created
//           perms int
//              optional, permission of the folder
func Mkdir(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		name  string
		perms = 0777
	)

	err := starlark.UnpackArgs(MkdirFuncName, args, kwargs, "name", &name, "perms?", &perms)
	if err != nil {
		return nil, err
	}

	return starlark.None, os.Mkdir(name, os.FileMode(perms))
}

// MkdirAll creates a directory named path, along with any necessary parents.
//
//   outline: os
//     functions:
//       mkdir_all(name, perms=0o777)
//         creates a new directory with the specified name and permission bits (before umask).
//         params:
//           name string
//             name of the folder to be created
//           perms int
//              optional, permission of the folder
func MkdirAll(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		path  string
		perms = 0777
	)

	err := starlark.UnpackArgs(MkdirAllFuncName, args, kwargs, "path", &path, "perms?", &perms)
	if err != nil {
		return nil, err
	}

	return starlark.None, os.MkdirAll(path, os.FileMode(perms))
}

// Remove removes the named file or (empty) directory.
//
//   outline: os
//     functions:
//       remove(name)
//         removes the named file or (empty) directory.
//         params:
//           name string
//             name of the file or directory to be deleted
func Remove(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var name string

	err := starlark.UnpackArgs(RemoveFuncName, args, kwargs, "name", &name)
	if err != nil {
		return nil, err
	}

	return starlark.None, os.Remove(name)
}

// RemoveAll removes path and any children it contains.
//
//   outline: os
//     functions:
//       remove_all(path)
//         removes path and any children it contains. It removes everything it
//         can but returns the first error it encounters.
//         params:
//           name string
//             path to be deleted
func RemoveAll(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var path string

	err := starlark.UnpackArgs(RemoveAllFuncName, args, kwargs, "path", &path)
	if err != nil {
		return nil, err
	}

	return starlark.None, os.RemoveAll(path)
}

// Rename renames (moves) oldpath to newpath. If
//
//   outline: os
//     functions:
//       rename(oldpath, newpath)
//         renames (moves) oldpath to newpath. If newpath already exists and is
//         not a directory, Rename replaces it. OS-specific restrictions may
//         apply when oldpath and newpath are in different directories.
//         params:
//           oldpath string
//             old path
//           newpath string
//             new path
func Rename(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		oldpath string
		newpath string
	)

	err := starlark.UnpackArgs(RenameFuncName, args, kwargs, "oldpath", &oldpath, "newpath", &newpath)
	if err != nil {
		return nil, err
	}

	return starlark.None, os.Rename(oldpath, newpath)
}
