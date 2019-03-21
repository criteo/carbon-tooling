pkgs = ./...

build-injector:
	go build -o cmd/injector injector/injector.go
build-docker-image-injector:
	docker build . -t injector -f build/injector/Dockerfile
 
build-sink:
	go build -o cmd/sink sink/sink.go
build-docker-image-sink:
	docker build . -t sink -f build/sink/Dockerfile 

build-docker-images: build-docker-image-injector build-docker-image-sink

all: format build

build: build-injector build-sink build-docker-images

style: format
	go get github.com/golang/lint/golint
	golint -set_exit_status $(shell go list $(pkgs))

format:
	go fmt $(pkgs)

vet:
	go vet $(pkgs)
