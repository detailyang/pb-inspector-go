# Background color
GREEN                := $(shell tput -Txterm setaf 2)
YELLOW               := $(shell tput -Txterm setaf 3)
BLUE                 := $(shell tput -Txterm setaf 4)
MAGENTA              := $(shell tput -Txterm setaf 5)
WHITE                := $(shell tput -Txterm setaf 7)
RESET                := $(shell tput -Txterm sgr0)
TARGET_MAX_CHAR_NUM  := 20


## Show help
help:
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET} ${MAGENTA}[variable=value]${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  ${YELLOW}%-$(TARGET_MAX_CHAR_NUM)s${RESET} ${GREEN}%s${RESET}\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.PHONY: format
## Format *.go by go format
format:
	go fmt ./...

.PHONY: format-pb
## Format *.proto
format-pb:
	prototool format -w ./proto

.PHONY: lint
## Lint *.go via golangci-lint
lint:
	golangci-lint run -v

## Generate go
generate:
	rm -rf ./proto/go && prototool generate

## Lint *.proto
lint-pb:
	prototool lint

.PHONY: build
## Build binary
build:
	go build -v -o bin/pb-inspector cmd/pb-inspector/*.go

.PHONY: test
## Testing everything
test:
	go test -failfast -race -v ./...

.PHONY: clean
## Clean artifact
clean:
	rm -rf bin/* *.test *.out
