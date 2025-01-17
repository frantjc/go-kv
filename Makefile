GO ?= go
GOLANGCI-LINT ?= golangci-lint

all: fmt lint

fmt generate:
	@$(GO) $@ ./...

lint:
	@$(GOLANGCI-LINT) run --fix

gen: generate

.PHONY: all fmt generate lint proto gen
