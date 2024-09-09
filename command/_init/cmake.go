package _init

import (
	"github.com/spf13/cobra"

	"github.com/anthony-dong/golang/pkg/utils"
)

func NewCmakeCommand() (*cobra.Command, error) {
	output := ""
	force := false
	cmd := &cobra.Command{
		Use:   "cmake",
		Short: `cmake project`,
		RunE: func(cmd *cobra.Command, args []string) error {
			files := getCmakeFiles()
			for _, file := range files {
				if err := file.Write(output, force); err != nil {
					return err
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&output, "output", "O", utils.GetPwd(), "the output dir")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "force replace file")
	return cmd, nil
}

func getCmakeFiles() []*File {
	return []*File{
		{
			Name:       "CMakeLists.txt",
			IsTemplate: false,
			Content: `cmake_minimum_required(VERSION 3.10)

set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_STANDARD_REQUIRED TRUE)
set(CMAKE_CXX_EXTENSIONS OFF)

if (CMAKE_CXX_COMPILER_ID MATCHES "Clang")
    set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -stdlib=libc++")
    set(CMAKE_EXE_LINKER_FLAGS "${CMAKE_EXE_LINKER_FLAGS} -stdlib=libc++ -lc++abi -pthread")
    set(CMAKE_SHARED_LINKER_FLAGS "${CMAKE_SHARED_LINKER_FLAGS} -stdlib=libc++ -lc++abi -pthread")
endif ()

message("CMAKE_VERSION: ${CMAKE_VERSION}")
message("CMAKE_C_COMPILER: ${CMAKE_C_COMPILER}")
message("CMAKE_CXX_COMPILER: ${CMAKE_CXX_COMPILER}")
message("CMAKE_CXX_COMPILER_ID: ${CMAKE_CXX_COMPILER_ID}")
message("CMAKE_CXX_STANDARD: ${CMAKE_CXX_STANDARD}")
message("CMAKE_INSTALL_PREFIX: ${CMAKE_INSTALL_PREFIX}")
message("CMAKE_SOURCE_DIR: ${CMAKE_SOURCE_DIR}")
message("CMAKE_CURRENT_SOURCE_DIR: ${CMAKE_CURRENT_SOURCE_DIR}")

include(cmake/bazel.cmake)

project("test")

#set(CUSTOM_CMAKE_SOURCE_DIR "${CMAKE_SOURCE_DIR}/src")
#set(GLOBAL_INCLUDE_DIRS "/opt/homebrew/include")
#set(GLOBAL_LINK_DIRS "/opt/homebrew/lib")

cc_binary(
        NAME main
        SRCS main.cpp
)
`,
		},
		{
			Name:       "main.cpp",
			IsTemplate: false,
			Content: `#include <iostream>

int main() {
  std::cout << "hello world"
            << "\n";
}`,
		},
		{
			Name:       "cmake/bazel.cmake",
			IsTemplate: false,
			Content: `# 参考自 https://github.com/google/ruy/tree/master/cmake

# CUSTOM_CMAKE_SOURCE_DIR 可以覆盖 CMAKE_SOURCE_DIR值！
# GLOBAL_INCLUDE_DIRS 可以配置全局 include dirs
# GLOBAL_LINK_DIRS 可以配置全局 include dirs

include(CMakeParseArguments)

# cc_binary()
#
# CMake function to imitate Bazel's cc_binary rule.
function(cc_binary)
    cmake_parse_arguments(
            _RULE
            ""
            "NAME"
            "SRCS;COPTS;LINKOPTS;DEPS;TAGS;INCLUDE_DIRS;LINK_DIRS"
            ${ARGN}
    )

    if (DEFINED CUSTOM_CMAKE_SOURCE_DIR)
        list(APPEND _RULE_INCLUDE_DIRS ${CUSTOM_CMAKE_SOURCE_DIR})
    else ()
        list(APPEND _RULE_INCLUDE_DIRS ${CMAKE_SOURCE_DIR})
    endif ()
    if (DEFINED GLOBAL_INCLUDE_DIRS)
        list(APPEND _RULE_INCLUDE_DIRS ${GLOBAL_INCLUDE_DIRS})
    endif ()

    if (DEFINED GLOBAL_LINK_DIRS)
        list(APPEND _RULE_LINK_DIRS ${GLOBAL_LINK_DIRS})
    endif ()

    message("# pwd: ${CMAKE_CURRENT_LIST_DIR}
# root: ${CMAKE_SOURCE_DIR}
cc_binary(
    name: ${_RULE_NAME},
    srcs: ${_RULE_SRCS},
    deps: ${_RULE_DEPS},
    includes: ${_RULE_INCLUDE_DIRS},
    linkopts: -L${_RULE_LINK_DIRS}
    )")
    add_executable(${_RULE_NAME} "")
    target_sources(${_RULE_NAME}
            PRIVATE
            ${_RULE_SRCS}
    )
    set_target_properties(${_RULE_NAME} PROPERTIES OUTPUT_NAME "${_RULE_NAME}")
    target_include_directories(${_RULE_NAME}
            PUBLIC
            "$<BUILD_INTERFACE:${_RULE_INCLUDE_DIRS}>"
            "$<INSTALL_INTERFACE:${CMAKE_INSTALL_INCLUDEDIR}>"
    )
    target_compile_options(${_RULE_NAME}
            PRIVATE
            ${_RULE_COPTS}
    )
    target_link_directories(${_RULE_NAME}
            PUBLIC
            "$<BUILD_INTERFACE:${_RULE_LINK_DIRS}>"
            "$<INSTALL_INTERFACE:${CMAKE_INSTALL_LIBDIR}>"
    )
    target_link_options(${_RULE_NAME}
            PRIVATE
            ${_RULE_LINKOPTS}
    )
    target_link_libraries(${_RULE_NAME}
            PUBLIC
            ${_RULE_DEPS}
    )
endfunction()

# CMake function to imitate Bazel's cc_library rule.
function(cc_library)
    cmake_parse_arguments(
            _RULE
            "TESTONLY"
            "NAME;ALIAS"
            "HDRS;SRCS;COPTS;DEFINES;LINKOPTS;DEPS;INCLUDE_DIRS;LINK_DIRS"
            ${ARGN}
    )

    # Check if this is a header-only library.
    if ("${_RULE_SRCS}" STREQUAL "")
        set(_RULE_IS_INTERFACE 1)
    else ()
        set(_RULE_IS_INTERFACE 0)
    endif ()

    if (DEFINED CUSTOM_CMAKE_SOURCE_DIR)
        list(APPEND _RULE_INCLUDE_DIRS ${CUSTOM_CMAKE_SOURCE_DIR})
    else ()
        list(APPEND _RULE_INCLUDE_DIRS ${CMAKE_SOURCE_DIR})
    endif ()
    if (DEFINED GLOBAL_INCLUDE_DIRS)
        list(APPEND _RULE_INCLUDE_DIRS ${GLOBAL_INCLUDE_DIRS})
    endif ()


    if (DEFINED GLOBAL_LINK_DIRS)
        list(APPEND _RULE_LINK_DIRS ${GLOBAL_LINK_DIRS})
    endif ()

    message("# pwd: ${CMAKE_CURRENT_LIST_DIR}
# root: ${CMAKE_SOURCE_DIR}
cc_library(
    name: ${_RULE_NAME},
    alias: ${_RULE_ALIAS},
    srcs: ${_RULE_SRCS},
    hdrs: ${_RULE_HDRS},
    deps: ${_RULE_DEPS},
    includes: ${_RULE_INCLUDE_DIRS},
    linkopts: -L${_RULE_LINK_DIRS}
    )")

    if (_RULE_IS_INTERFACE)
        # Generating a header-only library.
        add_library(${_RULE_NAME} INTERFACE)
        set_target_properties(${_RULE_NAME} PROPERTIES PUBLIC_HEADER "${_RULE_HDRS}")
        target_include_directories(${_RULE_NAME}
                INTERFACE
                "$<BUILD_INTERFACE:${_RULE_INCLUDE_DIRS}>" # 指定include路径
                "$<INSTALL_INTERFACE:${CMAKE_INSTALL_INCLUDEDIR}>" # 指定install路径
        )
        target_link_libraries(${_RULE_NAME}
                INTERFACE
                ${_RULE_DEPS}
                ${_RULE_LINKOPTS}
        )
        target_compile_definitions(${_RULE_NAME}
                INTERFACE
                ${_RULE_DEFINES}
        )
    else ()
        # Generating a static binary library.
        add_library(${_RULE_NAME} STATIC ${_RULE_SRCS} ${_RULE_HDRS})
        set_target_properties(${_RULE_NAME} PROPERTIES PUBLIC_HEADER "${_RULE_HDRS}")
        target_include_directories(${_RULE_NAME}
                PUBLIC
                "$<BUILD_INTERFACE:${_RULE_INCLUDE_DIRS}>" # CMAKE_SOURCE_DIR/PROJECT_SOURCE_DIR
                "$<INSTALL_INTERFACE:${CMAKE_INSTALL_INCLUDEDIR}>"
        )
        target_compile_options(${_RULE_NAME}
                PRIVATE
                ${_RULE_COPTS}
        )
        target_link_directories(${_RULE_NAME}
                PUBLIC
                "$<BUILD_INTERFACE:${_RULE_LINK_DIRS}>"
                "$<INSTALL_INTERFACE:${CMAKE_INSTALL_LIBDIR}>"
        )
        target_link_libraries(${_RULE_NAME}
                PUBLIC
                ${_RULE_DEPS}
        )
        target_link_options(${_RULE_NAME}
                PRIVATE
                ${_RULE_LINKOPTS}
        )
        target_compile_definitions(${_RULE_NAME}
                PUBLIC
                ${_RULE_DEFINES}
        )
    endif ()

    if (NOT _RULE_TESTONLY)
        # install dir
        if (DEFINED CUSTOM_CMAKE_SOURCE_DIR)
            file(RELATIVE_PATH _RULE_SUBDIR ${CUSTOM_CMAKE_SOURCE_DIR} ${CMAKE_CURRENT_LIST_DIR})
        else ()
            file(RELATIVE_PATH _RULE_SUBDIR ${CMAKE_SOURCE_DIR} ${CMAKE_CURRENT_LIST_DIR})
        endif ()
        install(
                TARGETS ${_RULE_NAME}
                LIBRARY DESTINATION ${CMAKE_INSTALL_LIBDIR}
                PUBLIC_HEADER DESTINATION ${CMAKE_INSTALL_INCLUDEDIR}/${CMAKE_PROJECT_NAME}/${_RULE_SUBDIR}
        )
    endif ()
    if (NOT "${_RULE_ALIAS}" STREQUAL "")
        add_library(${_RULE_ALIAS} ALIAS ${_RULE_NAME})
    endif ()
endfunction()

# CMake function to imitate Bazel's cc_test rule.
function(cc_test)
    cmake_parse_arguments(
            _RULE
            ""
            "NAME"
            "SRCS;COPTS;LINKOPTS;DEPS;TAGS;INCLUDE_DIRS;LINK_DIRS"
            ${ARGN}
    )
    if (DEFINED CUSTOM_CMAKE_SOURCE_DIR)
        list(APPEND _RULE_INCLUDE_DIRS ${CUSTOM_CMAKE_SOURCE_DIR})
    else ()
        list(APPEND _RULE_INCLUDE_DIRS ${CMAKE_SOURCE_DIR})
    endif ()
    if (DEFINED GLOBAL_INCLUDE_DIRS)
        list(APPEND _RULE_INCLUDE_DIRS ${GLOBAL_INCLUDE_DIRS})
    endif ()

    message("# pwd: ${CMAKE_CURRENT_LIST_DIR}
# root: ${CMAKE_SOURCE_DIR}
cc_test(
    name: ${_RULE_NAME},
    srcs: ${_RULE_SRCS},
    deps: ${_RULE_DEPS},
    includes: ${CPP_COMMON_INCLUDE_DIRS},
    linkopts: -L${_RULE_LINK_DIRS}
    )")
    add_executable(${_RULE_NAME} "")
    target_sources(${_RULE_NAME}
            PRIVATE
            ${_RULE_SRCS}
    )
    set_target_properties(${_RULE_NAME} PROPERTIES OUTPUT_NAME "${_RULE_NAME}")
    target_include_directories(${_RULE_NAME}
            PUBLIC
            "$<BUILD_INTERFACE:${_RULE_INCLUDE_DIRS}>"
            "$<INSTALL_INTERFACE:${CMAKE_INSTALL_INCLUDEDIR}>"
    )
    target_compile_options(${_RULE_NAME}
            PRIVATE
            ${_RULE_COPTS}
    )
    target_link_options(${_RULE_NAME}
            PRIVATE
            ${_RULE_LINKOPTS}
    )
    target_link_directories(${_RULE_NAME}
            PUBLIC
            "$<BUILD_INTERFACE:${_RULE_LINK_DIRS}>"
            "$<INSTALL_INTERFACE:${CMAKE_INSTALL_LIBDIR}>"
    )
    target_link_libraries(${_RULE_NAME}
            PUBLIC
            ${_RULE_DEPS}
    )
    add_test(
            NAME
            ${_RULE_NAME}
            COMMAND
            "$<TARGET_FILE:${_RULE_NAME}>"
    )
    if (_RULE_TAGS)
        set_property(TEST ${_NAME} PROPERTY LABELS ${_RULE_TAGS})
    endif ()
endfunction()`,
		},
	}
}
