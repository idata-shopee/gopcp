GOPATH := $(shell cd ../../../.. && pwd)
export GOPATH

test:
	@cd test && go test -cover

.PHONY:	test
