<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [github.com/golang/snappy](#githubcomgolangsnappy)
  - [Snappy的基本原理](#snappy%E7%9A%84%E5%9F%BA%E6%9C%AC%E5%8E%9F%E7%90%86)
  - [优势总结](#%E4%BC%98%E5%8A%BF%E6%80%BB%E7%BB%93)
  - [第三方应用-->prometheus](#%E7%AC%AC%E4%B8%89%E6%96%B9%E5%BA%94%E7%94%A8--prometheus)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# github.com/golang/snappy


Snappy是由Google开发的一种压缩/解压缩库。其设计目标是提供高效的数据压缩速度，而不是追求最大的压缩比。
这一特点使得Snappy在需要快速处理数据的场景中，如网络通信和大规模数据处理，表现尤为出色


## Snappy的基本原理
Snappy采用了一种基于字典压缩和LZ77算法的混合策略，通过识别和重用数据中的重复模式，实现高效压缩。
同时，Snappy优化了数据的处理流程，减少了计算开销，从而在保持较低压缩比的同时，显著提升了压缩和解压缩的速度

## 优势总结

- 快速：压缩速度大概在250MB/秒及更快的速度进行压缩。
- 稳定：在过去的几年中，Snappy在Google的生产环境中压缩并解压缩了数P字节（petabytes）的数据。Snappy位流格式是稳定的，不会在版本之间发生变化
- 健壮性：Snappy解压缩器设计为不会因遇到损坏或恶意输入而崩溃




## 第三方应用-->prometheus 


记录日志压缩过程
```go
// https://github.com/prometheus/prometheus/blob/775d90d5f87a64f3594c8b911ab2bc65a04f80c1/tsdb/wlog/wlog.go
func (w *WL) log(rec []byte, final bool) error {
	// When the last page flush failed the page will remain full.
	// When the page is full, need to flush it before trying to add more records to it.
	if w.page.full() {
		if err := w.flushPage(true); err != nil {
			return err
		}
	}

	// Compress the record before calculating if a new segment is needed.
	compressed := false
	if w.compress &&
		len(rec) > 0 &&
		// If MaxEncodedLen is less than 0 the record is too large to be compressed.
		snappy.MaxEncodedLen(len(rec)) >= 0 {
		// The snappy library uses `len` to calculate if we need a new buffer.
		// In order to allocate as few buffers as possible make the length
		// equal to the capacity.
		w.snappyBuf = w.snappyBuf[:cap(w.snappyBuf)]
		w.snappyBuf = snappy.Encode(w.snappyBuf, rec)
		if len(w.snappyBuf) < len(rec) {
			rec = w.snappyBuf
			compressed = true
		}
	}
	/// ...
}	
```

读取解压过程

```go
// https://github.com/prometheus/prometheus/blob/775d90d5f87a64f3594c8b911ab2bc65a04f80c1/tsdb/wlog/reader.go

func (r *Reader) next() (err error) {
	// We have to use r.buf since allocating byte arrays here fails escape
	// analysis and ends up on the heap, even though it seemingly should not.
	hdr := r.buf[:recordHeaderSize]
	buf := r.buf[recordHeaderSize:]

	r.rec = r.rec[:0]
	r.snappyBuf = r.snappyBuf[:0]

	i := 0
	for {
		if _, err = io.ReadFull(r.rdr, hdr[:1]); err != nil {
			return errors.Wrap(err, "read first header byte")
		}
		r.total++
		r.curRecTyp = recTypeFromHeader(hdr[0])
		compressed := hdr[0]&snappyMask != 0

        // ...

		if compressed {
			r.snappyBuf = append(r.snappyBuf, buf[:length]...)
		} else {
			r.rec = append(r.rec, buf[:length]...)
		}
        // ...
		if r.curRecTyp == recLast || r.curRecTyp == recFull {
			if compressed && len(r.snappyBuf) > 0 {
				// The snappy library uses `len` to calculate if we need a new buffer.
				// In order to allocate as few buffers as possible make the length
				// equal to the capacity.
				r.rec = r.rec[:cap(r.rec)]
				r.rec, err = snappy.Decode(r.rec, r.snappyBuf)
				return err
			}
			return nil
		}
		// ..
		i++
	}
}

```

## 参考
- https://segmentfault.com/a/1190000045387037