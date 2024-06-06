cc_library(
    name = "utils",
    srcs = [
        "times.cpp",
        "utils.cpp",
    ],
    hdrs = glob(include=["*.h"]),
)


cc_binary(
    name = "main",
    srcs = [
        "main.cpp",
    ],
    deps = [
        ":utils",
    ],
)