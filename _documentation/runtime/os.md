---
title: 'os'
---

os provides a platform-independent interface to operating system functionality.
## Functions

#### def <b>chdir</b>
```go
chdir(dir)
```
changes the current working directory to the named directory.

**parameters:**

| name | type | description |
|------|------|-------------|
| 'dir' | 'string' | target dir |


#### def <b>getenv</b>
```go
getenv(key) dir
```
retrieves the value of the environment variable named by the key.

**parameters:**

| name | type | description |
|------|------|-------------|
| 'key' | 'string' | name of the environment variable |


#### def <b>getwd</b>
```go
getwd() dir
```
returns a rooted path name corresponding to the current directory.

#### def <b>mkdir</b>
```go
mkdir(name, perms=0o777)
```
creates a new directory with the specified name and permission bits (before umask).

**parameters:**

| name | type | description |
|------|------|-------------|
| 'name' | 'string' | name of the folder to be created |
| 'perms' | 'int' | optional, permission of the folder |


#### def <b>mkdir_all</b>
```go
mkdir_all(name, perms=0o777)
```
creates a new directory with the specified name and permission bits (before umask).

**parameters:**

| name | type | description |
|------|------|-------------|
| 'name' | 'string' | name of the folder to be created |
| 'perms' | 'int' | optional, permission of the folder |


#### def <b>read_file</b>
```go
read_file(filename) string
```
reads the file named by filename and returns the contents.

**parameters:**

| name | type | description |
|------|------|-------------|
| 'filename' | 'string' | name of the file to be written |
| 'data' | 'string' | content to be witten to the file |
| 'perms' | 'int' | optional, permission of the file |


#### def <b>remove</b>
```go
remove(name)
```
removes the named file or (empty) directory.

**parameters:**

| name | type | description |
|------|------|-------------|
| 'name' | 'string' | name of the file or directory to be deleted |


#### def <b>remove_all</b>
```go
remove_all(path)
```
removes path and any children it contains. It removes everything it can but returns the first error it encounters.

**parameters:**

| name | type | description |
|------|------|-------------|
| 'name' | 'string' | path to be deleted |


#### def <b>rename</b>
```go
rename(oldpath, newpath)
```
renames (moves) oldpath to newpath. If newpath already exists and is not a directory, Rename replaces it. OS-specific restrictions may apply when oldpath and newpath are in different directories.

**parameters:**

| name | type | description |
|------|------|-------------|
| 'oldpath' | 'string' | old path |
| 'newpath' | 'string' | new path |


#### def <b>setenv</b>
```go
setenv(key, value) dir
```
sets the value of the environment variable named by the key.

**parameters:**

| name | type | description |
|------|------|-------------|
| 'key' | 'string' | name of the environment variable |
| 'value' | 'string' | value of the environment variable |


#### def <b>write_file</b>
```go
write_file(filename, data, perms=0o644)
```
retrieves the value of the environment variable named by the key.

**parameters:**

| name | type | description |
|------|------|-------------|
| 'filename' | 'string' | name of the file to be written |
| 'data' | 'string' | content to be witten to the file |
| 'perms' | 'int' | optional, permission of the file |



