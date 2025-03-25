package(default_visibility = ["PUBLIC"])

subinclude("//build_defs:remote")

GOOGLEAPI_REVISION = "5e5d500bc4f8c4d5591cd1fda817556b3fe663a7"

github_proto_srcs(
    name = "v1alpha1_srcs",
    dir = "google/api/expr/v1alpha1",
    files = [
        "checked.proto",
        "eval.proto",
        "explain.proto",
        "syntax.proto",
        "value.proto",
    ],
    repo = "googleapis/googleapis",
    revision = GOOGLEAPI_REVISION,
)

github_proto_srcs(
    name = "status_srcs",
    dir = "google/rpc",
    files = ["status.proto"],
    repo = "googleapis/googleapis",
    revision = GOOGLEAPI_REVISION,
)

proto_library(
    name = "status",
    srcs = [":status_srcs"],
    protoc_flags = ["-I" + package_name()],
)

proto_library(
    name = "v1alpha1",
    srcs = [":v1alpha1_srcs"],
    protoc_flags = ["-I" + package_name()],
    deps = [":status"],
)
