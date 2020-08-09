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
    printf "============Pakcage for %s============\n" $3
    local rootdir=${1}
    local release=${2}
    local osarch=(${3//_/ })
    local os=${osarch[0]}
    local arch=${osarch[1]}

    printf "[1/4] Cross compile@%s_%s\n" ${os} ${arch}
    GOOS=${os} GOARCH=${arch} gbb

    local bindir=""
    if [ $(get_os) = ${os} ] && [ $(get_arch) = ${arch} ]; then
        bindir=$GOPATH/bin
    else
        bindir=$GOPATH/bin/${os}_${arch}
    fi
    printf "[2/4] Change directory to %s\n" ${bindir}
    cd ${bindir}

    printf "[3/4] Package\n"
    if [ ${os} == "windows" ]; then
        zip $rootdir/g${release}.${os}-${arch}.zip ./g.exe
    else
        tar -czv -f $rootdir/g${release}.${os}-${arch}.tar.gz ./g
    fi

    printf "[4/4] Change directory to %s\n\n" ${rootdir}
    cd ${rootdir}
}

main() {
    export CGO_ENABLED="0"
    export GO111MODULE="on"
    export GOPROXY="https://goproxy.cn,direct"

    local release="1.2.0"
    local rootdir="$(pwd)"

    for item in "darwin_amd64" "linux_386" "linux_amd64" "linux_arm" "linux_arm64" "windows_386" "windows_amd64"; do
        package ${rootdir} ${release} ${item}
    done
}

main
