GOPATH := $(shell cd ../../../.. && pwd)
export GOPATH

test:
	@go test

cover:
	@go test -coverprofile=coverage.out
	@go tool cover -html=coverage.out

.PHONY:	test
