.PHONY: ci
ci: check lint build

AppName := devtool

.PHONY: build
build:
	bash build.sh build $(AppName)

.PHONY: install
install:
	bash build.sh install $(AppName)
	$$(go env GOPATH)/bin/$(AppName) --version

.PHONY: lint
lint:
	bash build.sh format

.PHONY: check
check:
	go run cli/scan_keyword/main.go -dir . -keywords "Ynl0ZWQub3Jn;Ynl0ZWRhbmNlLm5ldA==;Ynl0ZWRhbmNlLm9yZw==;ZmVpc2h1LmNu"

.PHONY: test
test: ## go tool cover -html=cover.out
	CGO_ENABLED=1 go test -coverprofile cover.out -count=1 ./pkg/...
	make -C command/tcpdump/test un_compress
	CGO_ENABLED=1 go test -coverprofile cover.out -count=1 ./command/...

.PHONY: help
help: ## help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf " \033[36m%-20s\033[0m  %s\n", $$1, $$2}' $(MAKEFILE_LIST)