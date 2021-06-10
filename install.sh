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
    local release="1.2.0"
    local install_dir="${HOME}/g"
    if [ -n "$G_HOME" ]; then
      export install_dir="${G_HOME}"
    fi
    local os=$(get_os)
    local arch=$(get_arch)
    local dest_file="${HOME}/g${release}.${os}-${arch}.tar.gz"
    local url="${GIT_MIRROR}https://github.com/voidint/g/releases/download/v${release}/g${release}.${os}-${arch}.tar.gz"

    echo "[1/3] Downloading ${url}"
    rm -f "${dest_file}"
    if [ -x "$(command -v wget)" ]; then
        wget -q -P "${HOME}" "${url}"
    else
        curl -s -S -L -o "${dest_file}" "${url}"
    fi

    echo "[2/3] Install g to the ${install_dir}"
    mkdir -p "${install_dir}"
    tar -xz -f "${dest_file}" -C "${install_dir}"
    chmod +x "${install_dir}/g"

    echo "[3/3] Set environment variables"
    echo export G_HOME="${install_dir}" >> ${install_dir}/g_profile
    echo export G_EXPERIMENTAL=true >> ${install_dir}/g_profile
    echo export GOROOT='$G_HOME/go' >> ${install_dir}/g_profile
    echo export PATH='$G_HOME:$GOROOT/bin:$PATH' >> ${install_dir}/g_profile
    if [ -n "$G_MIRROR" ]; then
        export G_MIRROR=${GO_MIRROR} >> ${install_dir}/g_profile
    fi
    if [ -x "$(command -v bash)" ]; then
        echo "# ===== set g environment variables =====" >> ${HOME}/.bashrc
        echo "[[ ! -f ${install_dir}/g_profile ]] || source ${install_dir}/g_profile" >> ${HOME}/.bashrc
    fi

    if [ -x "$(command -v zsh)" ]; then
        echo "# ===== set g environment variables =====" >> ${HOME}/.zshrc
        echo "[[ ! -f ${install_dir}/g_profile ]] || source ${install_dir}/g_profile" >> ${HOME}/.zshrc
    fi

    rm -f "${dest_file}"
    exit 0
}

main
