# 链路追踪(Distributed Tracing)
是一种用于分析和监控应用程序的方法，尤其是使用微服务架构构建的应用程序。分布式跟踪有助于精确定位故障发生的位置以及导致性能差的原因.

目前在Tracing技术这块比较有影响力的是两大开源技术框架：Netflix公司开源的OpenTracing和Google开源的OpenCensus.

两大框架都拥有比较高的开发者群体。为形成统一的技术标准，两大框架最终磨合成立了OpenTelemetry项目，简称otel.

## 常用的链路追踪系统
* Skywalking
* 阿里 鹰眼
* 大众点评 CAT
* Twitter Zipkin
* Naver pinpoint
* Uber Jaeger

## 链路监控组件要求
- 探针的性能消耗

  APM组件服务的影响应该做到足够小。在一些高度优化过的服务，即使一点点损耗也会很容易察觉到，而且有可能迫使在线服务的部署团队不得不将跟踪系统关停。

- 代码的侵入性

  对于应用的程序员来说，是不需要知道有跟踪系统这回事的。如果一个跟踪系统想生效，就必须需要依赖应用的开发者主动配合，那么这个跟踪系统也太脆弱了，往往由于跟踪系统在应用中植入代码的bug或疏忽导致应用出问题，这样才是无法满足对跟踪系统“无所不在的部署”这个需求。

- 可扩展性

  能够支持的组件越多当然越好。或者提供便捷的插件开发API，对于一些没有监控到的组件，应用开发者也可以自行扩展。

- 数据的分析

  数据的分析要快 ，分析的维度尽可能多。跟踪系统能提供足够快的信息反馈，就可以对生产环境下的异常状况做出快速反应。分析的全面，能够避免二次开发

## 基础概念

### Tracer
Tracer表示一次完整的追踪链路，用来创建Span，以及处理如何处理Inject(serialize) 和 Extract (deserialize),tracer由一个或多个span组成。
下图示例表示了一个由8个span组成的tracer:
```css
        [Span A]  ←←←(the root span)
            |
     +------+------+
     |             |
 [Span B]      [Span C] ←←←(Span C is a `ChildOf` Span A)
     |             |
 [Span D]      +---+-------+
               |           |
           [Span E]    [Span F] >>> [Span G] >>> [Span H]
                                       ↑
                                       ↑
                                       ↑
                         (Span G `FollowsFrom` Span F)
```
时间轴的展现方式
```css
––|–––––––|–––––––|–––––––|–––––––|–––––––|–––––––|–––––––|–> time
 [Span A···················································]
   [Span B··············································]
      [Span D··········································]
    [Span C········································]
         [Span E·······]        [Span F··] [Span G··] [Span H··]
```
代码方式
```go
otel.Tracer(tracerName)
```

### Span   
Span是一条追踪链路中的基本组成要素，一个span表示一个独立的工作单元，比如可以表示一次函数调用，一次http请求等等。span会记录如下基本要素:
- 服务名称（operation name）
- 服务的开始时间和结束时间
- K/V形式的Tags
- K/V形式的Logs
- SpanContext

代码方式
```go
otel.Tracer(tracerName).Start(ctx, spanName, opts ...)
```

### Attributes
Attributes以K/V键值对的形式保存用户自定义标签，主要用于链路追踪结果的查询过滤。例如： http.method="GET",http.status_code=200。
其中key值必须为字符串，value必须是字符串，布尔型或者数值型。 span中的Attributes仅自己可见，不会随着 SpanContext传递给后续span。

设置Attributes方式例如
```go
span.SetAttributes(
    label.String("http.remote", conn.RemoteAddr().String()),
    label.String("http.local", conn.LocalAddr().String()),
)
```

### events

Events与Attributes类似，也是K/V键值对形式。与Attributes不同的是，Events还会记录写入Events的时间，因此Events主要用于记录某些事件发生的时间。
Events的key值同样必须为字符串，但对value类型则没有限制

```go
span.AddEvent("http.request", trace.WithAttributes(
    label.Any("http.request.header", headers),
    label.Any("http.request.baggage", gtrace.GetBaggageMap(ctx)),
    label.String("http.request.body", bodyContent),
))
```
### Annotation: 注解,用来记录请求特定事件相关信息(例如时间)
通常包含四个注解信息：   
(1) cs：Client Start,表示客户端发起请求

(2) sr：Server Receive,表示服务端收到请求

(3) ss：Server Send,表示服务端完成处理，并将结果发送给客户端

(4) cr：Client Received,表示客户端获取到服务端返回信息


### SpanContext
SpanContext携带着一些用于跨服务通信的（跨进程）数据，主要包含：

- 足够在系统中标识该span的信息，比如：span_id, trace_id。
- Baggage - 为整条追踪连保存跨服务（跨进程）的K/V格式的用户自定义数据。Baggage 与 Attributes 类似，也是 K/V 键值对。与 Attributes 不同的是：
  - key跟value都只能是字符串格式
  - Baggage不仅当前span可见，其会随着SpanContext传递给后续所有的子span。要小心谨慎的使用Baggage - 因为在所有的span中传递这些K,V会带来不小的网络和CPU开销。

