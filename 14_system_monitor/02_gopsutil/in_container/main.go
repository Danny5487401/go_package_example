package main

import (
	"fmt"
	"github.com/shirou/gopsutil/process"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

const unifiedMountPoint = "/sys/fs/cgroup"

func main() {
	// 创建进程对象
	p, _ := process.NewProcess(int32(os.Getpid()))
	cpuPeriod, err := readUint(unifiedMountPoint + "/cpu/cpu.cfs_period_us")
	if err != nil {
		fmt.Println("获取路径信息失败", err.Error())
		return
	}

	cpuQuota, err := readUint(unifiedMountPoint + "/cpu/cpu.cfs_quota_us")
	if err != nil {
		fmt.Println("获取路径信息失败", err.Error())
		return
	}

	cpuNum := float64(cpuQuota) / float64(cpuPeriod)

	cpuPercent, err := p.Percent(time.Second)
	// cp := cpuPercent / float64(runtime.NumCPU())
	// 调整为
	cp := cpuPercent / cpuNum
	fmt.Println("单个核心的比例使用率", cp)

	// 容器的能使用的最大内存数，自然就是在memory.limit_in_bytes
	memLimit, err := readUint(unifiedMountPoint + "/memory/memory.limit_in_bytes")
	memInfo, err := p.MemoryInfo()
	//RSS叫常驻内存，是在RAM里分配给进程，允许进程访问的内存量
	// 进程占用内存的比例
	mp := memInfo.RSS * 100 / memLimit
	fmt.Println("获取进程占用内存的比例", mp)

}

func readUint(path string) (uint64, error) {
	v, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, err
	}
	return parseUint(strings.TrimSpace(string(v)), 10, 64)
}

func parseUint(s string, base, bitSize int) (uint64, error) {
	v, err := strconv.ParseUint(s, base, bitSize)
	if err != nil {
		intValue, intErr := strconv.ParseInt(s, base, bitSize)
		// 1. Handle negative values greater than MinInt64 (and)
		// 2. Handle negative values lesser than MinInt64
		if intErr == nil && intValue < 0 {
			return 0, nil
		} else if intErr != nil &&
			intErr.(*strconv.NumError).Err == strconv.ErrRange &&
			intValue < 0 {
			return 0, nil
		}
		return 0, err
	}
	return v, nil
}
