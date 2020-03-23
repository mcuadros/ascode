---
title: 'os'
---

os provides a platform-independent interface to operating system functionality.

## Index


* [def <b>chdir</b>(dir)](#def-iosibchdirb)
* [def <b>getenv</b>(key) dir](#def-iosibgetenvb)
* [def <b>getwd</b>() dir](#def-iosibgetwdb)
* [def <b>mkdir</b>(name, perms=0o777)](#def-iosibmkdirb)
* [def <b>mkdir_all</b>(name, perms=0o777)](#def-iosibmkdir_allb)
* [def <b>read_file</b>(filename) string](#def-iosibread_fileb)
* [def <b>remove</b>(name)](#def-iosibremoveb)
* [def <b>remove_all</b>(path)](#def-iosibremove_allb)
* [def <b>rename</b>(oldpath, newpath)](#def-iosibrenameb)
* [def <b>setenv</b>(key, value) dir](#def-iosibsetenvb)
* [def <b>write_file</b>(filename, data, perms=0o644)](#def-iosibwrite_fileb)


## Functions


#### def <i>os</i>.<b>chdir</b>
```go
os.chdir(dir)
```
changes the current working directory to the named directory.

###### Arguments

| name | type | description |
|------|------|-------------|
| `dir` | `string` | target dir |



#### def <i>os</i>.<b>getenv</b>
```go
os.getenv(key) dir
```
retrieves the value of the environment variable named by the key.

###### Arguments

| name | type | description |
|------|------|-------------|
| `key` | `string` | name of the environment variable |



#### def <i>os</i>.<b>getwd</b>
```go
os.getwd() dir
```
returns a rooted path name corresponding to the current directory.


#### def <i>os</i>.<b>mkdir</b>
```go
os.mkdir(name, perms=0o777)
```
creates a new directory with the specified name and permission bits (before umask).

###### Arguments

| name | type | description |
|------|------|-------------|
| `name` | `string` | name of the folder to be created |
| `perms` | `int` | optional, permission of the folder |



#### def <i>os</i>.<b>mkdir_all</b>
```go
os.mkdir_all(name, perms=0o777)
```
creates a new directory with the specified name and permission bits (before umask).

###### Arguments

| name | type | description |
|------|------|-------------|
| `name` | `string` | name of the folder to be created |
| `perms` | `int` | optional, permission of the folder |



#### def <i>os</i>.<b>read_file</b>
```go
os.read_file(filename) string
```
reads the file named by filename and returns the contents.

###### Arguments

| name | type | description |
|------|------|-------------|
| `filename` | `string` | name of the file to be written |
| `data` | `string` | content to be witten to the file |
| `perms` | `int` | optional, permission of the file |



#### def <i>os</i>.<b>remove</b>
```go
os.remove(name)
```
removes the named file or (empty) directory.

###### Arguments

| name | type | description |
|------|------|-------------|
| `name` | `string` | name of the file or directory to be deleted |



#### def <i>os</i>.<b>remove_all</b>
```go
os.remove_all(path)
```
removes path and any children it contains. It removes everything it can but returns the first error it encounters.

###### Arguments

| name | type | description |
|------|------|-------------|
| `name` | `string` | path to be deleted |



#### def <i>os</i>.<b>rename</b>
```go
os.rename(oldpath, newpath)
```
renames (moves) oldpath to newpath. If newpath already exists and is not a directory, Rename replaces it. OS-specific restrictions may apply when oldpath and newpath are in different directories.

###### Arguments

| name | type | description |
|------|------|-------------|
| `oldpath` | `string` | old path |
| `newpath` | `string` | new path |



#### def <i>os</i>.<b>setenv</b>
```go
os.setenv(key, value) dir
```
sets the value of the environment variable named by the key.

###### Arguments

| name | type | description |
|------|------|-------------|
| `key` | `string` | name of the environment variable |
| `value` | `string` | value of the environment variable |



#### def <i>os</i>.<b>write_file</b>
```go
os.write_file(filename, data, perms=0o644)
```
retrieves the value of the environment variable named by the key.

###### Arguments

| name | type | description |
|------|------|-------------|
| `filename` | `string` | name of the file to be written |
| `data` | `string` | content to be witten to the file |
| `perms` | `int` | optional, permission of the file |



