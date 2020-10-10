GO111MODULE := on
export GO111MODULE

clean:
	@go mod tidy

update:
	@go get -u

run-example:
	@make build && rm -rf example/stage/bin && cp -r ./bin example/stage/bin && cd example && docker-compose up --build

build:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o stage/bin/dpm_service .

test:
	@cd ./dpm && go test -v -race

.PHONY: test
