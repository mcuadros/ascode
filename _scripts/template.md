{{- define "functionName" }}
{{- if ne .Receiver "types" -}}
#
{{- end -}}
### def {{ if ne .Receiver "types" -}}
	 <i>{{ .Receiver }}</i>.
	{{- end -}}
	<b>{{- (index (split .Signature "(") 0) -}}</b>
{{- end -}}
{{- define "function" }}

{{ template "functionName" . }}
```go
{{if ne .Receiver "types" -}}{{.Receiver}}.{{- end }}{{ .Signature }}
```

{{- if ne .Description "" }}
{{ .Description }}
{{- end -}}

{{- if gt (len .Params) 0 }}

###### Arguments

| name | type | description |
|------|------|-------------|
{{ range .Params -}}
| `{{ .Name }}` | `{{ .Type }}` | {{ .Description }} |
{{ end -}}

{{- end -}}

{{- if gt (len .Examples) 0 }}
###### Examples
{{ range .Examples -}}
{{ .Description }}
```python
{{ .Code }}
```
{{ end -}}

{{- end -}}

{{- end -}}
{{- define "indexFunctionName" -}}
{{- $receiver := "" -}}
{{ if ne .Receiver "types" -}}{{ $receiver = printf "<i>%s</i>." .Receiver }}{{- end -}}
{{- $name := printf "def <b>%s</b>" (index (split .Signature "(") 0) -}}
{{- $anchor := printf "def %s<b>%s</b>" $receiver (index (split .Signature "(") 0) -}}
* [{{ $name }}({{ (index (split .Signature "(") 1) }}](#{{ sanitizeAnchor $anchor }})
{{- end -}}

{{- define "index" -}}
## Index

{{ range .Functions }}
{{ template "indexFunctionName" . }}
{{- end -}}
{{ range .Types -}}
{{ $name := printf "type <b>%s</b>" .Name }}
* [{{ $name }}](#{{ sanitizeAnchor $name }})
{{- range .Methods }}
    {{ template "indexFunctionName" . -}}
{{- end -}}
{{- end }}

{{ end -}}


{{- range . -}}
---
title: '{{ .Path }}'
---

{{ if ne .Description "" }}{{ .Description }}{{ end }}

{{ template "index" . }}


{{- if gt (len .Functions) 0 }}
## Functions
{{ range .Functions -}}
{{ template "function" . }}
{{ end -}}
{{- end }}

{{ if gt (len .Types) 0 }}
## Types
{{ range .Types -}}


### type <b>{{ .Name }}</b>
{{ if ne .Description "" }}{{ .Description }}{{ end -}}
{{ if gt (len .Fields) 0 }}

###### Properties

| name | type | description |
|------|------|-------------|
{{ range .Fields -}}
| `{{ .Name }}` | `{{ .Type }}` | {{ .Description }} |
{{ end -}}

{{ if gt (len .Examples) 0 }}
###### Examples
{{ range .Examples -}}
{{ .Description }}
```python
{{ .Code }}
```
{{ end -}}
{{ end -}}


{{ end -}}

{{ if gt (len .Methods) 0 }}

###### Methods

{{- range .Methods -}}
{{ template "function" . }}
{{ end -}}
{{- if gt (len .Operators) 0 }}

###### Operators

| operator | description |
|----------|-------------|
{{ range .Operators -}}
	| {{ .Opr }} | {{ .Description }} |
{{ end }}

{{ end }}

{{ end }}
{{- end -}}
{{- end -}}
{{ end }}