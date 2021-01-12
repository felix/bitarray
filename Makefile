
test: lint
	go test -short -coverprofile=coverage.txt -covermode=atomic ./... \
		&& go tool cover -func=coverage.txt

lint: ; go vet ./...

clean: ; rm -rf coverage*
