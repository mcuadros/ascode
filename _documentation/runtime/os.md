# os
os provides a platform-independent interface to operating system functionality.

## Functions

#### `chdir(dir)`
changes the current working directory to the named directory.

**parameters:**

| name | type | description |
|------|------|-------------|
| `dir` | `string` | target dir |


#### `getenv(key) dir`
retrieves the value of the environment variable named by the key.

**parameters:**

| name | type | description |
|------|------|-------------|
| `key` | `string` | name of the environment variable |


#### `getwd() dir`
returns a rooted path name corresponding to the current directory.

#### `mkdir(name, perms=0o777)`
creates a new directory with the specified name and permission bits (before umask).

**parameters:**

| name | type | description |
|------|------|-------------|
| `name` | `string` | name of the folder to be created |
| `perms` | `int` | optional, permission of the folder |


#### `mkdir_all(name, perms=0o777)`
creates a new directory with the specified name and permission bits (before umask).

**parameters:**

| name | type | description |
|------|------|-------------|
| `name` | `string` | name of the folder to be created |
| `perms` | `int` | optional, permission of the folder |


#### `read_file(filename) string`
reads the file named by filename and returns the contents.

**parameters:**

| name | type | description |
|------|------|-------------|
| `filename` | `string` | name of the file to be written |
| `data` | `string` | content to be witten to the file |
| `perms` | `int` | optional, permission of the file |


#### `remove(name)`
removes the named file or (empty) directory.

**parameters:**

| name | type | description |
|------|------|-------------|
| `name` | `string` | name of the file or directory to be deleted |


#### `remove_all(path)`
removes path and any children it contains. It removes everything it can but returns the first error it encounters.

**parameters:**

| name | type | description |
|------|------|-------------|
| `name` | `string` | path to be deleted |


#### `rename(oldpath, newpath)`
renames (moves) oldpath to newpath. If newpath already exists and is not a directory, Rename replaces it. OS-specific restrictions may apply when oldpath and newpath are in different directories.

**parameters:**

| name | type | description |
|------|------|-------------|
| `oldpath` | `string` | old path |
| `newpath` | `string` | new path |


#### `setenv(key, value) dir`
sets the value of the environment variable named by the key.

**parameters:**

| name | type | description |
|------|------|-------------|
| `key` | `string` | name of the environment variable |
| `value` | `string` | value of the environment variable |


#### `write_file(filename, data, perms=0o644)`
retrieves the value of the environment variable named by the key.

**parameters:**

| name | type | description |
|------|------|-------------|
| `filename` | `string` | name of the file to be written |
| `data` | `string` | content to be witten to the file |
| `perms` | `int` | optional, permission of the file |


