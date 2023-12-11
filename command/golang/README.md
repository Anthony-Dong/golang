# Go Tools
# 介绍

快速运行/Debug一个Go项目，适用于远程开发这种

# 使用说明
## run

```shell
➜  ~ devtool go run --help
Name: run golang project

Usage: devtool go run [flags]

Options:
      --debug         enable debug
      --env strings   go test env
  -h, --help          help for run
      --run string    go run pkg name (default ".")

Global Options:
  -v, --verbose              Turn on verbose mode
```


## Test

```shell
➜  ~ devtool go test --help
Name: test golang project

Usage: devtool go test [flags]

Options:
      --debug         enable debug
      --env strings   go test env
  -h, --help          help for test
      --pkg string    go test pkg
      --run string    go test name

Global Options:
  -v, --verbose              Turn on verbose mode      
```