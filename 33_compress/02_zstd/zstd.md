
# zstd



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



