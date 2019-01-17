ARTIFACTS_DIR ?= tmp

get:
	@go get -v -t ./...

vet:
	@go vet ./...

build:
	@go build ./...

test:
	@mkdir -p $(ARTIFACTS_DIR)
	@go test -race -coverprofile=$(ARTIFACTS_DIR)/cover.out .
	@go tool cover -func $(ARTIFACTS_DIR)/cover.out -o $(ARTIFACTS_DIR)/cover.txt
	@go tool cover -html $(ARTIFACTS_DIR)/cover.out -o $(ARTIFACTS_DIR)/cover.html

RUN_E2E_TESTS ?= false

ci: get vet test
	@if [ "$(RUN_E2E_TESTS)" != "true" ]; then \
	  echo "Skipping end to end tests."; else \
		go get github.com/segmentio/library-e2e-tester/cmd/tester; \
		tester -segment-write-key=$(SEGMENT_WRITE_KEY) -runscope-token=$(RUNSCOPE_TOKEN) -runscope-bucket=$(RUNSCOPE_BUCKET) -path='cli'; fi

.PHONY: get vet build test ci
