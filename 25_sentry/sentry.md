# Sentry
Sentry 是一个开源的实时错误追踪系统，可以帮助开发者实时监控并修复异常问题。它主要专注于持续集成、提高效率并且提升用户体验。
## 初始化客户端配置
```go
// ClientOptions that configures a SDK Client.
type ClientOptions struct {
	// The DSN to use. If the DSN is not set, the client is effectively
	// disabled.
	Dsn string

    // 调试模式会sentry打印结果到控制台
	Debug bool

    //栈信息追踪
	AttachStacktrace bool
	// The sample rate for event submission in the range [0.0, 1.0]. By default,
	// all events are sent. Thus, as a historical special case, the sample rate
	// 0.0 is treated as if it was 1.0. To drop all events, set the DSN to the
	// empty string.
	SampleRate float64
	// The sample rate for sampling traces in the range [0.0, 1.0].
	TracesSampleRate float64
	// Used to customize the sampling of traces, overrides TracesSampleRate.
	TracesSampler TracesSampler
	// List of regexp strings that will be used to match against event's message
	// and if applicable, caught errors type and value.
	// If the match is found, then a whole event will be dropped.
	IgnoreErrors []string
	// BeforeSend is called before error events are sent to Sentry.
	// Use it to mutate the event or return nil to discard the event.
	// See EventProcessor if you need to mutate transactions.
	BeforeSend func(event *Event, hint *EventHint) *Event
	// Before breadcrumb add callback.
	BeforeBreadcrumb func(breadcrumb *Breadcrumb, hint *BreadcrumbHint) *Breadcrumb
	// Integrations to be installed on the current Client, receives default
	// integrations.
	Integrations func([]Integration) []Integration
	// io.Writer implementation that should be used with the Debug mode.
	DebugWriter io.Writer
	// The transport to use. Defaults to HTTPTransport.
	Transport Transport
	// The server name to be reported.
	ServerName string
	// The release to be sent with events.
	Release string
	// The dist to be sent with events.
	Dist string
	// The environment to be sent with events.
	Environment string
	// Maximum number of breadcrumbs.
	MaxBreadcrumbs int
	// An optional pointer to http.Client that will be used with a default
	// HTTPTransport. Using your own client will make HTTPTransport, HTTPProxy,
	// HTTPSProxy and CaCerts options ignored.
	HTTPClient *http.Client
	// An optional pointer to http.Transport that will be used with a default
	// HTTPTransport. Using your own transport will make HTTPProxy, HTTPSProxy
	// and CaCerts options ignored.
	HTTPTransport http.RoundTripper
	// An optional HTTP proxy to use.
	// This will default to the HTTP_PROXY environment variable.
	HTTPProxy string
	// An optional HTTPS proxy to use.
	// This will default to the HTTPS_PROXY environment variable.
	// HTTPS_PROXY takes precedence over HTTP_PROXY for https requests.
	HTTPSProxy string
	// An optional set of SSL certificates to use.
	CaCerts *x509.CertPool
}
```