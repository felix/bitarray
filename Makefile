
LINTER	= go vet ./...
# Prefer golangci-lint if it exists
ifneq ($(shell command -v golangci-lint),)
	LINTER=golangci-lint run ./...
endif

test: lint
	go test -short -coverprofile=coverage.txt -covermode=atomic ./... \
		&& go tool cover -func=coverage.txt

lint: ; $(LINTER)

clean: ; rm -rf coverage*
