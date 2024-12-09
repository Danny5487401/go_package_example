package main

import (
	"fmt"
	"github.com/golang/snappy"
)

func main() {
	// 原始数据
	data := []byte("Hello, Snappy Compression!")

	fmt.Printf("data before compression len:%d\n", len(data))
	// 压缩数据
	var buf []byte // 缓冲区
	compressed := snappy.Encode(buf, data)
	fmt.Println("压缩后的数据:", compressed)
	fmt.Printf("data after compression len:%d\n", len(compressed)) //有时候压缩后会比压缩前字节数变大。这是和原字符串有关系

	// 解压缩数据
	decompressed, err := snappy.Decode(nil, compressed)
	if err != nil {
		panic(err)
	}
	fmt.Println("解压缩后的数据:", string(decompressed))
}
