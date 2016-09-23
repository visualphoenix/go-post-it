FROM golang:alpine

RUN apk add --no-cache ca-certificates file openssl

ENV CGO_ENABLED 0 BUILD_FLAGS="-v -ldflags '-d -s -w'"

ARG user
ARG app

COPY *.go /go/src/github.com/$user/$app/
WORKDIR /go/src/github.com/$user/$app

RUN set -x \
 && eval "GOARCH=amd64 go build $BUILD_FLAGS -o /go/bin/$app-amd64" \
 && file /go/bin/$app-amd64

RUN file /go/bin/$app-*
