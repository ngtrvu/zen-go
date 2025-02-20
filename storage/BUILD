load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "storage",
    srcs = ["cloud_storage.go"],
    importpath = "github.com/ngtrvu/zen-go/storage",
    visibility = ["//visibility:public"],
    deps = [
        "//frameworks/common/go/log",
        "@com_google_cloud_go_storage//:storage",
    ],
)

go_test(
    name = "storage_test",
    size = "small",
    srcs = ["cloud_storage_test.go"],
    embed = [
        ":storage",
    ],
    deps = [
        "@org_uber_go_mock//gomock",
    ],
)
