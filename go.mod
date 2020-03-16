module github.com/mcuadros/ascode

go 1.14

require (
	github.com/b5/outline v0.0.0-20190307020728-8cdd78996e40
	github.com/golang/protobuf v1.3.2 // indirect
	github.com/hashicorp/go-hclog v0.9.2
	github.com/hashicorp/go-plugin v1.0.1
	github.com/hashicorp/hcl2 v0.0.0-20190618163856-0b64543c968c
	github.com/hashicorp/terraform v0.12.23
	github.com/jessevdk/go-flags v1.4.0
	github.com/kr/pty v1.1.8 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/mcuadros/hcl2 v0.0.0-20190712010647-ace444864d00
	github.com/mitchellh/cli v1.0.0
	github.com/oklog/ulid v2.0.0+incompatible
	github.com/qri-io/starlib v0.4.2-0.20190710173850-cb41fc97dda5
	github.com/spf13/pflag v1.0.3 // indirect
	github.com/stretchr/objx v0.2.0 // indirect
	github.com/stretchr/testify v1.3.0
	github.com/vmihailenco/msgpack v4.0.4+incompatible // indirect
	github.com/zclconf/go-cty v1.2.1
	go.starlark.net v0.0.0-20200306205701-8dd3e2ee1dd5
	golang.org/x/crypto v0.0.0-20190701094942-4def268fd1a4 // indirect
	google.golang.org/appengine v1.6.1 // indirect
	gopkg.in/src-d/go-git.v4 v4.12.0
)

replace github.com/golangci/golangci-lint => github.com/golangci/golangci-lint v1.18.0

replace github.com/go-critic/go-critic v0.0.0-20181204210945-ee9bf5809ead => github.com/go-critic/go-critic v0.3.5-0.20190526074819-1df300866540

replace github.com/hashicorp/hcl2 => github.com/mcuadros/hcl2 v0.0.0-20190711172820-dd3dbf62a554

replace github.com/Unknwon/com => github.com/unknwon/com v1.0.1
