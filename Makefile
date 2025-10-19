SHELL := /usr/bin/env bash -o pipefail
GO_CMD := go
BINARY := goarapi
BUILD_DIR := bin
ENTRY := ./cmd/goarapi

.PHONY: help
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*##' Makefile | sed 's/:.*##/: /'

.PHONY: tidy
tidy: ## Download and tidy dependencies
	$(GO_CMD) mod tidy

.PHONY: test
test: ## Run unit tests
	$(GO_CMD) test ./...

.PHONY: build
build: tidy ## Build CLI binary
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 $(GO_CMD) build -o $(BUILD_DIR)/$(BINARY) $(ENTRY)

.PHONY: run
run: ## Run CLI locally
	$(GO_CMD) run $(ENTRY)

.PHONY: docker-build
docker-build: ## Build docker image
	docker build -t $(BINARY):latest .

.PHONY: docker-run
docker-run: ## Run docker container
	docker run --rm -it -v $$PWD/config_example.yaml:/etc/gokeenapi/config.yaml -p 8080:8080 $(BINARY):latest serve --addr :8080 --config /etc/gokeenapi/config.yaml
