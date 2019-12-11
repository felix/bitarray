
.PHONY: test
test: lint ## Run tests and create coverage report
	go test -short -coverprofile=coverage.txt -covermode=atomic ./... \
		&& go tool cover -func=coverage.txt

.PHONY: lint
lint: ## Lint code
	go vet ./... && revive ./...

.PHONY: clean
clean: ## Clean up temp files and binaries
	rm -rf coverage*
