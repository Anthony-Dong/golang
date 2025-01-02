# Bazel Rule

bazel 的规则定义语法实际上非常优与 cmake，cmake 会显得非常的丑陋且非常臃肿，因此很多开源项目都封装了 cmake rule，方便进行依赖管理.

# 快速开始

```cmake
# 添加依赖
include (cmake/cc_library.cmake)
include (cmake/cc_binary.cmake)
include (cmake/cc_test.cmake)

# 设置项目src路径，方便项目下别的文件include
list (APPEND CUSTOM_PROJECT_SOURCE_DIR ${CMAKE_CURRENT_SOURCE_DIR})
```

# cc_library

1. [bazel](https://bazel.build/reference/be/c-cpp#cc_library)

```build
load("@rules_cc//cc:defs.bzl", "cc_library")

cc_library(
    name = "times",
    srcs = [
        "times.cpp",
    ],
    hdrs = ["times.h"],
    visibility = ["//visibility:public"],
    deps = [
        "//utils:times",
        "@com_google_absl//absl/strings",
        "@com_google_absl//absl/strings:str_format",
    ],
)
```

2. [cmake](cc_library.cmake)

```cmake
cc_library (
        NAME cpp_network
        ALIAS cpp::network
        SRCS listener.cpp
        HDRS header.h utils.h task_queue.h listener.h
        DEPS event cpp::utils
)
```

# cc_binary

1. [bazel](https://bazel.build/reference/be/c-cpp#cc_binary)

```build
load("@rules_cc//cc:defs.bzl", "cc_binary") # 高版本的bazel不需要load cc_binary

cc_binary(
    name = "main",
    srcs = [
        "main.cpp",
    ],
    deps = [
        "//utils:times",
        "@com_google_absl//absl/strings",
        "@com_google_absl//absl/strings:str_format",
    ],
)
```

2. [cmake](cc_binary.cmake)

```cmake
cc_binary (
        NAME network_nio_main
        SRCS network/nio_main.cpp
        DEPS cpp::network cpp::log
)
```

# cc_test

1. [bazel](https://bazel.build/reference/be/c-cpp#cc_test)

> https://google.github.io/googletest/quickstart-bazel.html

```build
cc_test(
  name = "hello_test",
  size = "small",
  srcs = ["hello_test.cc"],
  deps = ["@com_google_googletest//:gtest_main"],
)
```

2. [cmake](cc_test.cmake)

> https://google.github.io/googletest/quickstart-cmake.html

```cmake
include(FetchContent)
FetchContent_Declare(
  googletest
  URL https://github.com/google/googletest/archive/03597a01ee50ed33e9dfd640b249b4be3799d395.zip
)
set(gtest_force_shared_crt ON CACHE BOOL "" FORCE)
FetchContent_MakeAvailable(googletest)

cc_test (
        NAME utils_time_test
        SRCS utils/time_test.cpp
        DEPS cpp::utils
)
```

# TODO

- 生成`pkgconfig`文件
