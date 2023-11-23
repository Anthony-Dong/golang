# Golang

# 介绍

本仓库是个人的一个日常学习Golang的一个仓库，内部包含了一些公共库，其中cli工具 `devtool` 方便平时日常开发！devtool是一个强大的cli工具，包罗万象！

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
└── pkg // 公共包
```

# [devtool](cli/devtool)

```shell
➜  ~ devtool --help
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
  upload      File upload tool

Options:
      --config-file string   Set the config file
  -h, --help                 help for devtool
      --log-level string     Set the log level in [fatal|error|warn|info|debug]
  -v, --verbose              Turn on verbose mode
      --version              version for devtool

Use "devtool COMMAND --help" for more information about a command.

To get more help with devtool, check out our guides at https://github.com/anthony-dong/golang
```

