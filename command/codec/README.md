# Codec

# 介绍

目前支持 `thrift`,`pb`,`br`,`base64`,`gizp`,`snappy`,`url`,`md5`,`hex` 等多种消息解析，比较适合我们日常开发中，经常性的会解析各种数据！使用这个命令可以帮助你实现快速的转换！

例如我们将一个 thrift/pb 的消息报文，是base64编码的，然后通过 base64 decode，然后通过 thrift/pb decode，最后通过 json pretty 打印可以看到如下结果！

```shell
➜  echo "AAAAEYIhAQRUZXN0HBwWAhUCAAAA" | devtool codec base64 --decode | devtool codec thrift | devtool json pretty
{
  "method": "Test",
  "seq_id": 1,
  "payload": {
    "1_STRUCT": {
      "1_STRUCT": {
        "1_I64": 1,
        "2_I32": 1
      }
    }
  },
  "message_type": "call",
  "protocol": "FramedCompact"
}

➜  echo "CgVoZWxsbxCIBEIDCIgE" | devtool codec base64 --decode | devtool codec pb | jq            
{
  "1": "hello",
  "2": 520,
  "8": {
    "1": 520
  }
}
```

# 使用说明

```shell
➜  devtool codec --help                                                                             
Name: The Encode and Decode data tool

Usage: devtool codec [OPTIONS] COMMAND

Commands:
  base64      base64 codec
  br          br codec
  gizp        gizp codec
  hex         hex codec
  md5         md5 codec
  pb          decode protobuf protocol
  snappy      snappy codec
  thrift      decode thrift protocol
  url         url codec

Options:
  -h, --help   help for codec
```