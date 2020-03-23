---
title: 'encoding/json'
---

json provides functions for working with json data

## Index


* [def <b>dumps</b>(obj) string](#def-ijsonibdumpsb)
* [def <b>loads</b>(source) object](#def-ijsonibloadsb)


## Functions


#### def <i>json</i>.<b>dumps</b>
```go
json.dumps(obj) string
```
serialize obj to a JSON string

###### Arguments

| name | type | description |
|------|------|-------------|
| `obj` | `object` | input object |



#### def <i>json</i>.<b>loads</b>
```go
json.loads(source) object
```
read a source JSON string to a starlark object

###### Arguments

| name | type | description |
|------|------|-------------|
| `source` | `string` | input string of json data |



