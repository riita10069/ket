.ONESHELL:

OS   := $(shell uname | awk '{print tolower($$0)}')
ARCH := $(shell case $$(uname -m) in (x86_64) echo amd64 ;; (aarch64) echo arm64 ;; (*) echo $$(uname -m) ;; esac)


GOLANGCI_LINT_VERSION := 1.42.0

DEV_DIR   := $(shell pwd)/_dev
BIN_DIR   := $(DEV_DIR)/bin
TOOLS_DIR := $(DEV_DIR)/tools
TOOLS_SUM := $(TOOLS_DIR)/go.sum

DELVE         := $(abspath $(BIN_DIR)/dlv)
GOFUMPT       := $(abspath $(BIN_DIR)/gofumpt)
GOLANGCI_LINT := $(abspath $(BIN_DIR)/golangci-lint)

BUILD_TOOLS := cd $(TOOLS_DIR) && go build -o


delve: $(DELVE) ## install delve to _dev/bin
$(DELVE): $(TOOLS_SUM)
	@$(BUILD_TOOLS) $(DELVE) github.com/go-delve/delve/cmd/dlv

gofumpt: $(GOFUMPT) ## install gofumpt to _dev/bin
$(GOFUMPT): $(TOOLS_SUM)
	@$(BUILD_TOOLS) $(GOFUMPT) mvdan.cc/gofumpt

golangci-lint: $(GOLANGCI_LINT) ## install golangci_lint to _dev/bin
$(GOLANGCI_LINT):
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(BIN_DIR) v$(GOLANGCI_LINT_VERSION)

.PHONY: test
test: ## go test for ./...
	@go test -count=1 -race --tags=test -v ./...

.PHONY: vet
vet: ## run go vet for ./...
	@go vet --tags=test ./...

.PHONY: fmt
fmt: $(GOFUMPT) ## run go fmt for ./
	@! $(GOFUMPT) -s -d ./ | grep -E '^'

.PHONY: lint
lint: $(GOLANGCI_LINT) ## run golangci lint for ./...
	@$(GOLANGCI_LINT) run --config ./.golangci.yml ./...

.PHONY: help
help: ## to show how to use the makefile
	@grep -E '^[a-zA-Z0-9/_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
