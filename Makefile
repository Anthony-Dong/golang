.PHONY: ci
ci: check lint build

AppName := devtool

.PHONY: build
build:
	./build.sh build $(AppName)

.PHONY: cors
cors:
	./build.sh cors $(AppName)

.PHONY: install
install:
	./build.sh install $(AppName)
	./build.sh install protoc-gen-console
	./build.sh install print-env
	CGO_ENABLED=1 IS_SUBMOD=1 bash ./build.sh install tcpdump_tools

.PHONY: lint
lint:
	bash build.sh format

.PHONY: check
check:
	go run cli/scan_keyword/main.go -dir . -keywords "Ynl0ZWQub3Jn;Ynl0ZWRhbmNlLm5ldA==;Ynl0ZWRhbmNlLm9yZw==;ZmVpc2h1LmNu"

.PHONY: test
test: ## go tool cover -html=cover.out
	go test -coverprofile cover.out -count=1 ./pkg/...

.PHONY: help
help: ## help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf " \033[36m%-20s\033[0m  %s\n", $$1, $$2}' $(MAKEFILE_LIST)