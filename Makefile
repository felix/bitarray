# If there is a GOPATH use it, otherwise we define our own. This is not exported
# and should not impact the Go modules system.
GOPATH?=	$(shell go env GOPATH)
SRC=		$(shell find . -type f -name '*.go')

export GO111MODULE=on
export CGO_ENABLED=0

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		|awk 'BEGIN{FS=":.*?## "};{printf "\033[36m%-30s\033[0m %s\n",$$1,$$2}'

.PHONY: test
test: lint ## Run tests and create coverage report
	go test -short -coverprofile=coverage.txt -covermode=atomic ./...
	go tool cover -html=coverage.txt -o coverage.html

lint: ## Lint code
	golint -set_exit_status ./...

.PHONY: golint
$(GOPATH)/bin/golint:
	@echo Building golint...
	go get -u golang.org/x/lint/golint

.PHONY: clean
clean: ## Clean up temp files and binaries
	rm -rf coverage*
