cmake_minimum_required(VERSION 3.20)
project(PROJECT_NAME)

set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_STANDARD_REQUIRED ON) # force stdc++=20
set(CMAKE_CXX_EXTENSIONS OFF) # 禁止使用GUN特性

include (cmake/cc_library.cmake)
include (cmake/cc_binary.cmake)
include (cmake/cc_test.cmake)

cc_binary(
    NAME main
    SRCS main.cpp
)