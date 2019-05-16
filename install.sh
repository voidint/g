#!/usr/bin/env bash

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
    local dest_file="${HOME}/g1.0.0.${os}-${arch}.tar.gz"

    echo "Download ${dest_file}"
    rm -f ${dest_file}
    wget -P ${HOME} "https://github.com/voidint/g/releases/download/v1.0.0/g1.0.0.${os}-${arch}.tar.gz"
    
    chmod +x ${dest_file}

    mkdir -p "${HOME}/bin"

    tar -xzv -f ${dest_file} -C ${HOME}/bin
    chmod +x ${HOME}/bin/g

    echo "Successfully installed to ${HOME}/bin/g"

    echo "export GOROOT=${HOME}/.g/go" >> ~/.bashrc
    echo "export PATH=${HOME}/bin:${HOME}/.g/go/bin:$PATH" >> ~/.bashrc
    echo "export G_MIRROR=https://golang.google.cn/dl/" >> ~/.bashrc 
    source ~/.bashrc 

    exit 0
}

main