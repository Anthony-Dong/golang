# Hexo

# 介绍

hexo命令支持快速将我的笔记打包成一个hexo项目，本人日常开发会写了很多的文章，但是只有少数需要发布到网上，因此我有两个目录，一个是工作笔记，一个是hexo-home，我只需要在工作笔记中加入 页眉，然后执行hexo命令就可以帮我需要发布的文章发布出去

# 使用说明

```shell
➜  ~ devtool hexo  build --help
Name: Build the markdown project to hexo

Usage: devtool hexo build [flags]

Options:
  -d, --dir string            The markdown project dir (required)
  -h, --help                  help for build
  -k, --keyword stringArray   The keyword need clear, eg: baidu => xxxxx, read from command and load config
  -t, --target_dir string     The hexo post dir (required)

Global Options:
      --config-file string   Set the config file
      --log-level string     Set the log level in [fatal|error|warn|info|debug]
  -v, --verbose              Turn on verbose mode
```

# 我是如何使用的

1. 我的笔记目录

```shell
➜  ~ cd ~/note/note
➜  note git:(master) ✗ tree -L 1 .
.
├── LICENSE
├── Makefile
├── README.md
├── bin // 一些二进制工具
├── code // 平时练习的代码
├── hexo-home // hexo 发布的目录
├── .... // 个人笔记
├── .... // 个人笔记
```

2. 执行命令

```shell
devtool hexo build --dir ./ --target_dir ./hexo-home/source/_posts
```

