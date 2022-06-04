# Sentry
Sentry 是一个开源的实时错误追踪系统，可以帮助开发者实时监控并修复异常问题。它主要专注于持续集成、提高效率并且提升用户体验。


sentry支持自动收集和手动收集两种错误收集方法.

使用sentry需要结合两个部分，客户端与sentry服务端；客户端就像你需要去监听的对象，比如公司的前端项目，而服务端就是给你展示已搜集的错误信息，项目管理，组员等功能的一个服务平台。


## 为什么用Sentry

1. 多项目，多用户,界面友好
2. 报警的规则多样性：可以配置异常出发规则，例如发送邮件，钉钉
3. 报警的及时性：不需要自己再去额外集成报警系统，一旦产生了 issue 便以邮件通知到项目组的每个成员。
4. 问题关联信息的聚合：每个问题不仅有一个整体直观的描绘，聚合的日志信息省略了人工从海量日志中寻找线索，免除大量无关信息的干扰。
5. 丰富的上下文：Sentry 不仅丰富还规范了上下文的内容，也让我们意识到更多的有效内容，提高日志的质量。
6. 支持主流语言接口：从Sentry的文档首页截下来的一张图，可以看到它支持目前主流的编程语言。
7. issues & events：在相同地方产生的异常会被归纳为一个「Issue」，每次在这个地方产生的异常叫做「Event」。所以在同一个地方触发两次异常，仍然只有一个Issue，但是可以在Event页面看到多个[Event]。
8. 聚合策略：Sentry 按照策略将日志事件进行聚合，从而提供一个 issue的events 。这么做就是为了智能地帮助我们组合关联的日志信息，减少人工的日志信息的提取工作量，关注一个 issue 首先关注这些聚合的事件。但是这个策略分组并不会那么智能，Sentry 主要按照以下几个方面，优先级从高到低进行日志事件的聚合：Stacktrace、Exception、Template、Messages。
9. alerts digest & limit：默认 Sentry 的 alerts 会发送邮件。当一个 issue 产生或者一组 issue 产生时，项目相关的成员都会受到邮件。但是并不是每次 issue 有更新就会产生 alert 。考虑到用户也不希望被一箩筐的报警邮件给轰炸，因为过多相当于没有， Sentry 除了对重复的报警进行抑制，还会追加一段时间内更新 issue 的摘要（digest）到下一个报警，这样，用户邮件上接收到的信息会充分压缩，不用苦恼于过多的邮件。另外，每个用户可以根据自己的喜好自行配置报警的时间间隔。


## 基本概念

### DSN
DSN是连接客户端(项目)与sentry服务端,让两者能够通信的钥匙；
每当我们在sentry服务端创建一个新的项目，都会得到一个独一无二的DSN，也就是密钥。在客户端初始化时会用到这个密钥，这样客户端报错，服务端就能抓到你对应项目的错误了


### event

每当项目产生一个错误，sentry服务端日志就会产生一个event，记录此次报错的具体信息。一个错误，对应一个event


### issue

同一类event的集合，一个错误可能会重复产生多次，sentry服务端会将这些错误聚集在一起，那么这个集合就是一个issue。




## sentry-Go 初始化客户端配置
```go
// ClientOptions that configures a SDK Client.
type ClientOptions struct {
	// DSN是连接客户端(项目)与sentry服务端,让两者能够通信的钥匙
	Dsn string

    // 调试模式会sentry打印结果到控制台
	Debug bool

    //栈信息追踪
	AttachStacktrace bool
	
	// 事件提交的采样率（0.0-1.0，默认为 1.0）
	SampleRate float64
	// The sample rate for sampling traces in the range [0.0, 1.0].
	TracesSampleRate float64
	// Used to customize the sampling of traces, overrides TracesSampleRate.
	TracesSampler TracesSampler
    // 用于与事件消息进行匹配的正则表达式字符串列表，如果适用，
    // 则捕获错误类型和值。如果找到匹配项，则将删除整个事件。
	IgnoreErrors []string
	
	// 发送回调之前
	// See EventProcessor if you need to mutate transactions.
	BeforeSend func(event *Event, hint *EventHint) *Event
	// 在面包屑之前添加回调.
	BeforeBreadcrumb func(breadcrumb *Breadcrumb, hint *BreadcrumbHint) *Breadcrumb
	// 要在当前客户端上安装的集成，接收默认集成
	Integrations func([]Integration) []Integration
	// io.Writer implementation that should be used with the Debug mode.
	DebugWriter io.Writer
	// The transport to use. Defaults to HTTPTransport.
	Transport Transport
	// 要报告的服务器名称
	ServerName string
	// 与事件一起发送的版本
	Release string
	// 与事件一起发送的 dist
	Dist string
	//  与事件一起发送的环境
	Environment string
	// 面包屑的最大数量
	MaxBreadcrumbs int
	// An optional pointer to http.Client that will be used with a default
	// HTTPTransport. Using your own client will make HTTPTransport, HTTPProxy,
	// HTTPSProxy and CaCerts options ignored.
	HTTPClient *http.Client
	// An optional pointer to http.Transport that will be used with a default
	// HTTPTransport. Using your own transport will make HTTPProxy, HTTPSProxy
	// and CaCerts options ignored.
	HTTPTransport http.RoundTripper
	// 要使用的可选 HTTP,HTTPS 代理。
	HTTPProxy string
	HTTPSProxy string
	
	//  要使用的可选 CaCert,默认为 `gocertifi.CACerts()`
	CaCerts *x509.CertPool
}
```