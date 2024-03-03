# thrift client

# 介绍

支持快速的发起thrift请求，就像curl命令一样 !!!

```shell
~/go/src/github.com/anthony-dong/golang devtool curl --help
Name: Send thrift like curl

Usage: devtool curl [flags]

Examples:
curl --url 'thrift://xxx.xxx.xxx/RPCMethod?addr=localhost:8888&env=prod&cluster=default' --header 'h1: v1' --header 'h2: v2' --data '{"k1": "v1"} -v'

Options:
      --branch string      The Remote IDL branch/version/commit(if cli supports it)
      --data string        The request body
      --example            New request example data
  -H, --header strings     The request header
  -h, --help               help for curl
      --idl string         The main IDL local path
      --modify             Enable the cli to modify the request
      --timeout duration   The request timeout (default 3m0s)
      --url string         The request url

Global Options:
      --config-file string   Set the config file
      --log-level string     Set the log level in [debug|info|notice|warn|error] (default "info")
  -v, --verbose              Turn on verbose mode

To get more help with devtool, check out our guides at https://github.com/anthony-dong/golang
```

# 快速发起测试

```shell
~/go/src/github.com/anthony-dong/golang devtool curl --url 'thrift://EchoService/echo?addr=localhost:8888' --idl example/kitex-examples/echo.thrift --data '{"message": "hello world"}'
[INFO] 17:49:19.548 request info
{
    "protocol": "thrift",
    "service": "EchoService",
    "rpc_method": "echo",
    "body": {
        "message": "hello world"
    },
    "addr": "localhost:8888",
    "timeout": "3m0s"
}
[INFO] 17:49:19.550 spend 1.605969ms
[INFO] 17:49:19.550 response body
{
    "message": "hello world"
}
```


# 获取测试示例
```shell
~/go/src/github.com/anthony-dong/golang devtool curl --url 'thrift://EchoService/echo?addr=localhost:8888' --idl example/kitex-examples/echo.thrift --example
[INFO] 17:58:43.218 request info
{
    "protocol": "thrift",
    "service": "EchoService",
    "rpc_method": "echo",
    "addr": "localhost:8888",
    "timeout": "3m0s"
}
[INFO] 17:58:43.219 new request example
{
    "message": ""
}
```