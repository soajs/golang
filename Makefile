.PHONY: lint test check

lint:
	@golangci-lint run --config .golangci.yml

test:
	@go test -cover ./...

check: lint test
