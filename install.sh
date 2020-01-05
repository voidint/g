#!/usr/bin/env bash
set -e

get_arch() {
    a=$(uname -m)
    case ${a} in
    "x86_64" | "amd64" )
        echo "amd64"
        ;;
    "i386" | "i486" | "i586")
        echo "386"
        ;;
    *)
        echo ${NIL}
        ;;
    esac
}

get_os(){
    echo $(uname -s | awk '{print tolower($0)}')
}

main() {
    local os=$(get_os)
    local arch=$(get_arch)
    local dest_file="${HOME}/g1.1.1.${os}-${arch}.tar.gz"
    local url="https://github.com/voidint/g/releases/download/v1.1.1/g1.1.1.${os}-${arch}.tar.gz"

    echo "[1/3] Download ${url}"
    rm -f ${dest_file}
    wget -q -P "${HOME}" "${url}" 
    chmod +x ${dest_file}

    echo "[2/3] Install g to the ${HOME}/bin"
    mkdir -p "${HOME}/bin"
    tar -xz -f ${dest_file} -C ${HOME}/bin
    chmod +x ${HOME}/bin/g

    echo "[3/3] Set environment variables"
    echo '# ===== set g environment variables =====' >> ${HOME}/.bashrc
    echo 'export GOROOT="${HOME}/.g/go"' >> ${HOME}/.bashrc
    echo 'export PATH="${HOME}/bin:${HOME}/.g/go/bin:$PATH"' >> ${HOME}/.bashrc
    echo 'export G_MIRROR=https://golang.google.cn/dl/' >> ${HOME}/.bashrc 

    echo '# ===== set g environment variables =====' >> ${HOME}/.zshrc
    echo 'export GOROOT="${HOME}/.g/go"' >> ${HOME}/.zshrc
    echo 'export PATH="${HOME}/bin:${HOME}/.g/go/bin:$PATH"' >> ${HOME}/.zshrc
    echo 'export G_MIRROR=https://golang.google.cn/dl/' >> ${HOME}/.zshrc

    exit 0
}

main