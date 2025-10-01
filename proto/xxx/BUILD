subinclude("//build_defs:utils")

go_copy(
    name = "pb_model.go",
    src = "//proto/exp:model",
)

go_library(
    name = "model",
    srcs = [
        ":pb_model.go",
    ],
    visibility = ["//..."],
    deps = [
        "//common/go/pbutil",
        "//common/go/postgres",
        "//proto/exp:api",
    ],
)
