# devtool

# 介绍

devtool 是一个强大的Cli工具，其中包含了日常开发中的一些可能涉及到的高频工具，这里避免重复造轮子，所以一般就是日常用的一些工具但是市场上没有符合自己需求的，得自己写！

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

# 如何下载

1. `go install`  下载

```shell
CGO_ENABLED=1 go install -v github.com/anthony-dong/golang/cli/devtool@latest
```

2. release 下载 https://github.com/anthony-dong/golang/releases
3. 源码下载

```shell
# 1. 下载源码
git clone --depth 1 https://github.com/anthony-dong/golang.git

# 2. 安装
cd golang && make install
```

4. 如果运行/编译报如下错误

```shell
github.com/google/gopacket/pcap
# github.com/google/gopacket/pcap
pkg/mod/github.com/google/gopacket@v1.1.19/pcap/pcap_unix.go:34:10: fatal error: pcap.h: No such file or directory
 #include <pcap.h>
          ^~~~~~~~
compilation terminated.
```

- linux(Debian) 环境可以执行, mac应该默认就自带了pcap

```shell
# 1. update
sudo apt update

# 2. install
sudo apt-get install -y libpcap-dev
```

# 配置文件

优先级顺序：

1. 读取  `--config-file` 参数传递的配置文件地址
2. 读取 `$HOME/.devtool/config.yaml`  
3. 读取 `$(pwd)/.devtool.yaml`

类型定义：[config.go](../../command/config.go)

# 工具介绍

## [编解码工具 - codec ](../../command/codec)

## [Go开发工具 - gotool](../../command/gotool) 

## [写博客工具 - hexo](../../command/hexo)

## [流量解析工具 - tcpdump](../../command/tcpdump)

## [任务模版工具 - run](../../command/run)

## [文件上传工具 - upload](../../command/upload)

## [JSON工具 - json](../../command/jsontool)
