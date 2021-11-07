#链路追踪(Distributed Tracing)
是一种用于分析和监控应用程序的方法，尤其是使用微服务架构构建的应用程序。分布式跟踪有助于精确定位故障发生的位置以及导致性能差的原因

##常用的链路追踪系统
    *Skywalking
    *阿里 鹰眼
    *大众点评 CAT
    *Twitter Zipkin
    *Naver pinpoint
    *Uber Jaeger

##分类
![](.introduction_images/openTracing_n_OpenCensus.png)
开源领域主要分为两派，一派是以 CNCF技术委员 会为主的 OpenTracing 的规范，例如 jaeger zipkin 都是遵循了OpenTracing 的规范。而另一派则是谷歌作为发起者的 OpenCensus，而且谷歌本身还是最早提出链路追踪概念的公司，后期连微软也加入了 OpenCensus

##OpenTelemetry 诞生
微软加入 OpenCensus 后，直接打破了之前平衡的局面，间接的导致了 OpenTelemetry 的诞生 谷歌和微软下定决心结束江湖之乱，首要的问题是如何整合两个两个社区已有的项目，OpenTelemetry 主要的理念就是，兼容 OpenCensus 和 OpenTracing ，可以让使用者无需改动或者很小的改动就可以接入 OpenTelemetry

## 链路监控组件要求
探针的性能消耗。
    APM组件服务的影响应该做到足够小。在一些高度优化过的服务，即使一点点损耗也会很容易察觉到，而且有可能迫使在线服务的部署团队不得不将跟踪系统关停。

代码的侵入性
    对于应用的程序员来说，是不需要知道有跟踪系统这回事的。如果一个跟踪系统想生效，就必须需要依赖应用的开发者主动配合，那么这个跟踪系统也太脆弱了，往往由于跟踪系统在应用中植入代码的bug或疏忽导致应用出问题，这样才是无法满足对跟踪系统“无所不在的部署”这个需求。

可扩展性
    能够支持的组件越多当然越好。或者提供便捷的插件开发API，对于一些没有监控到的组件，应用开发者也可以自行扩展。

数据的分析
    数据的分析要快 ，分析的维度尽可能多。跟踪系统能提供足够快的信息反馈，就可以对生产环境下的异常状况做出快速反应。分析的全面，能够避免二次开发

## 基础概念
Trace:
    类似于树结构的Span集合，表示一条调用链路，存在唯一标识。比如你运行的分布式大数据存储一次Trace就由你的一次请求组成。

Span   
    基本工作单元，一次链路调用(可以是RPC，DB等没有特定的限制)创建一个span，通过一个64位ID标识它，uuid较为方便，span中还有其他的数据，
    例如描述信息，时间戳，key-value对的(Annotation)tag信息，parent-id等,其中parent-id可以表示span调用链路来源。



Annotation: 注解,用来记录请求特定事件相关信息(例如时间)，通常包含四个注解信息：   
(1) cs：Client Start,表示客户端发起请求

(2) sr：Server Receive,表示服务端收到请求

(3) ss：Server Send,表示服务端完成处理，并将结果发送给客户端

(4) cr：Client Received,表示客户端获取到服务端返回信息