all: get_vendor_deps install test

get_vendor_deps:
	@rm -rf vendor/
	@echo "--> Running dep ensure"
	@dep ensure -v
	
install:
	go install ./cmd/iris

test:
	@go test `glide novendor`
test_cli:
	bash ./cmd/iris/sh_tests/stake.sh

build_linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/iris ./cmd/iris && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/iriscli ./cmd/iriscli

build:
	go build -o build/iris ./cmd/iris && \
    go build -o build/iris_cli ./cmd/iriscli