# Gopstuil
屏蔽了各个系统之间的差异，帮助我们方便地获取各种系统和硬件信息

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