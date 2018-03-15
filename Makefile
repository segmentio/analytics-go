ifndef CIRCLE_ARTIFACTS
CIRCLE_ARTIFACTS=tmp
endif

get:
	@go get -v -t ./...

vet:
	@go vet ./...

build:
	@go build

test:
	@mkdir -p ${CIRCLE_ARTIFACTS}
	@go test -race -coverprofile=${CIRCLE_ARTIFACTS}/cover.out .
	@go tool cover -func ${CIRCLE_ARTIFACTS}/cover.out -o ${CIRCLE_ARTIFACTS}/cover.txt
	@go tool cover -html ${CIRCLE_ARTIFACTS}/cover.out -o ${CIRCLE_ARTIFACTS}/cover.html

ci: get vet test

.PHONY: get vet build test ci
