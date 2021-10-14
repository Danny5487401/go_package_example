#nacos源码分析

客户端配置
```go
constant.ClientConfig{
	TimeoutMs            uint64 // 请求Nacos服务端的超时时间，默认是10000ms
	NamespaceId          string // ACM的命名空间Id
	Endpoint             string // 当使用ACM时，需要该配置. https://help.aliyun.com/document_detail/130146.html
	RegionId             string // ACM&KMS的regionId，用于配置中心的鉴权
	AccessKey            string // ACM&KMS的AccessKey，用于配置中心的鉴权
	SecretKey            string // ACM&KMS的SecretKey，用于配置中心的鉴权
	OpenKMS              bool   // 是否开启kms，默认不开启，kms可以参考文档 https://help.aliyun.com/product/28933.html
	                            // 同时DataId必须以"cipher-"作为前缀才会启动加解密逻辑
	CacheDir             string // 缓存service信息的目录，默认是当前运行目录
	UpdateThreadNum      int    // 监听service变化的并发数，默认20
	NotLoadCacheAtStart  bool   // 在启动的时候不读取缓存在CacheDir的service信息
	UpdateCacheWhenEmpty bool   // 当service返回的实例列表为空时，不更新缓存，用于推空保护
	Username             string // Nacos服务端的API鉴权Username
	Password             string // Nacos服务端的API鉴权Password
	LogDir               string // 日志存储路径
	RotateTime           string // 日志轮转周期，比如：30m, 1h, 24h, 默认是24h
	MaxAge               int64  // 日志最大文件数，默认3
	LogLevel             string // 日志默认级别，值必须是：debug,info,warn,error，默认值是info
}
```
服务端配置
```go
constant.ServerConfig{
	ContextPath string // Nacos的ContextPath
	IpAddr      string // Nacos的服务地址
	Port        uint64 // Nacos的服务端口
	Scheme      string // Nacos的服务地址前缀
}
```
Note：我们可以配置多个ServerConfig，客户端会对这些服务端做轮询请求