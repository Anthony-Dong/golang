# 参考自 https://github.com/google/ruy/tree/master/cmake

include(CMakeParseArguments)

# cc_test()
#
# CMake function to imitate Bazel's cc_test rule.
# CUSTOM_PROJECT_SOURCE_DIR: 业务自定义项目根路径
function(cc_test)
    cmake_parse_arguments(
            _RULE
            ""
            "NAME"
            "SRCS;COPTS;LINKOPTS;DEPS;TAGS;INCLUDE_DIRS"
            ${ARGN}
    )
    list(APPEND _RULE_DEPS "GTest" "GTest::gtest_main")

    if (DEFINED CUSTOM_PROJECT_SOURCE_DIR)
        list(APPEND _RULE_INCLUDE_DIRS ${CUSTOM_PROJECT_SOURCE_DIR})
    endif ()

    message("cc_test(
    # pwd: ${CMAKE_CURRENT_LIST_DIR}
    name: ${_RULE_NAME},
    srcs: ${_RULE_SRCS},
    deps: ${_RULE_DEPS},
    includes: ${_RULE_INCLUDE_DIRS},
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
endfunction()