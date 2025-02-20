load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "gormtest",
    srcs = ["gormtest.go"],
    importpath = "github.com/ngtrvu/zen-go/gormtest",
    visibility = ["//visibility:public"],
    deps = [
        "//frameworks/common/go/gorm",
        "//frameworks/common/go/gorm/migrator",
        "//frameworks/common/go/log",
        "//frameworks/common/go/utils",
        "@com_github_golang_migrate_migrate//source/go_bindata",
        "@com_github_google_uuid//:uuid",
        "@com_github_stretchr_testify//require",
        "@com_github_stretchr_testify//suite",
        "@io_gorm_driver_postgres//:postgres",
        "@io_gorm_gorm//:gorm",
    ],
)
