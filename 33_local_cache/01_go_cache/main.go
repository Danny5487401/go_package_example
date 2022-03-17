package main

import (
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
)

type MyStruct struct {
	Name string
}

func main() {
	// 设置超时时间和清理时间
	c := cache.New(5*time.Minute, 10*time.Minute)

	// 设置缓存值并带上过期时间
	c.Set("foo", "bar", cache.DefaultExpiration)

	// 设置没有过期时间的KEY，这个KEY不会被自动清除，想清除使用：c.Delete("baz")
	c.Set("baz", 42, cache.NoExpiration)

	var foo interface{}
	var found bool
	// 获取值
	foo, found = c.Get("foo")
	if found {
		fmt.Println(foo)
	}

	var foos string
	// 获取值， 并断言
	if x, found := c.Get("foo"); found {
		foos = x.(string)
		fmt.Println(foos)
	}
	// 对结构体指针进行操作
	var my *MyStruct
	c.Set("foo", &MyStruct{Name: "Danny"}, cache.DefaultExpiration)
	if x, found := c.Get("foo"); found {
		my = x.(*MyStruct)
		// ...
		fmt.Println(my)
	}

}
