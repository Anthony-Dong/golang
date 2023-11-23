# JSON

# 介绍

1. json-reader

```shell
➜  golib git:(master) ✗ echo '{"k1":{"k2":"v2"}}'  | devtool json --path k1 --pretty
{
  "k2": "v2"
}
```

2. curl + json 写文件

```json
curl --request GET 'https://xxxx.xxxx.org/api/v1/test?xxx=xxxx' \
--header 'Cookie: xxxxx' \
--header 'x-xxx-x: 1' |  devtool json --pretty --path k1.v1.v2  | devtool json writer 
```

# 使用说明

```shell
➜  ~ devtool json --help
Name: The Json Tool

  - Query JSON values: https://jqlang.github.io/jq
  - Terminal JSON viewer: https://github.com/antonmedv/fx
  - JSON Path: https://github.com/tidwall/gjson

Usage: devtool json [flags]

Options:
  -h, --help          help for json
      --path string   set specified path
      --pretty        set pretty json
```