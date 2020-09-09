.PHONY: build clean test

GO111MODULE=on
GOPROXY=https://goproxy.cn,direct

BIN=g
DIR_SRC=.

# GO_ENV=CGO_ENABLED=0
GO_FLAGS=-ldflags="-X 'github.com/voidint/g/build.Build=`date '+%Y-%m-%d %H:%M:%S'`' -X 'github.com/voidint/g/build.Commit=`git rev-parse --verify HEAD`' -X 'github.com/voidint/g/build.Branch=`git symbolic-ref --short -q HEAD`' -extldflags"

build:$(DIR_SRC)/main.go
	@echo "GO111MODULE=$(GO111MODULE)"
	@echo "GOPROXY=$(GOPROXY)"
	go build $(GO_FLAGS) -o $(BIN) $(DIR_SRC)

install: build
	go install $(GO_FLAGS) $(DIR_SRC)

test:
	go test ./...

clean:
	go clean -x ./...

all: clean build
