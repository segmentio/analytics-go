vet:
	@go vet ./...

build:
	@go build

test:
	@go test -race -cover ./...

.PHONY: vet build test
