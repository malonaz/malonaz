subinclude("//build_defs/codegen:go_proto_templates")

proto_library(
    name = "api",
    srcs = ["api.proto"],
    visibility = [
        "PUBLIC",
    ],
    deps = [
        "//proto/codegen:model",
    ],
)

go_proto_template_library(
    name = "model",
    srcs = ["api.proto"],
    template = "//build_defs/codegen/templates:model",
    visibility = ["PUBLIC"],
    deps = [
        "//proto/codegen:model",
    ],
)
