<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [grpc context](#grpc-context)
  - [源码分析](#%E6%BA%90%E7%A0%81%E5%88%86%E6%9E%90)
    - [流程](#%E6%B5%81%E7%A8%8B)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# grpc context
gRPC 是基于 HTTP/2 协议的。进程间传输定义了一个 metadata 对象，该对象放在 Request-Headers 内，所以通过 metadata 我们可以将上一个进程中的全局对象透传到下一个被调用的进程.

## 源码分析

```go
// /Users/xiaxin/go/pkg/mod/google.golang.org/grpc@v1.32.0/metadata/metadata.go
type MD map[string][]string
```

进程内部我们通过context来传输上下文数据，进程间传递MD的时候，我们也可以从ctx，取出来，进行传递.
```go
//set 数据到 metadata
md := metadata.Pairs("key", "val")
// 新建一个有 metadata 的 context
ctx := metadata.NewOutgoingContext(context.Background(), md)
```

为什么不直接把context里面的数据全取出来，传递给下游呢？
这是出于可维护性和安全性两方面的考虑，如果将ctx所有信息都传递下去，很有可能将一些内部信息泄漏，
另一方面，下游在取ctx的时候，不知道到底传了哪些数据。所以grpc定义了两个context：
- OutgoingContext  用于发送请求一方，包装下游依赖的数据，传递出去。
- IncomingContext  用于服务端接受客户端传递来的context信息.

context中间通过序列化成http2 header的方式进行传输.

```go
type mdIncomingKey struct{}
type mdOutgoingKey struct{}

// NewIncomingContext creates a new context with incoming md attached.
func NewIncomingContext(ctx context.Context, md MD) context.Context {
	return context.WithValue(ctx, mdIncomingKey{}, md)
}

// NewOutgoingContext creates a new context with outgoing md attached. If used
// in conjunction with AppendToOutgoingContext, NewOutgoingContext will
// overwrite any previously-appended metadata.
func NewOutgoingContext(ctx context.Context, md MD) context.Context {
	return context.WithValue(ctx, mdOutgoingKey{}, rawMD{md: md})
}
```
我们可以看到这两个context虽然也是通过context.WithValue 设置数据，通过context.Value来读取数据.

```go
// FromIncomingContext returns the incoming metadata in ctx if it exists.  The
// returned MD should not be modified. Writing to it may cause races.
// Modification should be made to copies of the returned MD.
func FromIncomingContext(ctx context.Context) (md MD, ok bool) {
	md, ok = ctx.Value(mdIncomingKey{}).(MD)
	return
}

// FromOutgoingContextRaw returns the un-merged, intermediary contents
// of rawMD. Remember to perform strings.ToLower on the keys. The returned
// MD should not be modified. Writing to it may cause races. Modification
// should be made to copies of the returned MD.
//
// This is intended for gRPC-internal use ONLY.
func FromOutgoingContextRaw(ctx context.Context) (MD, [][]string, bool) {
	raw, ok := ctx.Value(mdOutgoingKey{}).(rawMD)
	if !ok {
		return nil, nil, false
	}

	return raw.md, raw.added, true
}

// FromOutgoingContext returns the outgoing metadata in ctx if it exists.  The
// returned MD should not be modified. Writing to it may cause races.
// Modification should be made to copies of the returned MD.
func FromOutgoingContext(ctx context.Context) (MD, bool) {
	raw, ok := ctx.Value(mdOutgoingKey{}).(rawMD)
	if !ok {
		return nil, false
	}

	mds := make([]MD, 0, len(raw.added)+1)
	mds = append(mds, raw.md)
	for _, vv := range raw.added {
		mds = append(mds, Pairs(vv...))
	}
	return Join(mds...), ok
}

type rawMD struct {
	md    MD
	added [][]string
}
```

### 流程
直观理解，客户端在发送请求的时候，会初始化一个OutgoingContext，服务端在取的时候，用的是IncomingContext，
中间必然存在一个从OutgoingContext 取数据，放入http2 header，从http2 header 取数据存入IncomingContext 的过程。

1. server端构造IncomingContext 的过程：
从server.go文件ServeHTTP 函数开始
```go
// /Users/xiaxin/go/pkg/mod/google.golang.org/grpc@v1.32.0/server.go
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    st, err := transport.NewServerHandlerTransport(w, r, s.opts.statsHandler)
    // ...
    s.serveStreams(st)

}
       

func (s *Server) serveStreams(st transport.ServerTransport){
	//...
	st.HandleStreams(func(stream *transport.Stream)
	// ...
}

```

创建http2
```go
func (s *Server) newHTTP2Transport(c net.Conn, authInfo credentials.AuthInfo) transport.ServerTransport {
	config := &transport.ServerConfig{
		MaxStreams:            s.opts.maxConcurrentStreams,
		AuthInfo:              authInfo,
		InTapHandle:           s.opts.inTapHandle,
		StatsHandler:          s.opts.statsHandler,
		KeepaliveParams:       s.opts.keepaliveParams,
		KeepalivePolicy:       s.opts.keepalivePolicy,
		InitialWindowSize:     s.opts.initialWindowSize,
		InitialConnWindowSize: s.opts.initialConnWindowSize,
		WriteBufferSize:       s.opts.writeBufferSize,
		ReadBufferSize:        s.opts.readBufferSize,
		ChannelzParentID:      s.channelzID,
		MaxHeaderListSize:     s.opts.maxHeaderListSize,
		HeaderTableSize:       s.opts.headerTableSize,
	}
	st, err := transport.NewServerTransport("http2", c, config)
	if err != nil {
		s.mu.Lock()
		s.errorf("NewServerTransport(%q) failed: %v", c.RemoteAddr(), err)
		s.mu.Unlock()
		c.Close()
		channelz.Warning(logger, s.channelzID, "grpc: Server.Serve failed to create ServerTransport: ", err)
		return nil
	}

	return st
}
```

internal/transport/http2_server.go
```go
unc (t *http2Server) HandleStreams(handle func(*Stream), traceCtx func(context.Context, string) context.Context) {
	defer close(t.readerDone)
	for {
		t.controlBuf.throttle()
		frame, err := t.framer.fr.ReadFrame()
		atomic.StoreInt64(&t.lastRead, time.Now().UnixNano())
		if err != nil {
            // ....
		}
		switch frame := frame.(type) {
		case *http2.MetaHeadersFrame:
			if t.operateHeaders(frame, handle, traceCtx) {
				t.Close()
				break
			}
		case *http2.DataFrame:
			t.handleData(frame)
		case *http2.RSTStreamFrame:
			t.handleRSTStream(frame)
		case *http2.SettingsFrame:
			t.handleSettings(frame)
		case *http2.PingFrame:
			t.handlePing(frame)
		case *http2.WindowUpdateFrame:
			t.handleWindowUpdate(frame)
		case *http2.GoAwayFrame:
			// TODO: Handle GoAway from the client appropriately.
		default:
			if logger.V(logLevel) {
				logger.Errorf("transport: http2Server.HandleStreams found unhandled frame type %v.", frame)
			}
		}
	}
}
// operateHeader takes action on the decoded headers.
func (t *http2Server) operateHeaders(frame *http2.MetaHeadersFrame, handle func(*Stream), traceCtx func(context.Context, string) context.Context) (fatal bool) {
	// ...
	// Attach the received metadata to the context.
	if len(state.data.mdata) > 0 {
        s.ctx = metadata.NewIncomingContext(s.ctx, state.data.mdata)
    }
    // ...
}

```
通过http2的header构造了我们的IncomingContext。


2. 客户端

客户端的请求调用是从Invoke函数开始的
```go

// /Users/xiaxin/go/pkg/mod/google.golang.org/grpc@v1.32.0/call.go
func (cc *ClientConn) Invoke(ctx context.Context, method string, args, reply interface{}, opts ...CallOption) error{
	//... 
	return invoke(ctx, method, args, reply, cc, opts...)
}



func invoke(ctx context.Context, method string, req, reply interface{}, cc *ClientConn, opts ...CallOption) error {
	cs, err := newClientStream(ctx, unaryStreamDesc, cc, method, opts...)
	if err != nil {
		return err
	}
	if err := cs.SendMsg(req); err != nil {
		return err
	}
	return cs.RecvMsg(reply)
}
```
```go
func newClientStream(ctx context.Context, desc *StreamDesc, cc *ClientConn, method string, opts ...CallOption) (_ ClientStream, err error) {
    // ...
	cs := &clientStream{
        callHdr:      callHdr,
        ctx:          ctx,
        methodConfig: &mc,
        opts:         opts,
        callInfo:     c,
        cc:           cc,
        desc:         desc,
        codec:        c.codec,
        cp:           cp,
        comp:         comp,
        cancel:       cancel,
        beginTime:    beginTime,
        firstAttempt: true,
    }
    // ...
    
    
    op := func(a *csAttempt) error { return a.newStream() }
    if err := cs.withRetry(op, func() { cs.bufferForRetryLocked(0, op) }); err != nil {
        cs.finish(err)
        return nil, err
    }
}

func (a *csAttempt) newStream() error {
	cs := a.cs
	cs.callHdr.PreviousAttempts = cs.numRetries
	s, err := a.t.NewStream(cs.ctx, cs.callHdr)
	if err != nil {
		if _, ok := err.(transport.PerformedIOError); ok {
			// Return without converting to an RPC error so retry code can
			// inspect.
			return err
		}
		return toRPCErr(err)
	}
	cs.attempt.s = s
	cs.attempt.p = &parser{r: s}
	return nil
}
```
最终调用啦a.t.NewStream

实现在internal/transport/http2_client.go
```go

func (t *http2Client) NewStream(ctx context.Context, callHdr *CallHdr) (_ *Stream, err error){
    // ...
	headerFields, err := t.createHeaderFields(ctx, callHdr)
}



func (t *http2Client) createHeaderFields(ctx context.Context, callHdr *CallHdr) ([]hpack.HeaderField, error){
	// ... 
    md, added, ok := metadata.FromOutgoingContextRaw(ctx); ok
}

```