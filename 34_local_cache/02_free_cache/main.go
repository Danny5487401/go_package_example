package main

import (
	"fmt"
	"runtime/debug"

	"github.com/coocood/freecache"
)

func main() {
	// 缓存大小，100M
	cacheSize := 100 * 1024 * 1024
	cache := freecache.NewCache(cacheSize)
	debug.SetGCPercent(20)
	key := []byte("abc")
	val := []byte("def")

	// 设置KEY
	expire := 60 // expire in 60 seconds
	cache.Set(key, val, expire)
	// 获取key
	got, err := cache.Get(key)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%s\n", got)
	}
	fmt.Println("entry count ", cache.EntryCount())
	affected := cache.Del(key)
	fmt.Println("deleted key ", affected)
	fmt.Println("entry count ", cache.EntryCount())
}
