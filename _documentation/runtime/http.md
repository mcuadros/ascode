---
title: ''
---

http defines an HTTP client implementation

## Index


* [def <b>delete</b>(url,params={},headers={},body="",form_body={},json_body={},auth=](#def-ihttpibdeleteb)
* [def <b>get</b>(url,params={},headers={},auth=](#def-ihttpibgetb)
* [def <b>options</b>(url,params={},headers={},body="",form_body={},json_body={},auth=](#def-ihttpiboptionsb)
* [def <b>patch</b>(url,params={},headers={},body="",form_body={},json_body={},auth=](#def-ihttpibpatchb)
* [def <b>post</b>(url,params={},headers={},body="",form_body={},json_body={},auth=](#def-ihttpibpostb)
* [def <b>put</b>(url,params={},headers={},body="",form_body={},json_body={},auth=](#def-ihttpibputb)
* [type <b>response</b>](#type-bresponseb)
    * [def <b>body</b>() string](#def-iresponseibbodyb)
    * [def <b>json</b>()](#def-iresponseibjsonb)


## Functions


#### def <i>http</i>.<b>delete</b>
```go
http.delete(url,params={},headers={},body="",form_body={},json_body={},auth=()) response
```
perform an HTTP DELETE request, returning a response

###### Arguments

| name | type | description |
|------|------|-------------|
| `url` | `string` | url to request |
| `headers` | `dict` | optional. dictionary of headers to add to request |
| `body` | `string` | optional. raw string body to provide to the request |
| `form_body` | `dict` | optional. dict of values that will be encoded as form data |
| `json_body` | `any` | optional. json data to supply as a request. handy for working with JSON-API's |
| `auth` | `tuple` | optional. (username,password) tuple for http basic authorization |



#### def <i>http</i>.<b>get</b>
```go
http.get(url,params={},headers={},auth=()) response
```
perform an HTTP GET request, returning a response

###### Arguments

| name | type | description |
|------|------|-------------|
| `url` | `string` | url to request |
| `headers` | `dict` | optional. dictionary of headers to add to request |
| `auth` | `tuple` | optional. (username,password) tuple for http basic authorization |



#### def <i>http</i>.<b>options</b>
```go
http.options(url,params={},headers={},body="",form_body={},json_body={},auth=()) response
```
perform an HTTP OPTIONS request, returning a response

###### Arguments

| name | type | description |
|------|------|-------------|
| `url` | `string` | url to request |
| `headers` | `dict` | optional. dictionary of headers to add to request |
| `body` | `string` | optional. raw string body to provide to the request |
| `form_body` | `dict` | optional. dict of values that will be encoded as form data |
| `json_body` | `any` | optional. json data to supply as a request. handy for working with JSON-API's |
| `auth` | `tuple` | optional. (username,password) tuple for http basic authorization |



#### def <i>http</i>.<b>patch</b>
```go
http.patch(url,params={},headers={},body="",form_body={},json_body={},auth=()) response
```
perform an HTTP PATCH request, returning a response

###### Arguments

| name | type | description |
|------|------|-------------|
| `url` | `string` | url to request |
| `headers` | `dict` | optional. dictionary of headers to add to request |
| `body` | `string` | optional. raw string body to provide to the request |
| `form_body` | `dict` | optional. dict of values that will be encoded as form data |
| `json_body` | `any` | optional. json data to supply as a request. handy for working with JSON-API's |
| `auth` | `tuple` | optional. (username,password) tuple for http basic authorization |



#### def <i>http</i>.<b>post</b>
```go
http.post(url,params={},headers={},body="",form_body={},json_body={},auth=()) response
```
perform an HTTP POST request, returning a response

###### Arguments

| name | type | description |
|------|------|-------------|
| `url` | `string` | url to request |
| `headers` | `dict` | optional. dictionary of headers to add to request |
| `body` | `string` | optional. raw string body to provide to the request |
| `form_body` | `dict` | optional. dict of values that will be encoded as form data |
| `json_body` | `any` | optional. json data to supply as a request. handy for working with JSON-API's |
| `auth` | `tuple` | optional. (username,password) tuple for http basic authorization |



#### def <i>http</i>.<b>put</b>
```go
http.put(url,params={},headers={},body="",form_body={},json_body={},auth=()) response
```
perform an HTTP PUT request, returning a response

###### Arguments

| name | type | description |
|------|------|-------------|
| `url` | `string` | url to request |
| `headers` | `dict` | optional. dictionary of headers to add to request |
| `body` | `string` | optional. raw string body to provide to the request |
| `form_body` | `dict` | optional. dict of values that will be encoded as form data |
| `json_body` | `any` | optional. json data to supply as a request. handy for working with JSON-API's |
| `auth` | `tuple` | optional. (username,password) tuple for http basic authorization |




## Types
### type <b>response</b>
the result of performing a http request

###### Properties

| name | type | description |
|------|------|-------------|
| `url` | `string` | the url that was ultimately requested (may change after redirects) |
| `status_code` | `int` | response status code (for example: 200 == OK) |
| `headers` | `dict` | dictionary of response headers |
| `encoding` | `string` | transfer encoding. example: "octet-stream" or "application/json" |




###### Methods

#### def <i>response</i>.<b>body</b>
```go
response.body() string
```
output response body as a string


#### def <i>response</i>.<b>json</b>
```go
response.json()
```
attempt to parse resonse body as json, returning a JSON-decoded result


