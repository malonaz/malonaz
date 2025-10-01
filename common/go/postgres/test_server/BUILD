go_library(
    name = "test_server",
    srcs = [
        "server.go",
        "utils.go",
    ],
    test_only = True,
    visibility = ["//..."],
    deps = [
        "//common/go/binary",
        "//common/go/logging",
        "//common/go/postgres",
        "//common/go/postgres/migrator",
        "//common/go/postgres/migrator/migrations",
    ],
)
