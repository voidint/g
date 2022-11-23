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
    "armv6l" | "armv7l")
        echo "arm"
	;;
    "s390x")
        echo "s390x"
        ;;
    *)
        echo ${NIL}
        ;;
    esac
}

function get_os() {
    echo $(uname -s | awk '{print tolower($0)}')
}

main() {
    local release="1.4.0"
    local os=$(get_os)
    local arch=$(get_arch)
    local dest_file="${HOME}/.g/downloads/g${release}.${os}-${arch}.tar.gz"
    local url="https://github.com/voidint/g/releases/download/v${release}/g${release}.${os}-${arch}.tar.gz"

    echo "[1/3] Downloading ${url}"
    rm -f "${dest_file}"
    if [ -x "$(command -v wget)" ]; then
        wget -q -P "${HOME}/.g/downloads" "${url}"
    else
        curl -s -S -L --create-dirs -o "${dest_file}" "${url}"
    fi

    echo "[2/3] Install g to the ${HOME}/.g/bin"
    mkdir -p "${HOME}/.g/bin"
    tar -xz -f "${dest_file}" -C "${HOME}/.g/bin"
    chmod +x "${HOME}/.g/bin/g"

    echo "[3/3] Set environment variables"
    cat >${HOME}/.g/env <<-'EOF'
#!/bin/sh
# g shell setup
export GOROOT="${HOME}/.g/go"
export PATH="${HOME}/.g/bin:${HOME}/.g/go/bin:$PATH"
export G_MIRROR=https://golang.google.cn/dl/
	EOF


    if [ -x "$(command -v bash)" ]; then
        cat >>${HOME}/.bashrc <<-'EOF'
# g shell setup
if [ -f "${HOME}/.g/env" ]; then
    . "${HOME}/.g/env"
fi
		EOF
    fi

    if [ -x "$(command -v zsh)" ]; then
        cat >>${HOME}/.zshrc <<-'EOF'
# g shell setup
if [ -f "${HOME}/.g/env" ]; then
    . "${HOME}/.g/env"
fi
		EOF
    fi


    echo -e "\nTo configure your current shell, run:\nsource \"$HOME/.g/env\""

    exit 0
}

main
