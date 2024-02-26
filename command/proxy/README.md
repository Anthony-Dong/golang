# 代理工具

## 介绍
1. 目前支持快速发起thrift请求，其他的暂时不支持

```shell
~/go/src/github.com/anthony-dong/golang devtool proxy --help
Name: Proxy thrift requests

Usage: devtool proxy [flags]

Options:
  -d, --dial string     The dial addr
  -h, --help            help for proxy
  -l, --listen string   The listen addr

Global Options:
      --config-file string   Set the config file
      --log-level string     Set the log level in [debug|info|notice|warn|error] (default "info")
  -v, --verbose              Turn on verbose mode

To get more help with devtool, check out our guides at https://github.com/anthony-dong/golang
```


## 示例
```shell
~/go/src/github.com/anthony-dong/golang devtool proxy --listen :10086 --dial :8888
>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> [CALL] 2024-02-07 18:10:31 >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
{
    "method": "echo",
    "seq_id": 1,
    "protocol": "UnframedBinary",
    "message_type": "call",
    "payload": {
        "1_STRUCT": {
            "1_STRING": "my request"
        }
    },
    "meta_info": {}
}
<<<<<<<<<<<<<<<<<<<<<<<<<<<<< [REPLY] 2024-02-07 18:10:31 <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
{
    "method": "echo",
    "seq_id": 1,
    "protocol": "UnframedBinary",
    "message_type": "reply",
    "payload": {
        "0_STRUCT": {
            "1_STRING": "my request"
        }
    },
    "meta_info": {}
}
>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> [CALL] 2024-02-07 18:10:31 >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
{
    "method": "echo",
    "seq_id": 2,
    "protocol": "UnframedBinary",
    "message_type": "call",
    "payload": {
        "1_STRUCT": {
            "1_STRING": "my request"
        }
    },
    "meta_info": {}
}
<<<<<<<<<<<<<<<<<<<<<<<<<<<<< [REPLY] 2024-02-07 18:10:31 <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
{
    "method": "echo",
    "seq_id": 2,
    "protocol": "UnframedBinary",
    "message_type": "reply",
    "payload": {
        "0_STRUCT": {
            "1_STRING": "my request"
        }
    },
    "meta_info": {}
}
```

备注：如何启动测试的服务
1. go run hello/main.go hello/handler.go # 启动service
2. go run async_call/client/main.go # 注意把代码里的端口改成10086

## 其他高级特性

1. 代理tcp
```shell
devtool proxy -l :10086 -d 10.37.10.60:8888
```

2. 代理uds
```shell
devtool proxy -l :10086 -d /opt/app.socket
```