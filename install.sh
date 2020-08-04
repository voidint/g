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

main() {
    local release="1.1.3"
    local os=$(get_os)
    local arch=$(get_arch)
    local dest_file="${HOME}/g${release}.${os}-${arch}.tar.gz"
    local url="https://github.com/voidint/g/releases/download/v${release}/g${release}.${os}-${arch}.tar.gz"

    echo "[1/3] Downloading ${url}"
    rm -f "${dest_file}"
    if [ -x "$(command -v wget)" ]; then
        wget -q -P "${HOME}" "${url}"
    else
        curl -s -S -L -o "${dest_file}" "${url}"
    fi

    echo "[2/3] Install g to the ${HOME}/bin"
    mkdir -p "${HOME}/bin"
    tar -xz -f "${dest_file}" -C "${HOME}/bin"
    chmod +x "${HOME}/bin/g"

    echo "[3/3] Set environment variables"
    if [ -x "$(command -v bash)" ]; then
        cat >>${HOME}/.bashrc <<-'EOF'
		# ===== set g environment variables =====
		export GOROOT="${HOME}/.g/go"
		export PATH="${HOME}/bin:${HOME}/.g/go/bin:$PATH"
		export G_MIRROR=https://golang.google.cn/dl/
		EOF
    fi

    if [ -x "$(command -v zsh)" ]; then
        cat >>${HOME}/.zshrc <<-'EOF'
		# ===== set g environment variables =====
		export GOROOT="${HOME}/.g/go"
		export PATH="${HOME}/bin:${HOME}/.g/go/bin:$PATH"
		export G_MIRROR=https://golang.google.cn/dl/
		EOF
    fi

    exit 0
}

main
