#批处理
pipeline

    场景：参考 redis.io/topics/pipelining，用于解决高吞吐场景 RTT（Round Trip Time）导致的高延迟问题
    
    原理：client 一次性将多个命令写出，再依次读取 server 回复的数据。因为 client 在写完最后一条命令前都不读连接，
    故 server 视角 client socket 一直未发生EPOLLOUT事件，所有命令的处理结果被依次缓冲在client.reply链表

##接口
```go
type Pipeliner interface {
    StatefulCmdable                                   // 保证 pipeline 能直接执行面向用户的命令
    Do(ctx context.Context, args ...interface{}) *Cmd // 命令排队
    Exec(ctx context.Context) ([]Cmder, error)        // 处理队列中的命令
    Discard() error                                   // 清理命令队列，复用 pipeline
}
```
##管道逻辑实现
    Pipeline 本质是命令队列，Do时 append，Exec时 pop，最后按入队顺序逐个读取回复数据返回给调用方
```go
type pipelineExecer func(context.Context, []Cmder) error // 描述如何处理命令队列
type Pipeline struct {
    cmdable, statefulCmdable // 都指向 p.Process
    exec pipelineExecer      // 命令队列的实际执行并不关心，Pipeline 只负责管道逻辑
    cmds []Cmder             /*...*/
}
func (p *Pipeline) Process(ctx context.Context, cmd Cmder) error { // 直接执行命令也是入队等着
    p.cmds = append(p.cmds, cmd)
    return nil
}
func (p *Pipeline) Do(ctx context.Context, args ...interface{}) *Cmd {
    cmd := NewCmd(ctx, args...)
    _ = p.Process(ctx, cmd)
    return cmd
}
func (p *Pipeline) Exec(ctx context.Context) ([]Cmder, error) {
    if len(p.cmds) == 0 { return nil, nil }
    cmds, p.cmds := p.cmds, nil
    return cmds, p.exec(ctx, cmds) // 命令执行委托给 pipelineExecer
}
```
pipelineExecer

```go
func (c *baseClient) pipelineProcessCmds(ctx context.Context, cn *pool.Conn, cmds []Cmder) (bool, error) {
    cn.WithWriter(ctx, c.opt.WriteTimeout, func(wr *proto.Writer) error {
        return writeCmds(wr, cmds) // 将管道中所有命令都写出
    })
    cn.WithReader(ctx, c.opt.ReadTimeout, func(rd *proto.Reader) error {
        return pipelineReadCmds(rd, cmds)
    }) /*...*/
}
// 依次读取各命令的回复数据，设置对应的错误信息
func pipelineReadCmds(rd *proto.Reader, cmds []Cmder) error {
    for _, cmd := range cmds {
        err := cmd.readReply(rd)
        cmd.SetErr(err)
        if err != nil && !isRedisError(err) { // 非逻辑错误（-Err, $-1, *-1）则提前返回
            return err
        }
    } /*...*/
}
```

##事务

    场景：参考 redis.io/topics/transactions，用于在 server 端原子地执行一组命令
    命令：WATCH,UNWATCH监控 key，MULTI开启事务，DISCARD放弃事务，EXEC执行事务
    特性：事务不支持回滚，有 2 类错误用法会导致事务失败
    语法错误：会导致整个事务被放弃，Exec 时没有命令会被执行
    键类型错误：比如INCR string 值类型的 key，Exec 中途遇到此错误不会中断事务，也不回滚之前执行成功的命令；anteriz 认为这是用法问题，测试期间就应发现解决，不要等到生产环境事务执行失败才发现，故 Redis 并未实现事务内各 key 类型动态跟踪功能（增加了实现复杂度、增大了事务执行延迟）
    ACID：Redis 事务只实现了 ACID 中的 Isolation（事务由单线程原子执行） 和 Durability（事务记录会被持久化），不支持回滚，因而不满足 Atomicity 和 Consistency

###事务逻辑实现
函数 Client.Watch 作为事务入口，先执行WATCH keys命令，再把Tx.TxPipeline以闭包形式传给用户，自行往 TxPipeline 中加入事务命令
```go
func (c *Client) Watch(ctx context.Context, fn func(*Tx) error, keys ...string) error {
    tx := c.newTx(ctx)
    tx.Watch(ctx, keys...) /* handle errors...*/
    return fn(tx)
}
func (c *Tx) TxPipelined(ctx context.Context, fn func(Pipeliner) error) ([]Cmder, error) {
    return c.TxPipeline().Pipelined(ctx, fn)
}
func (c *Tx) TxPipeline() Pipeliner {
    return Pipeline{
        exec: func(ctx context.Context, cmds []Cmder) error {
            return c.hooks.processTxPipeline(ctx, cmds, c.baseClient.processTxPipeline)
        },
    } /*...*/
}
```
注意 TxPipeline 和 Pipeline 共用一套 Pipeliner 接口，区别是pipelineExecer执行逻辑有 3 处不同：

    TxPipeline 比 Pipeline 多了预先 WATCH
    TxPipeline 比 Pipeline 在执行 hook 前向命令队列 prepend MULTI 命令，append EXEC 命令
    TxPipeline 比 Pipeline 多了读取 QUEUED 回复，处理 *-1 事务执行失败的错误

最终事务执行会也是批量写命令，批量读
```go

func (c *baseClient) txPipelineProcessCmds(ctx context.Context, cn *pool.Conn, cmds []Cmder) (bool, error) {
    cn.WithWriter(ctx, c.opt.WriteTimeout, func(wr *proto.Writer) error {
        return writeCmds(wr, cmds)
    })
    cn.WithReader(ctx, c.opt.ReadTimeout, func(rd *proto.Reader) error {
        multiCmd := cmds[0].(*StatusCmd)
        cmds = cmds[1 : len(cmds)-1] //
        if err := txPipelineReadQueued(rd, multiCmd, cmds); err != nil { // 提前处理语法错误导致的事务取消
            return err
        }
        return pipelineReadCmds(rd, cmds) // 分解 EXEC 的数组回复依次作为各命令回复
    }) /*...*/
}
```