### Propagator

Propagator传播器用于端对端的数据编码/解码，例如：Client到Server端的数据传输，TraceId、SpanId和Baggage也是需要通过传播器来管理数据传输。
业务端开发者往往对Propagator无感知，只有中间件/拦截器的开发者需要知道它的作用。

OpenTelemetry的标准协议实现库提供了常用的TextMapPropagator，用于常见的文本数据端到端传输。
此外，为保证TextMapPropagator中的传输数据兼容性，不应当带有特殊字符


### OpenTelemetry Baggage
OpenTelemetry Baggage 是一个简单但通用的键值系统。一旦数据被添加为 Baggage，它就可以被所有下游服务访问。这允许有用的信息，如账户和项目 ID，在事务的后期变得可用，而不需要从数据库中重新获取它们。

例如，一个使用项目 ID 作为索引的前端服务可以将其作为 Baggage 添加，允许后端服务也通过项目 ID 对其跨度和指标进行索引。
这信息添加到了http header中，进行上下文传递，因此每增加一个项目都必须被编码为一个头，每增加一个项目都会增加事务中每一个后续网络请求的大小，因此不建议在将大量的非重要的信息添加到Baggage中。



## 分类
开源领域主要分为两派   
![](.introduction_images/openTracing_n_OpenCensus.png)

1. 一派是以 CNCF技术委员 会为主的 OpenTracing 的规范，例如 jaeger zipkin 都是遵循了OpenTracing 的规范。
2. 而另一派则是谷歌作为发起者的 OpenCensus，而且谷歌本身还是最早提出链路追踪概念的公司，后期连微软也加入了 OpenCensus.


## OpenTelemetry 诞生
Opentelemetry 源于 OpenTracing 与 OpenCensus 两大开源社区的合并而来。OpenTracing 在 2016 年由 Ben Sigelman 发起，旨在解决开源 Tracing 实现重复开发监测客户端， 数据模型不统一， Tracing 后端兼容成本高的问题。
OpenCensus 则是由 Google 内部实践而来，结合了 Tracing 和 Metrics 的监测客户端开源工具包。

由于两大开源社区各自的影响力都不小，而存在两个或多个 Tracing 的标准这个事情本身跟社区组建的主旨相违背。于是两大开源社区一拍即合，成立了 OpenTelemetry。


### openTelemetry架构
![](.introduction_images/opentelemetry_structure.png)


### Opentelemetry 项目组成
四个部分内容

- 跨语言规范说明
- 收集、转换、转发遥测数据的工具 Collector
- 各语言监测客户端 API & SDK
- 自动监测客户端与第三方库 Instrumentation & Contrib

### 1. 收集、转换、转发遥测数据的工具 Collector

从架构层面来说，Collector 有两种模式。
1. 一种是把 Collector 部署在应用相同的主机内（如 K8S 的 DaemonSet）， 或者部署在应用相同的 Pod 里面 （如 K8S 中的 Sidecar），应用采集到的遥测数据，直接通过回环网络传递给 Collector。这种模式统称为 Agent 模式。
2. 另一种模式是把 Collector 当作一个独立的中间件，应用把采集到的遥测数据往这个中间件里面传递。这种模式称之为 Gateway 模式。

![](.introduction_images/collector_pipeline.png)
在 Collector 内部设计中，一套数据的流入、处理、流出的过程称为 pipeline。一个 pipeline 有三部分组件组合而成，它们分别是 receiver/ processor/ exporter。

- receiver
负责按照对应的协议格式监听接收遥测数据，并把数据转给一个或者多个 processor
- processor
负责做遥测数据加工处理，如丢弃数据，增加信息，转批处理等，并把数据传递给下一个 processor 或者传递给一个或多个 exporter
- exporter
负责把数据往下一个接收端发送（一般是遥测后端），exporter 可以定义同时从多个不同的 processor 中获取遥测数据
  
### 2. 自动监测客户端与第三方库 Instrumentation & Contrib
如果单纯使用监测客户端 API & SDK 包，许多的操作是需要修改应用代码的。
如添加 Tracing 监测点，记录字段信息，元数据在进程/服务间传递的装箱拆箱等。这种方式具有代码侵入性，不易解耦，而且操作成本高，增加用户使用门槛。
这个时候就可以利用公共组件的设计模式或语言特性等来降低用户使用门槛。

利用公共组件的设计模式，例如在 Golang 的 Gin 组件，实现了 Middleware 责任链设计模式。
我们可以引用 github.com/gin-gonic/gin 库，创建一个 otelgin.Middleware，手动添加到 Middleware 链中，实现 Gin 的快速监测，


## 参考链接
1. https://www.bookstack.cn/read/goframe-1.16-zh/549208391059b05d.md