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
        "//proto/exp:api",
        "//third_party/go:github.com__malonaz__core__go__pbutil",
        "//third_party/go:github.com__malonaz__core__go__postgres",
    ],
)
