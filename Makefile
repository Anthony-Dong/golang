.PHONY: ci
ci: check lint build

AppName := devtool

.PHONY: build
build:
	./build.sh build $(AppName)
	./build.sh build protoc-gen-console
	./build.sh build print-env
	CGO_ENABLED=1 IS_SUBMOD=1 ./build.sh build tcpdump_tools

.PHONY: cors
cors:
	./build.sh cors $(AppName)
	./build.sh cors protoc-gen-console

.PHONY: install
install:
	./build.sh install $(AppName)
	./build.sh install protoc-gen-console
	./build.sh install print-env
	CGO_ENABLED=1 IS_SUBMOD=1 ./build.sh install tcpdump_tools

.PHONY: lint
lint:
	bash build.sh format

.PHONY: check
check:
	go run cli/scan_keyword/main.go -dir . -keywords "Ynl0ZWQub3Jn;Ynl0ZWRhbmNlLm5ldA==;Ynl0ZWRhbmNlLm9yZw==;ZmVpc2h1LmNu;bGFya29mZmljZS5jb20="

.PHONY: test
test: ## go tool cover -html=cover.out
	go test -coverprofile cover.out -count=1 ./pkg/...

.PHONY: release_assert
release_assert: build cors ## 创建 release assert
	mv bin/tcpdump_tools bin/$$(go env GOOS)_$$(go env GOARCH)
	zip -j bin/devtool_darwin_amd64.zip bin/darwin_amd64/*
	zip -j bin/devtool_darwin_arm64.zip bin/darwin_arm64/*
	zip -j bin/devtool_linux_amd64.zip bin/linux_amd64/*
	zip -j bin/devtool_windows_amd64.zip bin/windows_amd64/*

.PHONY: help
help: ## help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf " \033[36m%-20s\033[0m  %s\n", $$1, $$2}' $(MAKEFILE_LIST)