.PHONY: ci
ci: check lint build

.PHONY: build
build:
	./build.sh build devtool
	./build.sh build protoc-gen-console
	./build.sh build print-env

.PHONY: install
install:
	./build.sh install devtool
	./build.sh install protoc-gen-console
	./build.sh install print-env

.PHONY: lint
lint:
	bash build.sh format

.PHONY: check
check:
	go run cli/scan_keyword/main.go -dir . -keywords "Ynl0ZWQub3Jn;Ynl0ZWRhbmNlLm5ldA==;Ynl0ZWRhbmNlLm9yZw==;ZmVpc2h1LmNu;bGFya29mZmljZS5jb20="

.PHONY: test
test: ## go tool cover -html=cover.out
	CGO_ENABLED=1 go test -coverprofile cover.out -count=1 ./pkg/...

.PHONY: help
help: ## help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf " \033[36m%-20s\033[0m  %s\n", $$1, $$2}' $(MAKEFILE_LIST)