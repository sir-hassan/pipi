PKG_LIST := $(shell go list ./...)

all: build

clean:
	go clean
	rm -rf bin

test: clean
	test -z '$(shell gofmt -l .)'
	golint -set_exit_status $(PKG_LIST)
	go vet ./...
	go test ./... -v
	mkdir -p bin
	go test ./... -v --coverprofile bin/coverage.txt
	go tool cover -func bin/coverage.txt

build: test
	mkdir -p bin
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/pipi .
