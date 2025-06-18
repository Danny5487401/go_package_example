<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [ZSTANDARD（ZSTD）](#zstandardzstd)
  - [第三方实现](#%E7%AC%AC%E4%B8%89%E6%96%B9%E5%AE%9E%E7%8E%B0)
    - [gozstd-->基于cgo的封装](#gozstd--%E5%9F%BA%E4%BA%8Ecgo%E7%9A%84%E5%B0%81%E8%A3%85)
    - [klauspost/compress/zstd, pure go的实现](#klauspostcompresszstd-pure-go%E7%9A%84%E5%AE%9E%E7%8E%B0)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->


# ZSTANDARD（ZSTD）



## 第三方实现

### gozstd-->基于cgo的封装


victoriaMetrics 中使用


```go
// https://github.com/VictoriaMetrics/VictoriaMetrics/blob/cf23dc6480f77b79de500f145135a8f7be0ac065/lib/encoding/zstd/zstd_cgo.go
//go:build cgo

package zstd

import (
	"github.com/valyala/gozstd"
)

// Decompress appends decompressed src to dst and returns the result.
func Decompress(dst, src []byte) ([]byte, error) {
	return gozstd.Decompress(dst, src)
}

// CompressLevel appends compressed src to dst and returns the result.
//
// The given compressionLevel is used for the compression.
func CompressLevel(dst, src []byte, compressionLevel int) []byte {
	return gozstd.CompressLevel(dst, src, compressionLevel)
}

```

### klauspost/compress/zstd, pure go的实现





