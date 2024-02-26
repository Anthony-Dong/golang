## 如何快速启动kitex项目

1. 下载 kitex-example 项目 github.com/cloudwego/kitex-examples

```shell
> git clone github.com/cloudwego/kitex-examples

# 版本比较低 go1.16，最新的要求go1.21
> git checkout v0.2.3
```

2. 执行hello程序
```shell
~/go/src/github.com/anthony-dong/golang/example/kitex-examples go run hello/main.go hello/handler.go
2024/02/07 17:49:08.325556 server.go:83: [Info] KITEX: server listen at addr=[::]:8888
```