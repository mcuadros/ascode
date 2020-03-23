---
title: 'encoding/base64'
---

base64 defines base64 encoding & decoding functions, often used to represent binary as text.

## Index


* [def <b>decode</b>(src,encoding="standard") string](#def-ibase64ibdecodeb)
* [def <b>encode</b>(src,encoding="standard") string](#def-ibase64ibencodeb)


## Functions


#### def <i>base64</i>.<b>decode</b>
```go
base64.decode(src,encoding="standard") string
```
parse base64 input, giving back the plain string representation

###### Arguments

| name | type | description |
|------|------|-------------|
| `src` | `string` | source string of base64-encoded text |
| `encoding` | `string` | optional. string to set decoding dialect. allowed values are: standard,standard_raw,url,url_raw |



#### def <i>base64</i>.<b>encode</b>
```go
base64.encode(src,encoding="standard") string
```
return the base64 encoding of src

###### Arguments

| name | type | description |
|------|------|-------------|
| `src` | `string` | source string to encode to base64 |
| `encoding` | `string` | optional. string to set encoding dialect. allowed values are: standard,standard_raw,url,url_raw |



