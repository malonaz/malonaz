subinclude("//build_defs:utils")

go_copy(
    name = "pb_model.go",
    src = "//proto/exp",
)

go_library(
    name = "model",
    srcs = [
        ":pb_model.go",
    ],
    visibility = ["//..."],
    deps = [
    ],
)
