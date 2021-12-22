# Gopstuil
gopsutil是 Python 工具库psutil 的 Golang 移植版，可以帮助我们方便地获取各种系统和硬件信息。gopsutil为我们屏蔽了各个系统之间的差异，具有非常强悍的可移植性。
有了gopsutil，我们不再需要针对不同的系统使用syscall调用对应的系统方法。更棒的是gopsutil的实现中没有任何cgo的代码，使得交叉编译成为可能。

## 模块
* cpu：系统CPU 相关模块；
* disk：系统磁盘相关模块；
* docker：docker 相关模块；
* mem：内存相关模块；
* net：网络相关；
* process：进程相关模块；
* winservices：Windows 服务相关模块


## 容器环境下获取进程指标
在Linux中，Cgroups给用户暴露出来的操作接口是文件系统，它以文件和目录的方式组织在操作系统的/sys/fs/cgroup路径下，在 /sys/fs/cgroup下面有很多诸cpuset、cpu、 memory这样的子目录，每个子目录都代表系统当前可以被Cgroups进行限制的资源种类

针对我们监控Go进程内存和CPU指标的需求，我们只要知道cpu.cfs_period_us、cpu.cfs_quota_us 和memory.limit_in_bytes 就行。
前两个参数需要组合使用，可以用来限制进程在长度为cfs_period的一段时间内，只能被分配到总量为cfs_quota的CPU时间， 可以简单的理解为容器能使用的核心数 = cfs_quota / cfs_period。


## 使用

### CPU
我们知道 CPU 的核数有两种，一种是物理核数，一种是逻辑核数。物理核数就是主板上实际有多少个 CPU，一个物理 CPU 上可以有多个核心，这些核心被称为逻辑核。gopsutil中 CPU 相关功能在cpu子包中，cpu子包提供了获取物理和逻辑核数、CPU 使用率的接口：
- Counts(logical bool)：传入false，返回物理核数，传入true，返回逻辑核数；
- Percent(interval time.Duration, percpu bool)：表示获取interval时间间隔内的 CPU 使用率，percpu为false时，获取总的 CPU 使用率，percpu为true时，分别获取每个 CPU 的使用率，返回一个[]float64类型的值.
### 磁盘
子包disk用于获取磁盘信息。disk可获取 IO 统计、分区和使用率信息。
#### 1. IO
调用disk.IOCounters()函数，返回的 IO 统计信息用map[string]IOCountersStat类型表示。每个分区一个结构，键为分区名，值为统计信息。这里摘取统计结构的部分字段，主要有读写的次数、字节数和时间
```go
// src/github.com/shirou/gopsutil/disk/disk.go
type IOCountersStat struct {
  ReadCount        uint64 `json:"readCount"`
  MergedReadCount  uint64 `json:"mergedReadCount"`
  WriteCount       uint64 `json:"writeCount"`
  MergedWriteCount uint64 `json:"mergedWriteCount"`
  ReadBytes        uint64 `json:"readBytes"`
  WriteBytes       uint64 `json:"writeBytes"`
  ReadTime         uint64 `json:"readTime"`
  WriteTime        uint64 `json:"writeTime"`
  // ...
}
```

#### 2. 分区
调用disk.PartitionStat(all bool)函数，返回分区信息。如果all = false，只返回实际的物理分区（包括硬盘、CD-ROM、USB），忽略其它的虚拟分区。
如果all = true则返回所有的分区。返回类型为[]PartitionStat，每个分区对应一个PartitionStat结构：
```go
// src/github.com/shirou/gopsutil/disk/
type PartitionStat struct {
  Device     string `json:"device"`  //分区标识，在 Windows 上即为C:这类格式；
  Mountpoint string `json:"mountpoint"` //挂载点，即该分区的文件路径起始位置
  Fstype     string `json:"fstype"` // 文件系统类型，Windows 常用的有 FAT、NTFS 等，Linux 有 ext、ext2、ext3等
  Opts       string `json:"opts"`
}
```

#### 3. 使用率
调用disk.Usage(path string)即可获得路径path所在磁盘的使用情况，返回一个UsageStat结构：
```go
// src/github.com/shirou/gopsutil/disk.go
type UsageStat struct {
  Path              string  `json:"path"`
  Fstype            string  `json:"fstype"`
  Total             uint64  `json:"total"`  //该分区总容量；
  Free              uint64  `json:"free"`
  Used              uint64  `json:"used"`
  UsedPercent       float64 `json:"usedPercent"`
  InodesTotal       uint64  `json:"inodesTotal"`
  InodesUsed        uint64  `json:"inodesUsed"`
  InodesFree        uint64  `json:"inodesFree"`
  InodesUsedPercent float64 `json:"inodesUsedPercent"`
}
```

### 主机 
子包host可以获取主机相关信息，如开机时间、内核版本号、平台信息等等。

### 内存
使用mem.VirtualMemory()来获取内存信息。该函数返回的只是物理内存信息。我们还可以使用mem.SwapMemory()获取交换内存的信息，信息存储在结构SwapMemoryStat中
```go
// src/github.com/shirou/gopsutil/mem/
type SwapMemoryStat struct {
  Total       uint64  `json:"total"`
  Used        uint64  `json:"used"`
  Free        uint64  `json:"free"`
  UsedPercent float64 `json:"usedPercent"`
  Sin         uint64  `json:"sin"`
  Sout        uint64  `json:"sout"`
  PgIn        uint64  `json:"pgin"`  //载入页数
  PgOut       uint64  `json:"pgout"`  //淘汰页数
  PgFault     uint64  `json:"pgfault"` //缺页错误数
}
```
交换内存是以页为单位的，如果出现缺页错误(page fault)，操作系统会将磁盘中的某些页载入内存，同时会根据特定的机制淘汰一些内存中的页。PgIn表征载入页数，PgOut淘汰页数，PgFault缺页错误数。