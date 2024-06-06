#!/usr/bin/env bash

# 注意: 前置条件安装 spdlog
# wget https://github.com/gabime/spdlog/archive/refs/tags/v1.12.0.tar.gz -O- | tar -zxvf -
# cd spdlog
# mkdir -p build && cd build && cmake ../ && sudo make install

set -e

if [ "$1" == "bazel" ]; then
  # 安装 https://github.com/bazelbuild/bazelisk
  bazel run //:main
  exit 0
fi


if [ "$1" == "devtool" ]; then
  devtool cpp --src times.cpp --src utils.cpp --src main.cpp --output output/main --cxxopt -std=c++17 --cxxopt -I/usr/local/include --linkopt -L/usr/local/lib --linkopt -lspdlog -j8 -r -v
  exit 0
fi

echo "Usage:"
echo "  $0 [bazel|devtool]"
exit 1