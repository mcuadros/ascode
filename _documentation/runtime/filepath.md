---
title: 'path/filepath'
---

filepath implements utility routines for manipulating filename paths in a
way compatible with the target operating system-defined file path

## Index


* [def <b>abs</b>(path) string](#def-ifilepathibabsb)
* [def <b>base</b>(path) string](#def-ifilepathibbaseb)
* [def <b>clean</b>(path) string](#def-ifilepathibcleanb)
* [def <b>dir</b>(path) string](#def-ifilepathibdirb)
* [def <b>ext</b>(path) string](#def-ifilepathibextb)
* [def <b>glob</b>(pattern) list](#def-ifilepathibglobb)
* [def <b>is_abs</b>(path) bool](#def-ifilepathibis_absb)
* [def <b>join</b>(elements) string](#def-ifilepathibjoinb)
* [def <b>rel</b>(basepath, targpath) string](#def-ifilepathibrelb)


## Functions


#### def <i>filepath</i>.<b>abs</b>
```go
filepath.abs(path) string
```
returns an absolute representation of path. If the path is not
absolute it will be joined with the current working directory to turn
it into an absolute path.

###### Arguments

| name | type | description |
|------|------|-------------|
| `path` | `string` | relative or absolute path |



#### def <i>filepath</i>.<b>base</b>
```go
filepath.base(path) string
```
returns the last element of path. Trailing path separators are
removed before extracting the last element. If the path is empty,
`base` returns ".". If the path consists entirely of separators,
`base` returns a single separator.

###### Arguments

| name | type | description |
|------|------|-------------|
| `path` | `string` | input path |



#### def <i>filepath</i>.<b>clean</b>
```go
filepath.clean(path) string
```
returns the shortest path name equivalent to path by purely lexical
processing.

###### Arguments

| name | type | description |
|------|------|-------------|
| `path` | `string` | input path |



#### def <i>filepath</i>.<b>dir</b>
```go
filepath.dir(path) string
```
returns all but the last element of path, typically the path's
directory. After dropping the final element, `dir` calls `clean` on the
path and trailing slashes are removed. If the path is empty, `dir`
returns ".". If the path consists entirely of separators, `dir`
returns a single separator. The returned path does not end in a
separator unless it is the root directory.

###### Arguments

| name | type | description |
|------|------|-------------|
| `path` | `string` | input path |



#### def <i>filepath</i>.<b>ext</b>
```go
filepath.ext(path) string
```
returns the file name extension used by path. The extension is the
suffix beginning at the final dot in the final element of path; it
is empty if there is no dot.

###### Arguments

| name | type | description |
|------|------|-------------|
| `path` | `string` | input path |



#### def <i>filepath</i>.<b>glob</b>
```go
filepath.glob(pattern) list
```
returns the names of all files matching pattern or None if there is
no matching file.

###### Arguments

| name | type | description |
|------|------|-------------|
| `pattern` | `string` | pattern ([syntax](https://golang.org/pkg/path/filepath/#Match)) |



#### def <i>filepath</i>.<b>is_abs</b>
```go
filepath.is_abs(path) bool
```
reports whether the path is absolute.

###### Arguments

| name | type | description |
|------|------|-------------|
| `path` | `string` | input path |



#### def <i>filepath</i>.<b>join</b>
```go
filepath.join(elements) string
```
joins any number of path elements into a single path, adding a
`filepath.separator` if necessary. Join calls Clean on the result;
in particular, all empty strings are ignored. On Windows, the result
is a UNC path if and only if the first path element is a UNC path.

###### Arguments

| name | type | description |
|------|------|-------------|
| `elements` | `lists` | list of path elements to be joined |



#### def <i>filepath</i>.<b>rel</b>
```go
filepath.rel(basepath, targpath) string
```
returns a relative path that is lexically equivalent to targpath when
joined to basepath with an intervening separator. That is, `filepath.join(basepath, filepath.rel(basepath, targpath))`
is equivalent to targpath itself. On success, the returned path will
always be relative to basepath, even if basepath and targpath share
no elements. An error is returned if targpath can't be made relative
to basepath or if knowing the current working directory would be
necessary to compute it. Rel calls Clean on the result.

###### Arguments

| name | type | description |
|------|------|-------------|
| `basepath` | `string` | relative or absolute path |
| `targpath` | `string` | relative or absolute path |



