# Pinned tool versions. CI uses the same values so `make lint` reproduces it.
GOLANGCI_LINT_VERSION ?= v2.12.2
GOLANGCI_LINT         := go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
GOSEC                 := go run github.com/securego/gosec/v2/cmd/gosec@latest
GOVULNCHECK           := go run golang.org/x/vuln/cmd/govulncheck@latest

.DEFAULT_GOAL := build

.PHONY: build
build: ## Compile all packages and the node and bootstrap binaries.
	go build -o bin/bartering .
	go build -o bin/bootstrap ./bootstrap-node

.PHONY: test
test: ## Run the unit tests.
	go test ./...

.PHONY: test-race
test-race: ## Run the unit tests with the race detector.
	go test -race ./...

.PHONY: cover
cover: ## Run the tests and write a coverage profile to coverage.out.
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out | tail -1

.PHONY: fmt
fmt: ## Format the source tree.
	gofmt -w .

.PHONY: vet
vet: ## Run go vet.
	go vet ./...

.PHONY: lint
lint: ## Run golangci-lint (pinned version).
	$(GOLANGCI_LINT) run ./...

.PHONY: sec
sec: ## Run the gosec static security scanner.
	$(GOSEC) -quiet ./...

.PHONY: vuln
vuln: ## Check dependencies and the standard library for known vulnerabilities.
	$(GOVULNCHECK) ./...

.PHONY: tidy
tidy: ## Tidy and verify the module files.
	go mod tidy
	go mod verify

.PHONY: check
check: fmt vet lint test sec vuln ## Run the full local verification suite.

.PHONY: clean
clean: ## Remove build output.
	rm -rf bin coverage.out

.PHONY: help
help: ## Print this help.
	@grep -hE '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'
