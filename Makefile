.PHONY: build clean test

BIN=g
DIR_SRC=.

# GO_ENV=CGO_ENABLED=0
GO_FLAGS=-ldflags="-X 'build.Build=`date`' -X 'build.Commit=`git rev-parse --verify HEAD`' -X 'build.Branch=`git rev-parse --abbrev-ref HEAD`' -extldflags"
GO=$(shell which go)

build:$(DIR_SRC)/main.go
	@$(GO) build $(GO_FLAGS) -o $(BIN) $(DIR_SRC)

install: build
	@$(GO) install $(GO_FLAGS) $(DIR_SRC)

test:
	@$(GO) test ./...

# clean all build result
clean:
	@$(GO) clean ./...

all: clean build
