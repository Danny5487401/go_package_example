#go-redis源码分析

##连接池分析
```go
// 键不存在时返回内容
const Nil = proto.Nil
```

连接池接口
```go
type Pooler interface {
    NewConn(context.Context) (*Conn, error) // 创建连接
    CloseConn(*Conn) error // 关闭连接

    Get(context.Context) (*Conn, error) // 获取连接
    Put(*Conn) // 放回连接
    Remove(*Conn, error) // 移除连接

    Len() int // 连接池长度
    IdleLen() int // 空闲连接数量
    Stats() *Stats // 连接池统计

    Close() error // 关闭连接池
}
```
连接池结构体
```go
type ConnPool struct {
    opt *Options // 连接池配置

    dialErrorsNum uint32 // 连接错误次数，atomic

    lastDialErrorMu sync.RWMutex // 上一次连接错误锁，读写锁
    lastDialError   error // 上一次连接错误

    queue chan struct{} // 工作连接队列

    connsMu      sync.Mutex // 连接队列锁
    conns        []*Conn // 连接队列
    idleConns    []*Conn // 空闲连接队列
    poolSize     int // 连接池大小
    idleConnsLen int // 空闲连接队列长度

    stats Stats // 连接池统计

    _closed  uint32 // 连接池关闭标志，atomic
    closedCh chan struct{} // 通知连接池关闭通道
}
```
初始化
```
var _ Pooler = (*ConnPool)(nil)

func NewConnPool(opt *Options) *ConnPool {
    p := &ConnPool{
        opt: opt,

        queue:     make(chan struct{}, opt.PoolSize),
        conns:     make([]*Conn, 0, opt.PoolSize),
        idleConns: make([]*Conn, 0, opt.PoolSize),
        closedCh:  make(chan struct{}),
    }

    p.checkMinIdleConns()

    if opt.IdleTimeout > 0 && opt.IdleCheckFrequency > 0 {
        go p.reaper(opt.IdleCheckFrequency)
    }

    return p
}
```
    1。创建连接池，传入连接池配置选项参数 opt，工厂函数根据 opt 创建连接池实例。连接池主要依靠以下四个数据结构实现管理和通信：
        queue： 存储工作连接的缓冲通道
        conns：存储所有连接的切片
        idleConns：存储空闲连接的切片
        closed：用于通知所有协程连接池已经关闭的通道
    2。检查连接池的空闲连接数量是否满足最小空闲连接数量要求，若不满足，则创建足够的空闲连接。
    3。若连接池配置选项规定了空闲连接超时和检查空闲连接频率，则开启一个清理空闲连接的协程
    
关闭
```go
func (p *ConnPool) Close() error {
    if !atomic.CompareAndSwapUint32(&p._closed, 0, 1) {
        return ErrClosed
    }
    close(p.closedCh)

    var firstErr error
    p.connsMu.Lock()
    for _, cn := range p.conns {
        if err := p.closeConn(cn); err != nil && firstErr == nil {
            firstErr = err
        }
    }
    p.conns = nil
    p.poolSize = 0
    p.idleConns = nil
    p.idleConnsLen = 0
    p.connsMu.Unlock()

    return firstErr
}
```
    1。原子性检查连接池是否已经关闭，若没关闭，则将关闭标志置为1
    2。关闭 closedCh 通道，连接池中的所有协程都可以通过判断该通道是否关闭来确定连接池是否已经关闭。
    3。连接队列锁上锁，关闭队列中的所有连接，并置空所有维护连接池状态的数据结构，解锁。
    
过滤
```go
func (p *ConnPool) Filter(fn func(*Conn) bool) error {
    var firstErr error
    p.connsMu.Lock()
    for _, cn := range p.conns {
        if fn(cn) {
            if err := p.closeConn(cn); err != nil && firstErr == nil {
                firstErr = err
            }
        }
    }
    p.connsMu.Unlock()
    return firstErr
// 实质上是遍历连接池中的所有连接，并调用传入的 fn 过滤函数作用在每个连接上，过滤出符合业务要求的连接。
```
清理
```go
func (p *ConnPool) reaper(frequency time.Duration) {
    ticker := time.NewTicker(frequency)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            // It is possible that ticker and closedCh arrive together,
            // and select pseudo-randomly pick ticker case, we double
            // check here to prevent being executed after closed.
            if p.closed() {
                return
            }
            _, err := p.ReapStaleConns()
            if err != nil {
                internal.Logger.Printf("ReapStaleConns failed: %s", err)
                continue
            }
        case <-p.closedCh:
            return
        }
    }
}

func (p *ConnPool) ReapStaleConns() (int, error) {
    var n int
    for {
        p.getTurn()

        p.connsMu.Lock()
        cn := p.reapStaleConn()
        p.connsMu.Unlock()
        p.freeTurn()

        if cn != nil {
            _ = p.closeConn(cn)
            n++
        } else {
            break
        }
    }
    atomic.AddUint32(&p.stats.StaleConns, uint32(n))
    return n, nil
}

func (p *ConnPool) reapStaleConn() *Conn {
    if len(p.idleConns) == 0 {
        return nil
    }

    cn := p.idleConns[0]
    if !p.isStaleConn(cn) {
        return nil
    }

    p.idleConns = append(p.idleConns[:0], p.idleConns[1:]...)
    p.idleConnsLen--
    p.removeConn(cn)

    return cn
}
```

    1.开启一个用于检查并清理过期连接的 goroutine 每隔 frequency 时间遍历检查连接池中是否存在过期连接，并清理。
    2.创建一个时间间隔为 frequency 的计时器，在连接池关闭时关闭该计时器
    3.循环判断计时器是否到时和连接池是否关闭
    4.移除空闲连接队列中的过期连接

