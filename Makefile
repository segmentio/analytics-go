build:
	@godep go build

test:
	@godep go test -race -cover ./...

.PHONY: build test
