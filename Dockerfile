FROM golang:1.14-alpine 

LABEL MAINTAINER="Máximo Cuadros <mcuadros@gmail.com>"
LABEL "com.github.actions.description"="convert starlark files to HCL"
LABEL "com.github.actions.name"="ascode-action"
LABEL "com.github.actions.color"="blue"

RUN ["/bin/sh", "-c", "apk add --update --no-cache bash ca-certificates curl git jq openssh"]

RUN GO111MODULE=on go get github.com/mcuadros/ascode@08d1962

COPY entrypoint.sh /

ENTRYPOINT ["/entrypoint.sh"]
