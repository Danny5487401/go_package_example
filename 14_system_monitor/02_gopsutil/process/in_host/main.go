package main

import (
	"fmt"
	"github.com/shirou/gopsutil/process"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

// 下面获取进程资源占比的方法只有在虚拟机和物理机环境下才能准确。
// 类似Docker这样的Linux容器是靠着Linux的Namespace和Cgroups技术实现的进程隔离和资源限制，是不行的。

func main() {
	// 创建进程对象
	p, _ := process.NewProcess(int32(os.Getpid()))

	fmt.Println("进程对象", p)

	// 进程的CPU使用率
	cpuPercent, _ := p.Percent(time.Second)
	fmt.Println("进程的CPU使用率", cpuPercent)

	cp := cpuPercent / float64(runtime.NumCPU())
	fmt.Println("单个核心的比例使用率", cp)

	// 获取进程占用内存的比例
	mp, _ := p.MemoryPercent()
	fmt.Println("获取进程占用内存的比例", mp)
	// 创建的线程数
	threadCount := pprof.Lookup("threadcreate").Count()
	fmt.Println("创建的线程数", threadCount)
	// Goroutine数
	gNum := runtime.NumGoroutine()
	fmt.Println("Goroutine数", gNum)

}
