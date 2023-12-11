# CPP 构建流程

1. 源程序（source code）
2. 预处理器（preprocessor）- 预处理一般是 宏/头文件会直接展开到源文件， 例如：cpp main.c main.i 或者 gcc -E main.c -o main.i
3. 编译器（compiler）- 将预处理后的源文件翻译成汇编文件，例如：ccl main.i -Og -o main.s 或者 gcc -S main.i -o main.s
4. 汇编程序（assembler）- 汇编器(as)将main.s翻译为main.o，例如 as -o main.o main.s
5. 目标程序（object code） - 编译的产物：可重定位目标文件，实际上就是存在一些未定义的变量/函数等
6. 连接器（链接器，Linker）- 链接，实际上就是解决上面未定义的函数/变量等，例如：ld -o main main.o
7. 可执行程序（executables） - 链接的产物：静态库/动态库/可执行文件

其中 gcc/clang 等工具实际上就是整合了其全部过程，详细过程实际上你通过 -v 参数就能看到

上面参考了 https://open.toutiao.com/a6976592795305902629 这个文章

> 传统的CPP构建流程(非c++20的module模式)

1. 编译：输入源文件 (*.cpp *.hpp) 文件，输出 object(可重定位目标文件其实也就是汇编-二进制文件) 文件 （此阶段支持并发编译）

```shell
# -I 指定头文件的搜索路径（默认会携带上当前路径`.`，所以不需要指定）
clang++ -Wall -std=c++17 -O0 -g -I/usr/local/include -c times.cpp -o output/times.o
```

2. 链接:

- 可执行文件:

```shell
# -L 指定依赖库的搜索地址
# -l 表示依赖哪些库(/usr/local/lib/libspdlog.a 就可以通过 spdlog来表示)
clang++ -o output/main output/times.o output/utils.o output/main.o -L/usr/local/lib -lspdlog
```

- 静态库:  指的是将多个object文件，打包成静态库一个archive文件(.a)，动态库shared object文件一般是(.so or windows .ddl)

```shell
# -r 表示在库中插入模块
# -c 表示强制创建模块
ar -r -c output/utils.a output/times.o output/utils.o

# libtool 也是一个可以打包静态库的工具 
libtool -static -s  -o output/utils.a output/times.o output/utils.o
```

# 帮助

```shell
~/go/src/github.com/anthony-dong/golang devtool cpp --help
Name: Supports fast compile and running of a cpp file

Usage: devtool cpp [--src .cpp] [--hdr .h] [-o binary] [--type binary] [--thread number] [-r] [flags] -- [build flags ... ]

Options:
      --hdr strings            The source header files
  -h, --help                   help for cpp
  -I, --include strings        Add directory to include search path
  -l, --link strings           Add link library
  -L, --link_include strings   Add directory to library search path
  -o, --output string          The output file
      --release                Set the compile type is release
  -r, --run                    Exec output binary file
      --src strings            The source files
  -j, --thread int             The number of compiled threads (default 1)
      --type string            The link object type [binary|library] (default "binary")

Global Options:
      --config-file string   Set the config file
      --log-level string     Set the log level in [debug|info|notice|warn|error] (default "info")
  -v, --verbose              Turn on verbose mode

To get more help with devtool, check out our guides at https://github.com/anthony-dong/golang
```

# 测试

## 测试  devtool

```shell
~/go/src/github.com/anthony-dong/golang/command/cpp/test bash  build.sh devtool
[DEBUG] 15:18:37.789 init config success. filename: /home/anthony-dong/.devtool/config.yaml
[DEBUG] 15:18:37.789 start cmd: devtool, cmd.args: [], os.args: ["devtool","cpp","--src","times.cpp","--src","utils.cpp","--src","main.cpp","-v","-r","-L/usr/local/lib","-I/usr/local/include","-lspdlog","-j8"]
[DEBUG] 15:18:37.789 output: main, link type: binary, thread number: 8, tools config: {
    "CXX": "clang++",
    "CC": "clang",
    "Dir": "/home/anthony-dong/go/src/github.com/anthony-dong/golang/command/cpp/test",
    "SRCS": [
        "times.cpp",
        "utils.cpp",
        "main.cpp"
    ],
    "BuildIncludes": [
        "/usr/local/include"
    ],
    "LinkIncludes": [
        "/usr/local/lib"
    ],
    "LinkLibraries": [
        "spdlog"
    ],
    "CompileType": "debug"
}
[INFO] 15:18:37.789 Compile: clang++ -Wall -std=c++17 -O0 -g -I/usr/local/include -c main.cpp -o output/main.o
[INFO] 15:18:37.789 Compile: clang++ -Wall -std=c++17 -O0 -g -I/usr/local/include -c utils.cpp -o output/utils.o
[INFO] 15:18:37.789 Compile: clang++ -Wall -std=c++17 -O0 -g -I/usr/local/include -c times.cpp -o output/times.o
[INFO] 15:18:43.591 Link: clang++ -o output/main output/times.o output/utils.o output/main.o -L/usr/local/lib -lspdlog
[INFO] 15:18:44.066 Run: output/main
test::times::v1.0.0
test::utils::v1.0.0
[2023-12-03 15:18:44.069] [info] hello world
```

## 测试 bazel

```shell
~/go/src/github.com/anthony-dong/golang/command/cpp/test bash  build.sh bazel
INFO: Analyzed target //:main (36 packages loaded, 162 targets configured).
INFO: Found 1 target...
Target //:main up-to-date:
  bazel-bin/main
INFO: Elapsed time: 2.973s, Critical Path: 0.41s
INFO: 8 processes: 4 internal, 4 linux-sandbox.
INFO: Build completed successfully, 8 total actions
INFO: Running command line: bazel-bin/main
test::times::v1.0.0
test::utils::v1.0.0
```

