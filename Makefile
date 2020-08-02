.PHONY: build clean test

BIN=g
DIR_SRC=.

# GO_ENV=CGO_ENABLED=0
GO_FLAGS=-ldflags="-X 'github.com/voidint/g/build.Build=`date '+%Y-%m-%d %H:%M:%S'`' -X 'github.com/voidint/g/build.Commit=`git rev-parse --verify HEAD`' -X 'github.com/voidint/g/build.Branch=`git symbolic-ref --short -q HEAD`' -extldflags"
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
