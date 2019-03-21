ARG GO_VERSION=1.12
FROM golang:$GO_VERSION-alpine as builder
RUN apk add make
WORKDIR /go/src/github.com/criteo/carbon-tooling
COPY Makefile .
COPY vendor/ vendor/
COPY injector/ injector/
RUN make build-injector

FROM alpine
WORKDIR /app
ENV SLEEP_TIME="100ms"
ENV ENDPOINT_HOST="carbon-relay-ng"
ENV ENDPOINT_PORT=2003
COPY --from=builder /go/src/github.com/criteo/carbon-tooling/cmd/injector /app/injector
ENTRYPOINT ./injector -sendPause=$SLEEP_TIME -host=$ENDPOINT_HOST -port=$ENDPOINT_PORT
