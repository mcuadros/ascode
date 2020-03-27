FROM golang:1.14.0-alpine AS builder
RUN apk --no-cache add make
WORKDIR $GOPATH/src/github.com/mcuadros/ascode
COPY . .
RUN GO111MODULE=on CGO_ENABLED=0 GOPROXY="https://proxy.golang.org" \
    go build \
        -ldflags "$(make goldflags)" \
        -o /bin/ascode .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
COPY --from=builder /bin/ascode /bin/ascode
CMD ["ascode"]