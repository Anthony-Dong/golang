# 任务模版

# 介绍

**任务模版**在现实中还是非常刚需的，因为任务是高度复用化的，不论是github的workflow还是自研的任务流水线，或多或少都要定义一堆的任务，但是其实会发现一点需要大量定义一些重复的任务配置，这样就会导致维护困难。任务模版+`include`指令可以实现任务的高度复用化，具体可以看下文的例子～～

其中任务模版会有几个概念：

1. 模版变量：可以根据模版变化渲染任务配置

- HOME: 用户根路径

2. 模版函数：一些复杂的逻辑可以通过模版函数来实现

- IF: 三元表达式，例如：`{{ IF .DEBUG "1" "--debug" "-" }}` 表示：如果`DEBUG=1`则输出`--debug`否则是`-`, 注意代码中会将`-`
  认为是空值`""`

3. Include：可以引用别的模版配置，实现组合的效果

Todo：

1. 支持并行逻辑，例如下图，即我们需要一个 组任务的概念，目前还没有此类需要，所以暂时没有支持

![image-20231123125048011](https://tyut.oss-accelerate.aliyuncs.com/image/2023/11-23/22f740c98323415c83e128980b91baeb.png)

2. 解决循环引用的问题，既然我们支持了 `include` 指令，必不可少的就是循环引用检测

# 使用说明

```shell
➜  ~ devtool run --help
Name: Run task templates

Usage: devtool run [flags]

Options:
  -c, --config string     The task config file
  -d, --debug             Enable Debug
  -h, --help              help for run
  -I, --include strings   Add an config search path for includes
  -l, --list              List task config files
      --var strings       Define template variables

Global Options:
      --config-file string   Set the config file
      --log-level string     Set the log level in [fatal|error|warn|info|debug]
  -v, --verbose              Turn on verbose mode      
```

核心介绍一下`-I`参数，你可以指定多个include路径，类似于`gcc -I`这个参数一样搜索头文件，因为我们可能会定义一些公共模块在一个通用目录下。**include指令**的搜索逻辑是优先当前引用文件的路径，其次是指定-I的路径下搜索！

# 例子

> 相关配置文件在 [test](./test) 文件夹下

1. 背景：我们这边有多个API服务，他的框架都一样，因此大概它的启动脚本都一样，字节这边go-api服务基本使用的是hertz/ginex，这里以ginex框架为例子
   1. 第一步: 定义 `sd.yaml` , 它是一个服务发现[service discovery]的任务模版，主要是用于注册服务


```yaml
- name: "注册 {{ .SERVICE }} 服务" # 任务名称，其中 {{ .XXX }} 模版表达式，使用的是Go语言的自带的模版
  cmd: "sd" # 命令名称
  daemon: true # 表示后台运行这个程序，当整个程序组结束时会kill掉此进程的
  args:
    - "consul"
    - "--service"
    - "{{ .SERVICE }}"
    - "--env"
    - "{{ .ENV }}"
    - "--cluster"
    - "{{ .CLUSTER }}"
    - "--port"
    - "{{ .PORT }}"
  vars: # 模版变量
    ENV: "dev_xiaoming"
    CLUSTER: "default"
    PORT: 10086
```

2. 定义 `ginex.yaml` 为 ginex 框架的运行模版

```yaml
- include: "sd.yaml" # include sd.yaml 配置
  vars:
    PORT: 9999 #ginex 服务的端口是9999
- name: "启动 {{ .SERVICE }} 服务"
  cmd: "devtool"
  dir: "{{ .HOME }}/go/src/{{ .REPO }}" # HOME 为全局变量默认是用户根路径
  env:
    - "SSO_DEBUG_USER={{ .USERNAME }}"
    - "LOAD_SERVICE={{ .SERVICE }}"
  args:
    - "go"
    - "run"
    - '{{ IF .DEBUG "1" "--debug" "-" }}' # 通过IF表达式来动态注入--debug参数
    - "--run"
    - "."
    - "--"
    - "-service={{ .SERVICE }}"
    - "-conf-dir=conf"
    - "-log-dir=log"
  vars:
    DEBUG: 0 # 默认是0表示禁止DEBUG
    USERNAME: "xiaoming" # 默认是xiaoming
```

3. 定义Api服务 `github.search.api.yaml` 的配置文件

```yaml
- include: "ginex.yaml"
  vars:
    SERVICE: "github.search.api"
    REPO: "github.com/github/search_api"
    DEBUG: 1
```

4. 测试，可以执行  `devtool run -c command/run/test/github.search.api.yaml -d -v` 来输出配置信息，检查是否正确

```shell
➜  golang git:(master) ✗ devtool run -c command/run/test/github.search.api.yaml -d -v
[DEBUG] 13:25:44.496 start cmd: devtool, cmd.args: [], os.args: ["devtool","run","-c","command/run/test/github.search.api.yaml","-d","-v"]
[DEBUG] 13:25:44.496 includes: ["/Users/bytedance/go/src/github.com/anthony-dong/golang"]. debug file: true.
[INFO] 13:25:44.496 init template vars: {
    "HOME": "/Users/bytedance",
    "PWD": "/Users/bytedance/go/src/github.com/anthony-dong/golang"
}
[INFO] 13:25:44.497 file: /Users/bytedance/go/src/github.com/anthony-dong/golang/command/run/test/github.search.api.yaml
- name: 注册 github.search.api 服务
  cmd: sd
  args:
    - consul
    - --service
    - github.search.api
    - --env
    - dev_xiaoming
    - --cluster
    - default
    - --port
    - "9999"
  vars:
    CLUSTER: default
    ENV: dev_xiaoming
    PORT: "9999"
  daemon: true
- name: 启动 github.search.api 服务
  cmd: devtool
  dir: /Users/bytedance/go/src/github.com/github/search_api
  env:
    - SSO_DEBUG_USER=xiaoming
    - LOAD_SERVICE=github.search.api
  args:
    - go
    - run
    - --debug
    - --run
    - .
    - --
    - -service=github.search.api
    - -conf-dir=conf
    - -log-dir=log
  vars:
    DEBUG: "1"
    USERNAME: xiaoming
```

5. 执行，只需要把 `-d` 参数移除掉就会开始执行任务