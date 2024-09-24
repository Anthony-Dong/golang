.PHONY: init kitex

MAKEFILE_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

THRIFTGO := $(MAKEFILE_DIR)bin/thriftgo

KITEX := $(MAKEFILE_DIR)bin/kitex

init: export GOBIN = $(MAKEFILE_DIR)bin
init: ## https://www.cloudwego.io/docs/kitex/getting-started/prerequisite/
	go install github.com/cloudwego/thriftgo@latest
	$(THRIFTGO) --version
	go install github.com/cloudwego/kitex/tool/cmd/kitex@latest
	$(KITEX) --version

kitex:
	#$(KITEX) -service a.b.c ../../idl/test/api.thrift
	$(KITEX) ../../idl/test/api.thrift
	rm -rf handler.go

run:
	bash -e build.sh
	bash -e output/bootstrap.sh