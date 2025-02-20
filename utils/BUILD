load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "utils",
    srcs = [
        "env.go",
        "image.go",
        "numbers.go",
        "slices.go",
        "strings.go",
        "struct.go",
        "time.go",
    ],
    importpath = "github.com/ngtrvu/zen-go/utils",
    visibility = ["//visibility:public"],
    deps = [
        "@org_golang_x_text//language",
        "@org_golang_x_text//message",
    ],
)

go_test(
    name = "utils_test",
    size = "small",
    srcs = [
        "numbers_test.go",
        "slices_test.go",
        "strings_test.go",
        "struct_test.go",
        "time_test.go",
    ],
    deps = [
        ":utils",
        "@com_github_stretchr_testify//assert",
        "@com_github_stretchr_testify//require",
    ],
)
