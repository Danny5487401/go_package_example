<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [grpc配置retry自动重试](#grpc%E9%85%8D%E7%BD%AEretry%E8%87%AA%E5%8A%A8%E9%87%8D%E8%AF%95)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# grpc配置retry自动重试

grpc内置retry机制

##配置
Service Config 是以 JSON 格式配置的，具体文档见 [service_config.md](https://github.com/grpc/grpc/blob/master/doc/service_config.md)。

service_config.proto
```protobuf
// Configuration for a method.
message MethodConfig {
  // The names of the methods to which this configuration applies.
  // - MethodConfig without names (empty list) will be skipped.
  // - Each name entry must be unique across the entire ServiceConfig.
  // - If the 'method' field is empty, this MethodConfig specifies the defaults
  //   for all methods for the specified service.
  // - If the 'service' field is empty, the 'method' field must be empty, and
  //   this MethodConfig specifies the default for all methods (it's the default
  //   config).
  //
  // When determining which MethodConfig to use for a given RPC, the most
  // specific match wins. For example, let's say that the service config
  // contains the following MethodConfig entries:
  //
  // method_config { name { } ... }
  // method_config { name { service: "MyService" } ... }
  // method_config { name { service: "MyService" method: "Foo" } ... }
  //
  // MyService/Foo will use the third entry, because it exactly matches the
  // service and method name. MyService/Bar will use the second entry, because
  // it provides the default for all methods of MyService. AnotherService/Baz
  // will use the first entry, because it doesn't match the other two.
  //
  // In JSON representation, value "", value `null`, and not present are the
  // same. The following are the same Name:
  // - { "service": "s" }
  // - { "service": "s", "method": null }
  // - { "service": "s", "method": "" }
  message Name {
    string service = 1;  // Required. Includes proto package name.
    string method = 2;
  }
  repeated Name name = 1;

  // Whether RPCs sent to this method should wait until the connection is
  // ready by default. If false, the RPC will abort immediately if there is
  // a transient failure connecting to the server. Otherwise, gRPC will
  // attempt to connect until the deadline is exceeded.
  //
  // The value specified via the gRPC client API will override the value
  // set here. However, note that setting the value in the client API will
  // also affect transient errors encountered during name resolution, which
  // cannot be caught by the value here, since the service config is
  // obtained by the gRPC client via name resolution.
  google.protobuf.BoolValue wait_for_ready = 2;

  // The default timeout in seconds for RPCs sent to this method. This can be
  // overridden in code. If no reply is received in the specified amount of
  // time, the request is aborted and a DEADLINE_EXCEEDED error status
  // is returned to the caller.
  //
  // The actual deadline used will be the minimum of the value specified here
  // and the value set by the application via the gRPC client API.  If either
  // one is not set, then the other will be used.  If neither is set, then the
  // request has no deadline.
  google.protobuf.Duration timeout = 3;

  // The maximum allowed payload size for an individual request or object in a
  // stream (client->server) in bytes. The size which is measured is the
  // serialized payload after per-message compression (but before stream
  // compression) in bytes. This applies both to streaming and non-streaming
  // requests.
  //
  // The actual value used is the minimum of the value specified here and the
  // value set by the application via the gRPC client API.  If either one is
  // not set, then the other will be used.  If neither is set, then the
  // built-in default is used.
  //
  // If a client attempts to send an object larger than this value, it will not
  // be sent and the client will see a ClientError.
  // Note that 0 is a valid value, meaning that the request message
  // must be empty.
  google.protobuf.UInt32Value max_request_message_bytes = 4;

  // The maximum allowed payload size for an individual response or object in a
  // stream (server->client) in bytes. The size which is measured is the
  // serialized payload after per-message compression (but before stream
  // compression) in bytes. This applies both to streaming and non-streaming
  // requests.
  //
  // The actual value used is the minimum of the value specified here and the
  // value set by the application via the gRPC client API.  If either one is
  // not set, then the other will be used.  If neither is set, then the
  // built-in default is used.
  //
  // If a server attempts to send an object larger than this value, it will not
  // be sent, and a ServerError will be sent to the client instead.
  // Note that 0 is a valid value, meaning that the response message
  // must be empty.
  google.protobuf.UInt32Value max_response_message_bytes = 5;
  
  // 重试策略
  message RetryPolicy {
    // 最大尝试次数
    // This field is required and must be greater than 1.
    // Any value greater than 5 will be treated as if it were 5.
    uint32 max_attempts = 1;

    // 默认退避时间. The initial retry attempt will occur at
    // random(0, initial_backoff). In general, the nth attempt will occur at
    // random(0,
    //   min(initial_backoff*backoff_multiplier**(n-1), max_backoff)).
    // Required. Must be greater than zero.
    google.protobuf.Duration initial_backoff = 2;
    
    // 最大退避时间.
    google.protobuf.Duration max_backoff = 3;
    
    // 退避时间增加倍率
    float backoff_multiplier = 4;  // Required. Must be greater than zero.


    // 服务端返回什么错误码才重试.
    repeated google.rpc.Code retryable_status_codes = 5;
  }
  // 对冲策略. Hedged RPCs may execute more than
  // once on the server, so only idempotent methods should specify a hedging
  // policy.
  message HedgingPolicy {
    // The hedging policy will send up to max_requests RPCs.
    // This number represents the total number of all attempts, including
    // the original attempt.
    //
    // This field is required and must be greater than 1.
    // Any value greater than 5 will be treated as if it were 5.
    uint32 max_attempts = 1;

    // The first RPC will be sent immediately, but the max_requests-1 subsequent
    // hedged RPCs will be sent at intervals of every hedging_delay. Set this
    // to 0 to immediately send all max_requests RPCs.
    google.protobuf.Duration hedging_delay = 2;

    // The set of status codes which indicate other hedged RPCs may still
    // succeed. If a non-fatal status code is returned by the server, hedged
    // RPCs will continue. Otherwise, outstanding requests will be canceled and
    // the error returned to the client application layer.
    //
    // This field is optional.
    repeated google.rpc.Code non_fatal_status_codes = 3;
  }

  // Only one of retry_policy or hedging_policy may be set. If neither is set,
  // RPCs will not be retried or hedged.
  oneof retry_or_hedging_policy {
    RetryPolicy retry_policy = 6;
    HedgingPolicy hedging_policy = 7;
  }
}
```
注释写的还是很详细的，转换成 JSON 如下：
```json
{
		"methodConfig": [{
		  "name": [{"service": "echo.Echo","method":"UnaryEcho"}],
          "wait_for_ready": false,
          "timeout": 1000ms,
          "max_request_message_bytes": 1024,
          "max_response_message_bytes": 1024,
		  "retryPolicy": {
			  "maxAttempts": 4,
			  "initialBackoff": ".01s",
			  "maxBackoff": ".01s",
			  "backoffMultiplier": 1.0,
			  "retryableStatusCodes": [ "UNAVAILABLE" ]
		  },
		  "hedgingPolicy":{
              "maxAttempts":4,
              "hedgingDelay":"0.1s",
              "nonFatalStatusCodes": [ "" ]
          }}]
}
```
gRPC 的重试策略有两种分别是 重试(retryPolicy)和对冲(hedging)，一个RPC方法只能配置一种重试策略。

对冲是指在不等待响应的情况主动发送单次调用的多个请求，如果一个方法使用对冲策略，那么首先会像正常的 RPC 调用一样发送第一次请求，
如果 hedgingDelay 时间内没有响应，那么直接发送第二次请求，以此类推，直到发送了 maxAttempts 次。

对冲在超过指定时间没有响应就会直接发起请求，而重试则必须要服务端响应后才会发起请求。

name 指定下面的配置信息作用的 RPC 服务或方法
- service：通过服务名匹配，语法为<package>.<service> package就是proto文件中指定的package，service也是proto文件中指定的 Service Name。
- method：匹配具体某个方法，proto文件中定义的方法名。

如果不使用退避算法，失败后就一直重试只会增加服务器的压力。如果是因为服务器压力大，导致的请求失败，那么根据退避算法等待一定时间后再次请求可能就能成功。反之直接请求可能会因为压力过大导致服务崩溃。

- 第一次重试间隔是 random(0, initialBackoff)
- 第 n 次的重试间隔为 random(0, min( initialBackoff*backoffMultiplier**(n-1) , maxBackoff))