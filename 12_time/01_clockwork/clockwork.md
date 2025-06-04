<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [github.com/jonboulle/clockwork](#githubcomjonboulleclockwork)
  - [接口](#%E6%8E%A5%E5%8F%A3)
  - [适合用于测试以下场景](#%E9%80%82%E5%90%88%E7%94%A8%E4%BA%8E%E6%B5%8B%E8%AF%95%E4%BB%A5%E4%B8%8B%E5%9C%BA%E6%99%AF)
  - [第三方使用-->etcd](#%E7%AC%AC%E4%B8%89%E6%96%B9%E4%BD%BF%E7%94%A8--etcd)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# github.com/jonboulle/clockwork

在软件测试中，管理时间的流动往往是一项挑战，尤其是当你需要精确控制时间推进以测试某些定时行为时。这就是clockwork出现的原因——一个简洁易用的Go语言虚拟时钟库，它允许你在测试中独立于系统时间操作时间




## 接口 

clockwork的核心是一个FakeClock类型，它实现了Clock接口.
```go
type Clock interface {
	After(d time.Duration) <-chan time.Time
	Sleep(d time.Duration)
	Now() time.Time
}

// FakeClock provides an interface for a clock which can be
// manually advanced through time
type FakeClock interface {
	Clock
	// Advance advances the FakeClock to a new point in time, ensuring any existing
	// sleepers are notified appropriately before returning
	Advance(d time.Duration)
	// BlockUntil will block until the FakeClock has the given number of
	// sleepers (callers of Sleep or After)
	BlockUntil(n int)
}
```


在测试环境,可以自由地前进或暂停时间。例如，你可以创建一个FakeClock，然后在其上Sleep，并确保不会影响实际的时间流逝。更进一步，通过调用Advance方法，你可以立即跳转到未来的任意时刻，以便观察和验证你的代码如何响应时间的变化。

在生产环境中，你只需简单地切换回clockwork.NewRealClock()，就可以恢复使用真实的系统时间，无需任何修改


## 适合用于测试以下场景

- 定时任务: 当你需要测试定时器、调度器或者依赖于时间间隔的任何函数时。
- 超时处理: 测试函数或服务是否能在预期时间内正确响应或超时。
- 并发测试: 控制时间可以帮助你在多线程或goroutine环境中同步执行和检查。


## 第三方使用-->etcd

```go
// https://github.com/etcd-io/etcd/blob/34bd797e6754911ee540e8c87f708f88ffe89f37/etcdserver/api/v2discovery/discovery.go
type discovery struct {
	lg      *zap.Logger
	cluster string
	id      types.ID
	c       client.KeysAPI
	retries uint
	url     *url.URL

	clock clockwork.Clock
}


// 处理化
func newDiscovery(lg *zap.Logger, durl, dproxyurl string, id types.ID) (*discovery, error) {
	// ...
	return &discovery{
		lg:      lg,
		cluster: token,
		c:       dc,
		id:      id,
		url:     u,
		clock:   clockwork.NewRealClock(),
	}, nil
}
```

测试

```go
func TestWaitNodes(t *testing.T) {
	all := []*client.Node{
		0: {Key: "/1000/1", CreatedIndex: 2},
		1: {Key: "/1000/2", CreatedIndex: 3},
		2: {Key: "/1000/3", CreatedIndex: 4},
	}

	tests := []struct {
		nodes []*client.Node
		rs    []*client.Response
	}{
		{
			all,
			[]*client.Response{},
		},
		{
			all[:1],
			[]*client.Response{
				{Node: &client.Node{Key: "/1000/2", CreatedIndex: 3}},
				{Node: &client.Node{Key: "/1000/3", CreatedIndex: 4}},
			},
		},
		{
			all[:2],
			[]*client.Response{
				{Node: &client.Node{Key: "/1000/3", CreatedIndex: 4}},
			},
		},
		{
			append(all, &client.Node{Key: "/1000/4", CreatedIndex: 5}),
			[]*client.Response{
				{Node: &client.Node{Key: "/1000/3", CreatedIndex: 4}},
			},
		},
	}

	for i, tt := range tests {
		// Basic case
		c := &clientWithResp{rs: nil, w: &watcherWithResp{rs: tt.rs}}
		dBase := &discovery{cluster: "1000", c: c}

		// Retry case
		var retryScanResp []*client.Response
		if len(tt.nodes) > 0 {
			retryScanResp = append(retryScanResp, &client.Response{
				Node: &client.Node{
					Key:   "1000",
					Value: strconv.Itoa(3),
				},
			})
			retryScanResp = append(retryScanResp, &client.Response{
				Node: &client.Node{
					Nodes: tt.nodes,
				},
			})
		}
		cRetry := &clientWithResp{
			rs: retryScanResp,
			w:  &watcherWithRetry{rs: tt.rs, failTimes: 2},
		}
		fc := clockwork.NewFakeClock()
		dRetry := &discovery{
			cluster: "1000",
			c:       cRetry,
			clock:   fc,
		}

		for _, d := range []*discovery{dBase, dRetry} {
			go func() {
				for i := uint(1); i <= maxRetryInTest; i++ {
					fc.BlockUntil(1)
					fc.Advance(time.Second * (0x1 << i))
				}
			}()
			g, err := d.waitNodes(tt.nodes, 3, 0) // we do not care about index in this test
			if err != nil {
				t.Errorf("#%d: err = %v, want %v", i, err, nil)
			}
			if !reflect.DeepEqual(g, all) {
				t.Errorf("#%d: all = %v, want %v", i, g, all)
			}
		}
	}
}
```

## 参考
