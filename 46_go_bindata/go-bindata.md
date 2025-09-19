<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [go-bindata](#go-bindata)
  - [需求](#%E9%9C%80%E6%B1%82)
  - [第三方使用](#%E7%AC%AC%E4%B8%89%E6%96%B9%E4%BD%BF%E7%94%A8)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# go-bindata

嵌入静态文件,在 Go 1.16 中包含了 go embed 的功能，使用它可以代替 go-bindata 将文件嵌入到可执行的二进制文件中.

## 需求

在日常开发工作中，有一些配置文件之类的静态内容我们是剥离在项目之外存在的，而如果想要实现该项目打出来的二进制包能够不依赖本地静态配置.


把静态资源嵌入在程序里，原因无外乎以下几点：

- 布署程序更简单。传统部署要么需要把静态资源和编译好的程序一起打包上传，要么使用docker和dockerfile自动化.
- 保证程序完整性。运行中发生静态资源损坏或丢失往往会影响程序的正常运行.
- 可以自主控制程序需要的静态资源.



## 第三方使用

github.com/shuLhan/go-bindata使用: 基于 eBPF 的开源项目 eCapture
```makefile
# https://github.com/gojue/ecapture/blob/v1.4.1/Makefile

.PHONY: assets
assets: .checkver_$(CMD_GO) ebpf ebpf_noncore
	$(CMD_GO) run github.com/shuLhan/go-bindata/cmd/go-bindata $(IGNORE_LESS52) -pkg assets -o "assets/ebpf_probe.go" $(wildcard ./user/bytecode/*.o)

```



## 参考

- [https://cloud.tencent.com/developer/article/2295688](https://cloud.tencent.com/developer/article/2295688)
- [深入浅出 Golang 资源嵌入方案：go-bindata篇](https://juejin.cn/post/7053644550488719396)
