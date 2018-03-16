ifndef CIRCLE_ARTIFACTS
CIRCLE_ARTIFACTS=tmp
endif

get:
	@go get -v -t ./...

vet:
	@go vet ./...

build:
	@go build ./...

test:
	@mkdir -p ${CIRCLE_ARTIFACTS}
	@go test -race -coverprofile=${CIRCLE_ARTIFACTS}/cover.out .
	@go tool cover -func ${CIRCLE_ARTIFACTS}/cover.out -o ${CIRCLE_ARTIFACTS}/cover.txt
	@go tool cover -html ${CIRCLE_ARTIFACTS}/cover.out -o ${CIRCLE_ARTIFACTS}/cover.html

ci: get vet test
	@if [ "$(RUN_E2E_TESTS)" != "true" ]; then \
	  echo "Skipping end to end tests."; else \
		@wget https://github.com/segmentio/library-e2e-tester/releases/download/0.1.1/tester_linux_amd64; \
		@chmod +x tester_linux_amd64; \
		@go build ./cmd/...; \
		@tester_linux_amd64 -segment-write-key=$(SEGMENT_WRITE_KEY) -runscope-token=$(RUNSCOPE_TOKEN) -runscope-bucket=$(RUNSCOPE_BUCKET) -path='cli'; fi

.PHONY: get vet build test ci
