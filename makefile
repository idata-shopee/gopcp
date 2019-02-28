GOPATH := $(shell cd ../../../.. && pwd)
export GOPATH

test:
	@go test -v -race

cover:
	@go test -coverprofile=coverage.out
	@go tool cover -html=coverage.out

.PHONY:	test
