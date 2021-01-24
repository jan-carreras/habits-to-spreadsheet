
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: build
build: test ## Build app
	go build -o bin/hsync cmd/import/*.go

.PHONY: test
test: ## Run the tests with coverage
	@go test -cover ./...

.PHONY: run
run: ## Execute the cli APP
	@go run cmd/import/*.go

.PHONY: fmt
fmt: ## Reformat code and imports
	@gofmt -l -w $(SRC)
	@goimports -w $(SRC)

.PHONY: check
check: test ## Run linters & gofmt check
	@test -z $(shell gofmt -l $(SRC) | tee /dev/stderr) || (echo "[ERR] Fix formatting issues with 'make fmt'" && false)
	@which golangci-lint > /dev/null 2>/dev/null || (echo "ERROR: golangci-lint not found" && false)
	@golangci-lint run

.PHONY: help
help: ## This help message
	@echo 'usage: make [target] ...'
	@echo 'targets:'
	@egrep '^(.+)\:\ ##\ (.+)' ${MAKEFILE_LIST} | column -t -c 2 -s ':#'