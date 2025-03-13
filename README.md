# Golang

# 介绍

本仓库是个人的一个日常学习Golang的一个仓库，内部包含了一些公共库，其中cli工具 `devtool` 方便平时日常开发！devtool是一个强大的cli工具，包罗万象！

# 如何使用

```shell
go get -v github.com/anthony-dong/golang
```

# 项目结构

```shell
➜  golang git:(master) tree -L 1 .
.
├── Makefile // 开发脚本
├── README.md
├── bin // 二进制产物
├── build.sh
├── cli // cli工具
├── command // 命令
├── go.mod 
├── go.sum
└── pkg
    ├── bufutils
    ├── codec // 编解码
    ├── collections // 集合
    ├── consts // 常量
    ├── httpclient // http client
    ├── idl // idl
    ├── internal // 内部包
    ├── logs // 日志
    ├── rpc // rpc 调用
    ├── tcpdump // tcpduump
    ├── tools // 外部工具
    └── utils // 工具
```

# [devtool](cli/devtool)

如何下载:  `go install -v github.com/anthony-dong/golang/cli/devtool@master`  或者参考[此文档](cli/devtool)

```shell
➜  devtool git:(master) devtool --help
Usage: devtool [OPTIONS] COMMAND

Commands:
  codec       The Encode and Decode data tool
  gen         Auto compile thrift、protobuf IDL
  go          The golang tools
  help        Help about any command
  hexo        The Hexo tool
  json        The Json tool
  run         Run task templates
  tcpdump     Decode tcpdump file & stream
  turl        Send thrift request like curl
  upload      File upload tool

Options:
      --config-file string   Set the config file
  -h, --help                 help for devtool
      --log-level string     Set the log level in [debug|info|notice|warn|error] (default "info")
  -v, --verbose              Turn on verbose mode
      --version              version for devtool

Use "devtool COMMAND --help" for more information about a command.

To get more help with devtool, check out our guides at https://github.com/anthony-dong/golang
```

# [tcpdump_tools](cli/tcpdump_tools)

1. 安装

```shell
CGO_ENABLED=1 go install -v github.com/anthony-dong/golang/cli/tcpdump_tools@master
```

2. 使用

```shell
~ tcpdump_tools -h
decode tcpdump file, help doc: https://github.com/anthony-dong/golang/tree/master/cli/tcpdump_tools

Usage:
  tcpdump_tools [-r file] [-v] [-X] [--max dump size] [flags]

Examples:
  tcpdump 'port 8080' -X -l -n | tcpdump_tools

Flags:
  -X, --dump          Enable Display payload details with hexdump.
  -r, --file string   The packets file, eg: tcpdump_xxx_file.pcap.
  -h, --help          help for tcpdump_tools
      --max int       The hexdump max size
  -v, --verbose       Enable Display decoded details.
```

# protoc-gen-console

> 将protoc输出到文件中

1. 安装

```shell
go install -v github.com/anthony-dong/golang/cli/protoc-gen-console@master
```

2. 使用

```shell
protoc -I . --plugin=protoc-gen-console=${HOME}/go/bin/protoc-gen-console --console_opt=output=output/out.
json  --console_opt=disable_source_code=1  --console_opt=params1=1,params2=2  --console_out=. pkg/idl/test/text.proto
```