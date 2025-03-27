package main

import (
	"encoding/json"
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

func main() {
	// cpu
	fmt.Println(cpu.Counts(true))
	info1, _ := cpu.Info()
	for _, info := range info1 {
		data, _ := json.MarshalIndent(info, "", " ")
		fmt.Print(string(data))
	}

	// 磁盘信息
	info2, _ := disk.Partitions(false)
	for _, info := range info2 {
		data, _ := json.MarshalIndent(info, "", "  ")
		fmt.Println(string(data))
	}

	// 内存信息
	v, _ := mem.VirtualMemory()

	fmt.Printf("Total: %v, Available: %v, UsedPercent:%f%%\n", v.Total, v.Available, v.UsedPercent)

}
