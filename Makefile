SOURCE_FILES?=./...
TEST_PATTERN?=.
TEST_OPTIONS?=

export GOPROXY 		:= https://proxy.golang.org,https://gocenter.io,direct
export PATH 		:= ./bin:$(PATH)
export GO111MODULE 	:= on

# Install all the build and lint dependencies
setup:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh
	go mod tidy
.PHONY: setup

test:
	go test $(TEST_OPTIONS) -v -failfast -race -coverpkg=./... -covermode=atomic -coverprofile=coverage.out $(SOURCE_FILES) -run $(TEST_PATTERN) -timeout=2m

cover: test
	go tool cover -html=coverage.out

fmt:
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do gofmt -w -s "$$file"; goimports -w "$$file"; done

lint:
	./bin/golangci-lint run ./...

ci: lint test

build:
	go build -o fork-cleaner ./cmd/fork-cleaner/main.go

.DEFAULT_GOAL := build
