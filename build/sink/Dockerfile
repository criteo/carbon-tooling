ARG GO_VERSION=1.12
FROM golang:$GO_VERSION-alpine as builder
RUN apk add --update make
WORKDIR /go/src/github.com/criteo/carbon-tooling
COPY Makefile .
COPY vendor/ vendor/
COPY sink/ sink/
RUN make build-sink 

FROM alpine
WORKDIR /app
COPY --from=builder go/src/github.com/criteo/carbon-tooling/cmd/sink /app/sink
ENTRYPOINT ./sink
