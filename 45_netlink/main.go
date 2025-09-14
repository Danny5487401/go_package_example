package main

import (
	"fmt"
	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

func main() {
	// 创建 Netlink Handle 对象
	handle, err := netlink.NewHandle(unix.NETLINK_ROUTE)
	if err != nil {
		fmt.Println("创建 Netlink Handle 失败:", err)
		return
	}
	defer handle.Close() // 确保在程序结束时关闭 Netlink Handle

	// 获取所有路由信息
	routes, err := handle.RouteList(nil, unix.AF_INET)
	if err != nil {
		fmt.Println("获取路由列表失败:", err)
		return
	}

	fmt.Println("路由列表:")
	for _, route := range routes {
		fmt.Printf("目标网段: %s, 下一跳: %s, 接口索引: %d\n", route.Dst.String(), route.Gw.String(), route.LinkIndex)
	}
}
