github.com/qiniu/api.v7 (Qiniu Go SDK v7.x)
===============

[![LICENSE](https://img.shields.io/github/license/qiniu/api.v7.svg)](https://github.com/qiniu/api.v7/blob/master/LICENSE)
[![Build Status](https://travis-ci.org/qiniu/api.v7.svg?branch=master)](https://travis-ci.org/qiniu/api.v7)
[![Go Report Card](https://goreportcard.com/badge/github.com/qiniu/api.v7)](https://goreportcard.com/report/github.com/qiniu/api.v7)
[![GitHub release](https://img.shields.io/github/v/tag/qiniu/api.v7.svg?label=release)](https://github.com/qiniu/api.v7/releases)
[![codecov](https://codecov.io/gh/qiniu/api.v7/branch/master/graph/badge.svg)](https://codecov.io/gh/qiniu/api.v7)
[![GoDoc](https://godoc.org/github.com/qiniu/api.v7?status.svg)](https://godoc.org/github.com/qiniu/api.v7)

[![Qiniu Logo](http://open.qiniudn.com/logo.png)](http://qiniu.com/)

# 【迁移】

Qiniu Go SDK 代码库已经被迁移到 `github.com/qiniu/go-sdk`，当前代码库 `github.com/qiniu/api.v7` 将不再更新，请尽快将您项目中依赖的 `github.com/qiniu/api.v7/v7` 直接替换为 `github.com/qiniu/go-sdk/v7`，该替换不会造成不兼容的问题。

## 过时警告

当前代码库 `github.com/qiniu/api.v7` 在初始化时会输出过时警告，如果您暂时无法迁移又不希望看到该警告，可以在执行应用程序前设置环境变量 `SUPPRESS_DEPRECATION_WARNING=1` 来屏蔽该警告。

# 下载

## 使用 Go mod【推荐】

在您的项目中的 `go.mod` 文件内添加这行代码

```
require github.com/qiniu/api.v7/v7 v7.8.2
```

并且在项目中使用 `"github.com/qiniu/api.v7/v7"` 引用 Qiniu Go SDK。

例如

```go
import (
    "github.com/qiniu/api.v7/v7/auth"
    "github.com/qiniu/api.v7/v7/storage"
)
```

## 不使用 Go mod【不推荐，且只能获取 v7.2.5 及其以下版本】

```bash
go get -u github.com/qiniu/api.v7
```

# go版本需求

需要 go1.10 或者 1.10 以上

#  文档

[七牛SDK文档站](https://developer.qiniu.com/kodo/sdk/1238/go) 或者 [项目WIKI](https://github.com/qiniu/api.v7/wiki)

# 示例

[参考代码](https://github.com/qiniu/api.v7/tree/master/examples)
