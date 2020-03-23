---
title: 'encoding/yaml'
---

yaml provides functions for working with yaml data

## Index


* [def <b>dumps</b>(obj) string](#def-iyamlibdumpsb)
* [def <b>loads</b>(source) object](#def-iyamlibloadsb)


## Functions


#### def <i>yaml</i>.<b>dumps</b>
```go
yaml.dumps(obj) string
```
serialize obj to a yaml string

###### Arguments

| name | type | description |
|------|------|-------------|
| `obj` | `object` | input object |



#### def <i>yaml</i>.<b>loads</b>
```go
yaml.loads(source) object
```
read a source yaml string to a starlark object

###### Arguments

| name | type | description |
|------|------|-------------|
| `source` | `string` | input string of yaml data |



