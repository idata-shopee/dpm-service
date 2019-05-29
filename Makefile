GOPATH := $(shell cd ../../../.. && pwd)
export GOPATH

init-dep:
	@dep init

dep:
	@dep ensure

status-dep:
	@dep status

update-dep:
	@dep ensure -update

run-example:
	@make build && rm -rf example/stage/bin && cp -r ./bin example/stage/bin && cd example && docker-compose up --build

build:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o stage/bin/dpm_service .

test:
	@echo "test"

.PHONY: test
