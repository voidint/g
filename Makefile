GO = CGO_ENABLED=0 GO111MODULE=on GOPROXY=https://goproxy.cn,direct go
BUILD_DATE := $(shell date '+%Y-%m-%d %H:%M:%S')
GIT_BRANCH := $(shell git symbolic-ref --short -q HEAD)
GIT_COMMIT_HASH := $(shell git rev-parse HEAD|cut -c 1-8)
GO_FLAGS := -v -ldflags="-X 'github.com/voidint/g/build.Built=$(BUILD_DATE)' -X 'github.com/voidint/g/build.GitCommit=$(GIT_COMMIT_HASH)' -X 'github.com/voidint/g/build.GitBranch=$(GIT_BRANCH)'"


all: install lint test clean

build:
	mkdir -p bin
	$(GO) build $(GO_FLAGS) -o ./bin/g

install: build
	$(GO) install $(GO_FLAGS)

build-all: build-linux build-darwin build-windows

build-linux: build-linux-386 build-linux-amd64 build-linux-arm build-linux-arm64 build-linux-s390x
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

install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install github.com/securego/gosec/v2/cmd/gosec@latest
	go install github.com/google/addlicense@latest


lint:
	# Please make sure you are using the latest go version before executing lint
	@go version
	go vet ./...
	golangci-lint run ./...
	staticcheck ./...
	gosec -exclude=G107,G204,G304,G401,G505 -quiet ./...

test:
	go test -v -gcflags=all=-l ./...

test-coverage:
	go test -gcflags=all=-l -race -coverprofile=coverage.txt -covermode=atomic ./...

view-coverage: test-coverage
	go tool cover -html=coverage.txt
	rm -f coverage.txt

addlicense:
	addlicense -v -c "voidint <voidint@126.com>" -l mit -ignore '.github/**' -ignore 'vendor/**' -ignore '**/*.yml' -ignore '**/*.html' -ignore '**/*.sh' .

clean:
	$(GO) clean -x
	rm -f sha256sum.txt
	rm -rf bin
	rm -f coverage.txt

upgrade-deps:
	go get -u -v github.com/urfave/cli/v2@latest
	go get -u -v github.com/Masterminds/semver/v3@latest
	go get -u -v github.com/PuerkitoBio/goquery@latest
	go get -u -v github.com/mholt/archiver/v3@latest
	go get -u -v github.com/schollz/progressbar/v3@latest
	go get -u -v github.com/daviddengcn/go-colortext@latest
	go get -u -v github.com/fatih/color@latest
	go get -u -v github.com/k0kubun/go-ansi@latest
	go get -u -v github.com/agiledragon/gomonkey/v2@latest
	go get -u -v github.com/stretchr/testify@latest
	go get -u -v golang.org/x/text@latest

.PHONY: all build install install-tools lint test test-coverage view-coverage addlicense package clean upgrade-deps build-linux build-darwin build-windows build-linux-386 build-linux-amd64 build-linux-arm build-linux-arm64 build-linux-s390x build-darwin-amd64 build-darwin-arm64 build-windows-386 build-windows-amd64 build-windows-arm build-windows-arm64
