#Go-redis源码分析

##初始化
Client结构体
```go
type Client struct {
	baseClient
	cmdable

	ctx context.Context
}

// baseClient 是真正的客户端值，负责获取连接。
type baseClient struct {
    opt      *Options
    connPool pool.Pooler
    limiter  Limiter
    
    process           func(Cmder) error
    processPipeline   func([]Cmder) error
    processTxPipeline func([]Cmder) error
    
    onClose func() error // hook called when client is closed
}

type Cmder interface {
    Name() string
    Args() []interface{}
    stringArg(int) string
    
    readReply(rd *proto.Reader) error
    setErr(error)
    
    readTimeout() *time.Duration
    
    Err() error
}

```


初始化客户端,注意：defaultProcess 这个实际处理 CMD 的函数会贯穿整个客户端
```go
func NewClient(opt *Options) *Client {
    opt.init()
 
    c := Client{
        baseClient: baseClient{
            opt:      opt,
            connPool: newConnPool(opt), // 生成管理连接实例
        },
    }
    c.baseClient.init() // 客户端初始化，比较重要的是设置了 process：c.process = c.defaultProcess
    c.init() // 就是把 Client的  process 设置为 c.defaultProcess
 
    return &c
}
```

##执行命令过程
以string类型的 Set 方法为例子
```go
func (c *cmdable) Set(key string, value interface{}, expiration time.Duration) *StatusCmd {
	args := make([]interface{}, 3, 4)
	args[0] = "set"
	args[1] = key
	args[2] = value
	if expiration > 0 {
		if usePrecise(expiration) {
			args = append(args, "px", formatMs(expiration))
		} else {
			args = append(args, "ex", formatSec(expiration))
		}
	}
	cmd := NewStatusCmd(args...)
	c.process(cmd)
	return cmd
}
```
我们在执行命令的时候会调用到 c.process(cmd) 这个函数，在之前我们提到了 c.process = c.defaultProcess，defaultProcess 接收一个 Cmder 接口
```go
func (c *baseClient) init() {
	c.process = c.defaultProcess
	c.processPipeline = c.defaultProcessPipeline
	c.processTxPipeline = c.defaultProcessTxPipeline
}
```
```go
func (c *baseClient) defaultProcess(cmd Cmder) error {
    // 重试次数
    for attempt := 0; attempt <= c.opt.MaxRetries; attempt++ {
        if attempt > 0 {
            time.Sleep(c.retryBackoff(attempt))
        }
        // 获取到连接
        cn, err := c.getConn()
        if err != nil {
            cmd.setErr(err)
            if internal.IsRetryableError(err, true) {
                continue
            }
            return err
        }
        // 往网络连接中写入数据
        err = cn.WithWriter(c.opt.WriteTimeout, func(wr *proto.Writer) error {
            return writeCmd(wr, cmd)
        })
        if err != nil {
           // 释放连接
            c.releaseConn(cn, err)
            cmd.setErr(err)
            if internal.IsRetryableError(err, true) {
                continue
            }
            return err
        }
        // 读取数据
        err = cn.WithReader(c.cmdTimeout(cmd), cmd.readReply)
        c.releaseConn(cn, err)
        
        if err != nil && internal.IsRetryableError(err, cmd.readTimeout() == nil) {
            continue
        }
 
        return err
    }
 
    return cmd.Err()
}

```
总体的流程就是：获取到连接 → 往网络连接中写入数据 → 读取数据→ 释放连接→ 返回结果。其中比较重要的是 c.getConn()，这里会从链接池里面获取链接
