GO = CGO_ENABLED=0 GO111MODULE=on GOPROXY=https://goproxy.cn,direct go
BUILD_DATE := $(shell date '+%Y-%m-%d %H:%M:%S')
GIT_BRANCH := $(shell git symbolic-ref --short -q HEAD)
GIT_COMMIT_HASH := $(shell git rev-parse --verify HEAD)
GO_FLAGS := -v -ldflags="-X 'github.com/voidint/g/build.Build=$(BUILD_DATE)' -X 'github.com/voidint/g/build.Commit=$(GIT_COMMIT_HASH)' -X 'github.com/voidint/g/build.Branch=$(GIT_BRANCH)'"


all: install test clean

build:
	$(GO) build $(GO_FLAGS)

install: build
	$(GO) install $(GO_FLAGS)

build-all: build-linux build-darwin build-windows

build-linux: build-linux-386 build-linux-amd64 build-linux-arm build-linux-arm64
build-linux-386:
	GOOS=linux GOARCH=386 $(GO) build $(GO_FLAGS) -o bin/linux-386/g
build-linux-amd64:
	GOOS=linux GOARCH=amd64 $(GO) build $(GO_FLAGS) -o bin/linux-amd64/g
build-linux-arm:
	GOOS=linux GOARCH=arm $(GO) build $(GO_FLAGS) -o bin/linux-arm/g
build-linux-arm64:
	GOOS=linux GOARCH=arm64 $(GO) build $(GO_FLAGS) -o bin/linux-arm64/g
build-linux-s390x:
        GOOS=linux GOARCH=s390x $(GO) build $(GO_FLAGS) -o  bin/linux-s390x/g


build-darwin: build-darwin-amd64 build-darwin-arm64
build-darwin-amd64:
	GOOS=darwin GOARCH=amd64 $(GO) build $(GO_FLAGS) -o bin/darwin-amd64/g
build-darwin-arm64:
	GOOS=darwin GOARCH=arm64 $(GO) build $(GO_FLAGS) -o bin/darwin-arm64/g


build-windows: build-windows-386 build-windows-amd64 build-windows-arm build-windows-arm64
build-windows-386:
	GOOS=windows GOARCH=386 $(GO) build $(GO_FLAGS) -o bin/windows-386/g.exe
build-windows-amd64:
	GOOS=windows GOARCH=amd64 $(GO) build $(GO_FLAGS) -o bin/windows-amd64/g.exe
build-windows-arm:
	GOOS=windows GOARCH=arm $(GO) build $(GO_FLAGS) -o bin/windows-arm/g.exe
build-windows-arm64:
	GOOS=windows GOARCH=arm64 $(GO) build $(GO_FLAGS) -o bin/windows-arm64/g.exe

package:
	sh ./package.sh

test:
	$(GO) test -v ./...

clean:
	$(GO) clean -x
	rm -f sha256sum.txt
	rm -rf bin

.PHONY: all build install test package clean build-linux build-darwin build-windows build-linux-386 build-linux-amd64 build-linux-arm build-linux-arm64 build-darwin-amd64 build-darwin-arm64 build-windows-386 build-windows-amd64 build-windows-arm build-windows-arm64
