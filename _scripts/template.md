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

{{- range . -}}
---
title: '{{ .Path }}'
---

{{ if ne .Description "" }}{{ .Description }}{{ end }}

{{- if gt (len .Functions) 0 }}
## Functions
{{ range .Functions -}}
{{ template "function" . }}
{{ end -}}
{{- end }}

{{ if gt (len .Types) 0 }}
## Types
{{ range .Types -}}


### <b>{{ .Name }}</b>
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