#!/usr/bin/env bash
set -e

function get_arch() {
    a=$(uname -m)
    case ${a} in
    "x86_64" | "amd64")
        echo "amd64"
        ;;
    "i386" | "i486" | "i586")
        echo "386"
        ;;
    "aarch64" | "arm64")
        echo "arm64"
        ;;
    *)
        echo ${NIL}
        ;;
    esac
}

function get_os() {
    echo $(uname -s | awk '{print tolower($0)}')
}

function package() {
    printf "============Pakcage for %s============\n" $2
    local release=${1}
    local osarch=(${2//_/ })
    local os=${osarch[0]}
    local arch=${osarch[1]}

    printf "[1/2] Cross compile@%s_%s\n" ${os} ${arch}
    GOOS=${os} GOARCH=${arch} gbb

    printf "[2/2] Package\n"
    if [ ${os} == "windows" ]; then
        zip g${release}.${os}-${arch}.zip ./g.exe
    else
        tar -czv -f g${release}.${os}-${arch}.tar.gz ./g
    fi
}

main() {
    export CGO_ENABLED="0"
    export GO111MODULE="on"
    export GOPROXY="https://goproxy.cn,direct"

    local release="1.2.1"

    for item in "darwin_amd64" "darwin_arm64" "linux_386" "linux_amd64" "linux_arm" "linux_arm64" "windows_386" "windows_amd64" "windows_arm" "windows_arm64"; do
        package ${release} ${item}
    done
    go clean
}

main
