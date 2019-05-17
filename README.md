# g
[![Build Status](https://travis-ci.org/voidint/g.svg?branch=master)](https://travis-ci.org/voidint/g)
[![GoDoc](https://godoc.org/github.com/voidint/g?status.svg)](https://godoc.org/github.com/voidint/g)
[![codecov](https://codecov.io/gh/voidint/g/branch/master/graph/badge.svg)](https://codecov.io/gh/voidint/g)
[![codebeat badge](https://codebeat.co/badges/0b4bf243-95da-444c-b163-6cb8a35d1f8d)](https://codebeat.co/projects/github-com-voidint-g-master)
[![Go Report Card](https://goreportcard.com/badge/github.com/voidint/g)](https://goreportcard.com/report/github.com/voidint/g)

`g`是一个Linux、macOS、Windows下的命令行工具，可以提供一个便捷的多版本[go](https://golang.org/)环境的管理和切换。


## 特性
- 支持列出可供安装的go版本号
- 支持列出已安装的go版本号
- 支持在本地安装多个go版本
- 支持卸载已安装的go版本
- 支持在已安装的go版本之间自由切换

## 安装
### 自动化安装
- Linux/macOS

    ```shell
    $ wget -qO- https://github.com/voidint/g/blob/master/install.sh | bash
    $ echo "unalias g" >> ~/.bashrc # 可选。若其他程序（如'git'）使用了'g'作为别名。
    $ source ~/.bashrc # 或者 source ~/.zshrc
    ```

### 手动安装
- 下载对应平台的[二进制压缩包](https://github.com/voidint/g/releases)。
- 将压缩包解压至`PATH`环境变量目录下，如`/usr/local/bin`。
- 编辑shell环境配置文件（`~/.bashrc`、`~/.zshrc`...）

    ```shell
    $ cat>>~/.bashrc<<EOF
    export GOROOT="${HOME}/.g/go"
    export PATH="${HOME}/.g/go/bin:$PATH"
    export G_MIRROR=https://golang.google.cn/dl/
    EOF
    ```

## 使用
查询当前可供安装的`stable`状态的go版本

```shell
$ g ls-remote stable
1.11.9
1.12.4
```

安装目标go版本`1.12.4`

```shell
$ g install 1.12.4
installed successfully
$ go version
go version go1.12.4 darwin/amd64
```


查询已安装的go版本

```shell
$ g ls
1.12.4
```

查询可供安装的所有go版本

```shell
$ g ls-remote
1
1.2.2
1.3
1.3.1
...    // 省略若干版本
1.11.7
1.11.8
1.11.9
1.12
1.12.1
1.12.2
1.12.3
1.12.4
```

安装目标go版本`1.11.9`

```shell
$ g install 1.11.9
installed successfully
$ go version
go version go1.11.9 darwin/amd64
```

切换到另一个已安装的go版本

```shell
$ g ls
1.11.9
1.12.4
$ g use 1.12.4
go version go1.12.4 darwin/amd64

```

卸载一个已安装的go版本

```shell
g uninstall 1.11.9
Uninstall successfully
```
## FAQ
- 环境变量`G_MIRROR`有什么作用？

    由于中国大陆无法自由访问Golang官网，导致查询及下载go版本都变得困难，因此可以通过该环境变量指定一个镜像站点（如`https://golang.google.cn/dl/`），g将从该站点查询、下载可用的go版本。

- 支持源代码编译安装吗？

    不支持


## 鸣谢
感谢[nvm](https://github.com/nvm-sh/nvm)、[n](https://github.com/tj/n)、[rvm](https://github.com/rvm/rvm)等工具提供的宝贵思路。