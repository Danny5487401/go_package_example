#Config包

##Config 特征
1 动态加载：根据需要动态加载多个资源文件。 go config 在后台管理并监控配置文件，并自动更新到内存中
2 资源可插拔： 从任意数量的源中进行选择以加载和合并配置。后台资源源被抽象为内部使用的标准格式，并通过编码器进行解码。源可以是环境变量，标志，文件，etcd，k8s configmap等
3 可合并的配置：如果指定了多种配置，它们会合并到一个视图中。
4 监控变化：可以选择是否监控配置的指定值，热重启
5 安全恢复： 万一配置加载错误或者由于未知原因而清除，可以指定回退值进行回退

##Source 是资源加载来源
    *cli 命令行
    *consul
    *env
    *etcd
    *file
    *flag
    *memory
    *也有一些社区支持的插件：
    *configmap - nread from k8s configmap
    *grpc - read from grpc server
    *runtimevar - read from Go Cloud Development Kit runtime variable
    *url - read from URL
    *vault - read from Vault server

##ChangeSet
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
##Encoder 处理配置文件的编码与解码，支持
    json
    yaml
    toml
    xml
    xml
    hcl

##Config动态配置接口抽象
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
