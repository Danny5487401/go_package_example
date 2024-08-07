<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Config包](#config%E5%8C%85)
  - [Config 特征](#config-%E7%89%B9%E5%BE%81)
  - [Source 是资源加载来源](#source-%E6%98%AF%E8%B5%84%E6%BA%90%E5%8A%A0%E8%BD%BD%E6%9D%A5%E6%BA%90)
  - [ChangeSet](#changeset)
  - [Encoder 处理配置文件的编码与解码，支持](#encoder-%E5%A4%84%E7%90%86%E9%85%8D%E7%BD%AE%E6%96%87%E4%BB%B6%E7%9A%84%E7%BC%96%E7%A0%81%E4%B8%8E%E8%A7%A3%E7%A0%81%E6%94%AF%E6%8C%81)
  - [Config动态配置接口抽象](#config%E5%8A%A8%E6%80%81%E9%85%8D%E7%BD%AE%E6%8E%A5%E5%8F%A3%E6%8A%BD%E8%B1%A1)
  - [Reader](#reader)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Config包

## Config 特征
1 动态加载：根据需要动态加载多个资源文件。 go config 在后台管理并监控配置文件，并自动更新到内存中

2 资源可插拔： 从任意数量的源中进行选择以加载和合并配置。后台资源源被抽象为内部使用的标准格式，并通过编码器进行解码。源可以是环境变量，标志，文件，etcd，k8s configmap等

3 可合并的配置：如果指定了多种配置，它们会合并到一个视图中。

4 监控变化：可以选择是否监控配置的指定值，热重启

5 安全恢复： 万一配置加载错误或者由于未知原因而清除，可以指定回退值进行回退

## Source 是资源加载来源
* cli 命令行
* consul
* env
* etcd
* file
* flag
* memory

也有一些社区支持的插件：

* configmap - nread from k8s configmap
* grpc - read from grpc server
* runtimevar - read from Go Cloud Development Kit runtime variable
* url - read from URL
* vault - read from Vault server

## ChangeSet
Source返回ChangeSet，是多个后端的单例内部抽象
```go
// ChangeSet represents a set of changes from a source
type ChangeSet struct {
    // Raw encoded config data
    Data      []byte
    // MD5 checksum of the data
    Checksum  string
    // Encoding format e.g json, yaml, toml, xml
    Format    string
    // Source of the config e.g file, consul, etcd
    Source    string
    // Time of loading or update
    Timestamp time.Time
}
```
## Encoder 处理配置文件的编码与解码，支持
    json
    yaml
    toml
    xml
    xml
    hcl
Go-Micro默认的解析器是json

## Config动态配置接口抽象
它扮演的是一种管理者的角色，它负责配置的新建、配置的读取、配置的同步、配置的监听等功能
```go
// Config is an interface abstraction for dynamic configuration
type Config interface {
    // provide the reader.Values interface
    reader.Values
    // Stop the config loader/watcher
    Close() error
    // Load config sources
    Load(source ...source.Source) error
    // Force a source changeset sync
    Sync() error
    // Watch a value for changes
    Watch(path ...string) (Watcher, error)
}
```

## Reader
它的作用是将多个源的配置合并为一个。因为对于json、yaml、或者是consul，其配置结构都是类似的，最终Go-Micro都会将其转换成map[string]interface{}形式，并进行合并