---
title: 'encoding/json'
---

json provides functions for working with json data
## Functions

#### def <b>dumps</b>
```go
dumps(obj) string
```
serialize obj to a JSON string

**parameters:**

| name | type | description |
|------|------|-------------|
| 'obj' | 'object' | input object |


#### def <b>loads</b>
```go
loads(source) object
```
read a source JSON string to a starlark object

**parameters:**

| name | type | description |
|------|------|-------------|
| 'source' | 'string' | input string of json data |



