module github.com/mcuadros/ascode

go 1.14

require (
	github.com/Masterminds/semver/v3 v3.0.3
	github.com/b5/outline v0.0.0-20190307020728-8cdd78996e40
	github.com/containers/image/v5 v5.2.1
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/go-metrics v0.0.1 // indirect
	github.com/go-git/go-git/v5 v5.0.0
	github.com/gophercloud/utils v0.0.0-20190128072930-fbb6ab446f01
	github.com/gopherjs/gopherjs v0.0.0-20181103185306-d547d1d9531e // indirect
	github.com/hashicorp/go-hclog v0.11.0
	github.com/hashicorp/go-plugin v1.0.1-0.20190610192547-a1bc61569a26
	github.com/hashicorp/hcl2 v0.0.0-20190618163856-0b64543c968c
	github.com/hashicorp/terraform v0.12.23
	github.com/jessevdk/go-flags v1.4.0
	github.com/mitchellh/cli v1.0.0
	github.com/oklog/ulid/v2 v2.0.2
	github.com/qri-io/starlib v0.4.2-0.20200213133954-ff2e8cd5ef8d
	github.com/smartystreets/assertions v0.0.0-20190116191733-b6c0e53d7304 // indirect
	github.com/smartystreets/goconvey v0.0.0-20181108003508-044398e4856c // indirect
	github.com/stretchr/testify v1.5.1
	github.com/vmihailenco/msgpack v4.0.4+incompatible // indirect
	github.com/zclconf/go-cty v1.3.1
	go.starlark.net v0.0.0-20200306205701-8dd3e2ee1dd5
	gopkg.in/yaml.v2 v2.2.8 // indirect
)

replace github.com/hashicorp/hcl2 => github.com/mcuadros/hcl2 v0.0.0-20190711172820-dd3dbf62a554
