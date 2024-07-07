# g

![GitHub release (latest by date)](https://img.shields.io/github/v/release/voidint/g)
[![GoDoc](https://godoc.org/github.com/voidint/g?status.svg)](https://godoc.org/github.com/voidint/g)
[![codecov](https://codecov.io/gh/voidint/g/branch/master/graph/badge.svg)](https://codecov.io/gh/voidint/g)
[![codebeat badge](https://codebeat.co/badges/0b4bf243-95da-444c-b163-6cb8a35d1f8d)](https://codebeat.co/projects/github-com-voidint-g-master)
[![Go Report Card](https://goreportcard.com/badge/github.com/voidint/g)](https://goreportcard.com/report/github.com/voidint/g)

[简体中文 🇨🇳](./README_CN.md)

**Note:** The master branch may still be under development and may not represent a stable version. Please download stable versions of the source code through tags or download compiled binary executables through [release](https://github.com/voidint/g/releases).

`g` is a command-line tool for Linux, macOS, and Windows that provides convenient management and switching of multiple versions of the [Go](https://golang.org/) environment.

[![asciicast](https://asciinema.org/a/356685.svg)](https://asciinema.org/a/356685)

## Features

- Support for listing available versions of Go for installation
- Support for listing installed versions of Go
- Support for installing multiple versions of Go locally
- Support for uninstalling installed versions of Go
- Support for freely switching between installed versions of Go
- Support for clearing package file cache
- Support for self-updating software (>= 1.5.0)
- Support for clean uninstallation of the software (>= 1.5.0)

## Installation

### Automated Installation

- Linux/macOS (bash/zsh)

  ```shell
  # It is recommended to clear the `GOROOT`, `GOBIN`, and other environment variables before installation.
  $ curl -sSL https://raw.githubusercontent.com/voidint/g/master/install.sh | bash
  $ echo "unalias g" >> ~/.bashrc # Optional. If other programs (such as `git`) have used `g` as an alias.
  $ source "$HOME/.g/env"
  ```

- Windows (pwsh)

  ```pwsh
  $ iwr https://raw.githubusercontent.com/voidint/g/master/install.ps1 -useb | iex
  ```

### Manual Installation(for Linux/macOS)
- Create a directory for `g` (recommended: `~/.g`)
- Download the binary compressed file from [releases](https://github.com/voidint/g/releases) and unzip it into the `bin` subdirectory of the `g` directory (i.e. `~/.g/bin`).
- Write necessary environment variables into `~/.g/env` file.

  ```shell
  $ cat >~/.g/env <<'EOF'
  #!/bin/sh
  # g shell setup
  export GOROOT="${HOME}/.g/go"
  export PATH="${HOME}/.g/bin:${GOROOT}/bin:$PATH"
  export G_MIRROR=https://golang.google.cn/dl/
  EOF
  ```

- Import `~/.g/env` into the shell environment configuration files (e.g. `~/.bashrc`, `~/.zshrc`...).

  ```shell
  $ cat >>~/.bashrc <<'EOF'
  # g shell setup
  if [ -f "${HOME}/.g/env" ]; then
      . "${HOME}/.g/env"
  fi
  EOF
  ```

- Enable environment variables.
  ```shell
  $ source ~/.bashrc # source ~/.zshrc
  ```

### Manual Installation (for Windows PowerShell)

- Create a directory: `mkdir ~/.g/bin`
- Download the binary compressed file for Windows version from [releases](https://github.com/voidint/g/releases), and after unzipping it, put it in the ~/.g/bin directory.
- The default binary file name is `g.exe`, if you have already used `g` as an abbreviation for Git command, you can change `g.exe` to another name, such as `gvm.exe`.
- Run the command `code $PROFILE`, this command will open the default PowerShell configuration file using VSCode.
- Add the following content to the default PowerShell configuration file:

  ```ps1
  $env:GOROOT="$HOME\.g\go"
  $env:Path=-join("$HOME\.g\bin;", "$env:GOROOT\bin;", "$env:Path")
  ```

- Open the PowerShell terminal again, and you can use the `g` or `gvm` command.

## Usage

To query the currently available stable versions of Go for installation:

```shell
$ g ls-remote stable
  1.19.10
  1.20.5
```

To install a specific version of Go (e.g., 1.20.5):

```shell
$ g install 1.14.7
Downloading 100% [===============] (92/92 MB, 12 MB/s)               
Computing checksum with SHA256
Checksums matched
Now using go1.20.5
```

To query the list of installed Go versions:

```shell
$ g ls
  1.19.10
* 1.20.5
```

To list all available Go versions for installation:

```shell
$ g ls-remote
  1
  1.2.2
  1.3
  1.3.1
  ...    
  1.19.10
  1.20rc1
  1.20rc2
  1.20rc3
  1.20
  1.20.1
  1.20.2
  1.20.3
  1.20.4
* 1.20.5
```

To switch to another installed Go version:

```shell
$ g use 1.19.10
go version go1.19.10 darwin/arm64
```

To uninstall a specific installed Go version:

```shell
$ g uninstall 1.19.10
Uninstalled go1.19.10
```

To clear the package file cache for Go installations:

```shell
$ g clean 
Remove go1.18.10.darwin-arm64.tar.gz
Remove go1.19.10.darwin-arm64.tar.gz
Remove go1.20.5.darwin-arm64.tar.gz
```

To view the version information of `g` itself:

``` shell
g version 1.5.0
build: 2023-01-01T21:01:52+08:00
branch: master
commit: cec84a3f4f927adb05018731a6f60063fd2fa216
```

To update `g` software itself:

```shell
$ g self update
You are up to date! g v1.5.0 is the latest version.
```

To uninstall the `g` software itself:

```shell
$ g self uninstall
Are you sure you want to uninstall g? (Y/n)
y
Remove /Users/voidint/.g/bin/g
Remove /Users/voidint/.g
```

## FAQ

- What is the purpose of the environment variable `G_MIRROR`?

  Due to the restricted access to the Golang official website in mainland China, it has become difficult to query and download go versions. Therefore, the environment variable `G_MIRROR` can be used to specify one or multiple mirror sites (separated by commas) from which g will query and download available go versions. The known available mirror sites are as follows:

  - Go official mirror site: https://golang.google.cn/dl/
  - Alibaba Cloud: https://mirrors.aliyun.com/golang/
  - Nanjing University: https://mirrors.nju.edu.cn/golang/
  - Huazhong University of Science and Technology: https://mirrors.hust.edu.cn/golang/
  - University of Science and Technology of China: https://mirrors.ustc.edu.cn/golang/

- What is the purpose of the environment variable `G_EXPERIMENTAL`?

  When the value of this environment variable is set to true, it enables all experimental features.

- What is the purpose of the environment variable `G_HOME`?

  By convention, g uses the `~/.g` directory as its home directory. If you want to customize the home directory (especially for Windows users), you can use the G_HOME environment variable to switch to another directory. Since this feature is still experimental, it requires enabling the experimental feature switch `G_EXPERIMENTAL=true` to take effect. Please note that this solution is not perfect, which is why it is classified as an experimental feature. For more details, please refer to [#18](https://github.com/voidint/g/issues/18).

- On macOS, when installing a go version, g throws an error message saying `[g] Installation package not found.` What is the reason?

  The Go official support for ARM architecture on macOS was introduced in version [1.16](https://go.dev/doc/go1.16#darwin). Therefore, go installation packages of version 1.15 and earlier cannot be installed on ARM-based macOS systems. If you attempt to install these versions, g will throw an error message `[g] Installation package not found.`

- Does g support network proxy?

  Yes, it supports network proxy. You can set the network proxy address in environment variables such as `HTTP_PROXY`, `HTTPS_PROXY`, `http_proxy`, and `https_proxy`.

- Which versions of Windows are supported?

  Since g relies on symbolic links, the operating system must be Windows Vista or above.

- Why doesn't g work after installing it on Windows?

  This may be because the downloaded and installed files are not added to the `$Path`. You need to manually add `$Path` to the user's environment variables. For convenience, you can run the `path.ps1` PowerShell script provided in the project and then restart your computer.

- After installing a go version using g, when running the `go version` command, the output shows a different version than the one installed. Is this a bug?

  This is likely due to an incorrect setting of the `PATH` environment variable in the current shell environment (it is recommended to run `which go` to see the path of the go binary file). By default, the path to the go binary file should be `~/.g/go/bin/go`. If it is not this path, it means that the PATH environment variable is set incorrectly.
  
- Does g support compiling and installing from source code?

  No, it does not support compiling and installing from source code.

## Acknowledgement

Thanks to tools like [nvm](https://github.com/nvm-sh/nvm), [n](https://github.com/tj/n), [rvm](https://github.com/rvm/rvm) for providing valuable ideas.
