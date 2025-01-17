GO ?= go
GOLANGCI-LINT ?= golangci-lint

all: fmt lint

fmt generate test:
	@$(GO) $@ ./...

lint:
	@$(GOLANGCI-LINT) run --fix

gen: generate

.PHONY: all fmt generate test lint proto gen
