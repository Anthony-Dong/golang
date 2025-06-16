# devtool

# 介绍

devtool 是一个强大的开发者工具，其中包含了日常开发中的一些可能涉及到的高频工具，比如编解码工具、抓包工具、发送请求工具、代理工具、文件存储工具等！

```shell
➜  git:(master) devtool 
Name: A collection of development tools

Usage: devtool [OPTIONS] COMMAND

Commands:
  codec       Tools for encoding and decoding data in various formats
  codegen     Generate code from IDL files (Thrift, Protobuf)
  curl        Send Thrift requests similar to cURL
  diff        Diff utilities for various data formats
  git         Git related tools and utilities
  golang      Go language related tools and utilities
  help        Help about any command
  hexo        Hexo static site generator tools
  mock        Generate mock data for various formats
  proxy       HTTP/HTTPS, Thrift, and FileSystem proxy tools
  run         Run tasks defined in template files
  static      Serve static files
  strace      Trace system calls and signals
  tcp         TCP related tools
  trimer      Trim or filter data from various formats
  upload      File upload utility
  wal         Commands for interacting with Write-Ahead Log (WAL) files

Options:
      --config-file string   Set the config file
  -h, --help                 help for devtool
      --log-level string     Set the log level in [debug|info|notice|warn|error] (default "info")
  -v, --verbose              Turn on verbose mode
      --version              version for devtool

Use "devtool COMMAND --help" for more information about a command.

To get more help with devtool, check out our guides at https://github.com/anthony-dong/golang
```

# 如何下载

1. `go install`  下载

```shell
# 注意Go版本大于等于1.18
CGO_ENABLED=1 go install -v github.com/anthony-dong/golang/cli/devtool@latest
```

2. release 下载 https://github.com/anthony-dong/golang/releases

```shell
# 1. update
sudo apt update

# 2. install
sudo apt-get install -y libpcap-dev
```

# 配置文件

优先级顺序：

1. 读取  `--config-file` 参数传递的配置文件地址
2. 读取 `$(pwd)/.devtool.yaml`
3. 读取 `dirname($0)/.devtool.yaml`
4. 读取 `$HOME/.devtool/config.yaml`

类型定义：[config.go](../../command/config.go)

# 工具介绍

## [编解码工具 - codec ](../../command/codec)

## [Go开发工具 - golang](../../command/golang)

## [写博客工具 - hexo](../../command/hexo)

## [流量解析工具 - tcpdump](../../command/tcpdump)

## [任务模版工具 - run](../../command/run)

## [文件上传工具 - upload](../../command/upload)

## [Thrift/HTTPS/HTTP代理和抓包工具](../../command/proxy)

## [像curl一样发起Thrift请求](../../command/curl)