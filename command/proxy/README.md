# 代理工具

## 介绍

1. 支持 thrift 代理，支持 thrift 抓包
2. 支持 HTTP/HTTPS 代理，支持 HTTPS 抓包(使用的MITM)，压测后测试基本没啥劣化(10ms以内吧)

备注: 如果你环境允许的话推荐用 [mitmproxy](https://mitmproxy.org/) 实现HTTPS 抓包，这里造轮子的原因是因为 mitmproxy 依赖太重了，服务器安装太麻烦了！

```shell
~/go/src/github.com/anthony-dong/golang devtool proxy --help
Name: Proxy thrift requests

Usage: devtool proxy [flags]

Options:
  -h, --help            help for proxy
  -l, --listen string   The proxy listen addr. (default ":8080")
      --output string   the output position of the packet (simple/format/json/@file). (default "format")
      --remote string   The remote(thrift) addr.
      --type string     the proxy type. thrift/http/https (default "http")
```

## HTTPS/HTTP 代理

1. 开启代理

```shell
~/go/src/github.com/anthony-dong/golang devtool proxy
[INFO] 03:01:18.829 注意: 安装手册
# 如何配置代理
export http_proxy=http://localhost:8080
export https_proxy=http://localhost:8080
# 如何下载证书 (仅需要操作一次)
wget -e http_proxy=localhost:8080 http://devtool.mitm/cert/pem -O devtool-ca-cert.pem
sudo mv devtool-ca-cert.pem /usr/local/share/ca-certificates/devtool.crt
sudo update-ca-certificates
.... 请初始化后再使用才能生效 ....


```

3. 发起请求

```shell
export http_proxy=http://localhost:8080
export https_proxy=http://localhost:8080
curl 'https://api.github.com/search/users?q=anthony-dong' -v
```

3. 收到抓包信息

```shell
[ID-3] [14:50:31.242] [GET] [https://api.github.com/search/users] [200] [270.258ms]
> GET /search/users?q=anthony-dong HTTP/1.1
> Host: api.github.com
> User-Agent: curl/7.77.0
> Accept: */*
>
< HTTP/1.1 200
< Content-Length: 1201
< Content-Type: application/json; charset=utf-8
< Server: GitHub.com
< Date: Sun, 03 Mar 2024 06:50:28 GMT
< Cache-Control: no-cache
< Vary: Accept, Accept-Encoding, Accept, X-Requested-With
< X-Github-Media-Type: github.v3; format=json
< X-Github-Api-Version-Selected: 2022-11-28
< Access-Control-Expose-Headers: ETag, Link, Location, Retry-After, X-GitHub-OTP, X-RateLimit-Limit, X-RateLimit-Remaining, X-RateLimit-Used, X-RateLimit-Resource, X-RateLimit-Reset, X-OAuth-Scopes, X-Accepted-OAuth-Scopes, X-Poll-Interval, X-GitHub-Media-Type, X-GitHub-SSO, X-GitHub-Request-Id, Deprecation, Sunset
< Access-Control-Allow-Origin: *
< Strict-Transport-Security: max-age=31536000; includeSubdomains; preload
< X-Frame-Options: deny
< X-Content-Type-Options: nosniff
< X-Xss-Protection: 0
< Referrer-Policy: origin-when-cross-origin, strict-origin-when-cross-origin
< Content-Security-Policy: default-src 'none'
< X-Ratelimit-Limit: 10
< X-Ratelimit-Remaining: 6
< X-Ratelimit-Reset: 1709448660
< X-Ratelimit-Resource: search
< X-Ratelimit-Used: 4
< Accept-Ranges: bytes
< X-Github-Request-Id: CB8A:22F5AD:14B9C5E:161F3D0:65E41DB7
<
{
  "total_count": 1,
  "incomplete_results": false,
  "items": [
    {
      "login": "Anthony-Dong",
      "id": 52623370,
      "node_id": "MDQ6VXNlcjUyNjIzMzcw",
      "avatar_url": "https://avatars.githubusercontent.com/u/52623370?v=4",
      "gravatar_id": "",
      "url": "https://api.github.com/users/Anthony-Dong",
      "html_url": "https://github.com/Anthony-Dong",
      "followers_url": "https://api.github.com/users/Anthony-Dong/followers",
      "following_url": "https://api.github.com/users/Anthony-Dong/following{/other_user}",
      "gists_url": "https://api.github.com/users/Anthony-Dong/gists{/gist_id}",
      "starred_url": "https://api.github.com/users/Anthony-Dong/starred{/owner}{/repo}",
      "subscriptions_url": "https://api.github.com/users/Anthony-Dong/subscriptions",
      "organizations_url": "https://api.github.com/users/Anthony-Dong/orgs",
      "repos_url": "https://api.github.com/users/Anthony-Dong/repos",
      "events_url": "https://api.github.com/users/Anthony-Dong/events{/privacy}",
      "received_events_url": "https://api.github.com/users/Anthony-Dong/received_events",
      "type": "User",
      "site_admin": false,
      "score": 1.0
    }
  ]
}
```

## Thrift 代理

```shell
~ devtool proxy --type thrift --remote localhost:8888 -v
[ID-1] [14:54:16.435] [::1]:53307->[::1]:8888 [echo] [OK] [430.65µs]
> CALL echo UnframedBinary
> ID: 1
>
{
	"1_STRUCT": {
		"1_STRING": "my request"
	}
}
< REPLY echo UnframedBinary
< ID: 1
<
{
	"0_STRUCT": {
		"1_STRING": "my request"
	}
}

[ID-2] [14:54:16.436] [::1]:53307->[::1]:8888 [echo] [OK] [184.46µs]
> CALL echo UnframedBinary
> ID: 2
>
{
	"1_STRUCT": {
		"1_STRING": "my request"
	}
}
< REPLY echo UnframedBinary
< ID: 2
<
{
	"0_STRUCT": {
		"1_STRING": "my request"
	}
}
```

备注：如何启动测试的服务, 下载 https://github.com/cloudwego/kitex-examples 项目

- go run hello/main.go hello/handler.go # 启动service
- go run async_call/client/main.go # 注意把代码里的端口改成8080