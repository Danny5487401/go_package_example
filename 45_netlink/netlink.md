<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [netlink](#netlink)
  - [Linux用户空间与内核空间交互方式](#linux%E7%94%A8%E6%88%B7%E7%A9%BA%E9%97%B4%E4%B8%8E%E5%86%85%E6%A0%B8%E7%A9%BA%E9%97%B4%E4%BA%A4%E4%BA%92%E6%96%B9%E5%BC%8F)
  - [优点](#%E4%BC%98%E7%82%B9)
  - [github.com/vishvananda/netlink](#githubcomvishvanandanetlink)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# netlink


Netlink 是一种用于进程间通信（IPC）的机制，主要用于网络协议栈与用户空间之间的通信。它支持用户程序与内核的双向通信，常用于网络配置和管理。

## Linux用户空间与内核空间交互方式

CPU 将指令分为特权指令和非特权指令，对于那些危险的指令，只允许操作系统及其相关模块使用，普通应用程序只能使用那些不会造成灾难的指令。比如 Intel 的 CPU 将特权等级分为 4 个级别：Ring0~Ring3


当进程运行在 Ring3 级别时被称为运行在用户态，而运行在 Ring0 级别时被称为运行在内核态。

Linux应用程序与内核程序交互主要有以下几种通信方式：

1. 系统调用

2. 虚拟文件系统
   - proc文件: /proc 是一个虚拟文件系统，允许用户空间程序通过读取和写入文件的方式与内核进行交互。它提供了内核状态信息，例如进程信息、内存使用情况等。
   - sysfs文件系统
   - debugfs文件系统

3. IOCTL（Input/Output Control）: IOCTL 是一种允许用户空间程序与设备驱动程序通信的机制。它通常用于向驱动程序传递特定的命令或配置设备参数。

4. netlink通信 （ip命令通过该方式）

5. 内存映像 mmap共享内存


## 优点

- netlink使用简单，只需要在include/linux/netlink.h中增加一个新类型的 netlink 协议定义即可,(如 #define NETLINK_TEST 20然后，内核和用户态应用就可以立即通过 socket API 使用该 netlink 协议类型进行数据交换)
- netlink是一种异步通信机制，在内核与用户态应用之间传递的消息保存在socket缓存队列中，发送消息只是把消息保存在接收者的socket的接收队列，而不需要等待接收者收到消息
- 使用 netlink 的内核部分可以采用模块的方式实现，使用 netlink 的应用部分和内核部分没有编译时依赖
- netlink 支持多播，内核模块或应用可以把消息多播给一个netlink组，属于该neilink 组的任何内核模块或应用都能接收到该消息，内核事件向用户态的通知机制就使用了这一特性
- 内核可以使用 netlink 首先发起会话


Netlink协议基于BSD socket和AF_NETLINK地址簇，使用32位的端口号寻址，每个Netlink协议通常与一个或一组内核服务/组件相关联，
如NETLINK_ROUTE用于获取和设置路由与链路信息、NETLINK_KOBJECT_UEVENT用于内核向用户空间的udev进程发送通知等。

```go
// golang.org/x/sys@v0.30.0/unix/zerrors_linux.go
import "syscall"

const (
	// ...

	NETLINK_ADD_MEMBERSHIP                      = 0x1
	NETLINK_AUDIT                               = 0x9
	NETLINK_BROADCAST_ERROR                     = 0x4
	NETLINK_CAP_ACK                             = 0xa
	NETLINK_CONNECTOR                           = 0xb
	NETLINK_CRYPTO                              = 0x15
	NETLINK_DNRTMSG                             = 0xe
	NETLINK_DROP_MEMBERSHIP                     = 0x2
	NETLINK_ECRYPTFS                            = 0x13
	NETLINK_EXT_ACK                             = 0xb
	NETLINK_FIB_LOOKUP                          = 0xa
	NETLINK_FIREWALL                            = 0x3 // 接收IPV4防火墙代码发送的数据包。
	NETLINK_GENERIC                             = 0x10
	NETLINK_GET_STRICT_CHK                      = 0xc
	NETLINK_INET_DIAG                           = 0x4
	NETLINK_IP6_FW                              = 0xd
	NETLINK_ISCSI                               = 0x8
	NETLINK_KOBJECT_UEVENT                      = 0xf
	NETLINK_LISTEN_ALL_NSID                     = 0x8
	NETLINK_LIST_MEMBERSHIPS                    = 0x9
	NETLINK_NETFILTER                           = 0xc
	NETLINK_NFLOG                               = 0x5 // 用户态的iptables管理工具和内核中的netfilter模块之间通讯的通
	NETLINK_NO_ENOBUFS                          = 0x5
	NETLINK_PKTINFO                             = 0x3
	NETLINK_RDMA                                = 0x14
	NETLINK_ROUTE                               = 0x0
	NETLINK_RX_RING                             = 0x6
	NETLINK_SCSITRANSPORT                       = 0x12
	NETLINK_SELINUX                             = 0x7
	NETLINK_SMC                                 = 0x16
	NETLINK_SOCK_DIAG                           = 0x4
	NETLINK_TX_RING                             = 0x7
	NETLINK_UNUSED                              = 0x1
	NETLINK_USERSOCK                            = 0x2
	NETLINK_XFRM                                = 0x6
	
	// ..
}
```

## github.com/vishvananda/netlink

netlink 包为 go 提供了一个简单的 netlink 库。
Netlink 是 linux用户态程序用来与内核通信的接口。
它可用于添加和删除接口、设置 ip 地址和路由以及配置 ipsec。Netlink 通信需要提升权限，因此在大多数情况下，此代码需要以 root 身份运行。


初始化 handle 

```go
func newHandle(newNs, curNs netns.NsHandle, nlFamilies ...int) (*Handle, error) {
	h := &Handle{sockets: map[int]*nl.SocketHandle{}}
	fams := nl.SupportedNlFamilies
	if len(nlFamilies) != 0 {
		fams = nlFamilies
	}
	for _, f := range fams {
		// 创建NETLINK socket
		s, err := nl.GetNetlinkSocketAt(newNs, curNs, f)
		if err != nil {
			return nil, err
		}
		// 根据 netlink family 区分 socket
		h.sockets[f] = &nl.SocketHandle{Socket: s}
	}
	return h, nil
}
```

```go
func GetNetlinkSocketAt(newNs, curNs netns.NsHandle, protocol int) (*NetlinkSocket, error) {
	// 获取恢复的操作
	c, err := executeInNetns(newNs, curNs)
	if err != nil {
		return nil, err
	}
	defer c()
	return getNetlinkSocket(protocol)
}


func executeInNetns(newNs, curNs netns.NsHandle) (func(), error) {
	var (
		err       error
		moveBack  func(netns.NsHandle) error
		closeNs   func() error
		unlockThd func()
	)
	// 恢复
	restore := func() {
		// order matters
		if moveBack != nil {
			moveBack(curNs)
		}
		if closeNs != nil {
			closeNs()
		}
		if unlockThd != nil {
			unlockThd()
		}
	}
	if newNs.IsOpen() {
		// 锁定当前线程
		runtime.LockOSThread()
		unlockThd = runtime.UnlockOSThread
		if !curNs.IsOpen() {
			if curNs, err = netns.Get(); err != nil {
				restore()
				return nil, fmt.Errorf("could not get current namespace while creating netlink socket: %v", err)
			}
			closeNs = curNs.Close
		}
		if err := netns.Set(newNs); err != nil {
			restore()
			return nil, fmt.Errorf("failed to set into network namespace %d while creating netlink socket: %v", newNs, err)
		}
		moveBack = netns.Set
	}
	return restore, nil
}


func getNetlinkSocket(protocol int) (*NetlinkSocket, error) {
	// 使用socket()函数创建一个socket，
	// socket域(地址族)是AF_NETLINK,socket的类型是SOCK_RAW或者SOCK_DGRAM,因为netlink是一种面向数据包的服务。
	// protocol 协议类型选择netlink要使用的类型即可。
	fd, err := unix.Socket(unix.AF_NETLINK, unix.SOCK_RAW|unix.SOCK_CLOEXEC, protocol)
	if err != nil {
		return nil, err
	}
	err = unix.SetNonblock(fd, true)
	if err != nil {
		return nil, err
	}
	s := &NetlinkSocket{
		fd:   int32(fd),
		file: os.NewFile(uintptr(fd), "netlink"),
	}
	s.lsa.Family = unix.AF_NETLINK
	// netlink的bind()函数把一个本地socket地址(源socket地址)与一个打开的socket进行关联
	if err := unix.Bind(fd, &s.lsa); err != nil {
		unix.Close(fd)
		return nil, err
	}

	return s, nil
}
```

## 参考

- [linux下netlink的使用简介](https://www.cnblogs.com/wanghuaijun/p/15712343.html)