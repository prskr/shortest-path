VERSION = $(shell git describe --dirty --tags --always)
REPO = github.com/baez90/shortest-path
BUILD_PATH = $(REPO)/cmd/shortest-path
PKGS = $(shell go list ./...)
TEST_PKGS = $(shell find . -type f -name "*_test.go" -printf '%h\n' | sort -u)
GOARGS = GOOS=linux GOARCH=amd64
GO_BUILD_ARGS = -ldflags="-w -s"
BINARY_NAME = shortest-path
DIR = $(dir $(realpath $(firstword $(MAKEFILE_LIST))))
DEBUG_PORT = 2345

.PHONY: all clean

all: format compile

rebuild: clean format compile

format:
	@go fmt $(PKGS)

revive: ensure-revive
	@revive --config $(DIR)assets/lint/config.toml -formatter friendly $(DIR)...

clean:
	@rm -f debug $(BINARY_NAME)
	@rm -rf dist

test:
	@go test -coverprofile=./cov-raw.out -v $(TEST_PKGS)
	@cat ./cov-raw.out | grep -v "generated" > ./cov.out

cli-cover-report:
	@go tool cover -func=cov.out

html-cover-report:
	@go tool cover -html=cov.out -o .coverage.html

deps:
	@go build -v ./...

compile: deps
	@$(GOARGS) go build $(GO_BUILD_ARGS) -o $(DIR)/$(BINARY_NAME) $(BUILD_PATH)

watch-test: ensure-reflex
	@reflex -r '\.go$$' -s -- sh -c 'make test'

serve-godoc: ensure-godoc
	@godoc -http=:6060

serve-docs: ensure-reflex docs
	@reflex -r '\.md$$' -s -- sh -c 'mdbook serve -d $(DIR)/public -n 127.0.0.1 $(DIR)/docs'

docs:
	@mdbook build -d $(DIR)/public $(DIR)/docs`

test-release: ensure-goreleaser ensure-packr2
	@goreleaser --snapshot --skip-publish --rm-dist

ensure-revive:
ifeq (, $(shell which revive))
	$(shell go get -u github.com/mgechev/revive)
endif

ensure-delve:
ifeq (, $(shell which dlv))
	$(shell go get -u github.com/go-delve/delve/cmd/dlv)
endif

ensure-reflex:
ifeq (, $(shell which reflex))
	$(shell go get -u github.com/cespare/reflex)
endif

ensure-godoc:
ifeq (, $(shell which godoc))
	$(shell go get -u golang.org/x/tools/cmd/godoc)
endif

ensure-goreleaser:
ifeq (, $(shell which goreleaser))
	$(shell curl -sL https://github.com/goreleaser/goreleaser/releases/download/v$(GORELEASER_VERSION)/goreleaser_Linux_x86_64.tar.gz | tar -xvz --exclude "*.md" -C $$GOPATH/bin)
endif