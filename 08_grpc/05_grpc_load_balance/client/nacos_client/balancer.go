package nacos_client

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"math/rand"
	"sync"
)

// Weighted Round Robin实现
// 如果一个conn的权重为n，那么就在加权结果集中加入n个conn，这样在后续Pick时不需要考虑加权的问题，只需向普通Round Robin那样逐个Pick出来即可。
type wrrPickerBuilder struct{}

func (*wrrPickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	if len(info.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}

	var scs []balancer.SubConn
	// 提取已经就绪的connection的权重信息，作为Picker实例的输入
	for subConn, addr := range info.ReadySCs {
		weight := addr.Address.Attributes.Value("weight").(int)
		if weight <= 0 {
			weight = 1
		}
		for i := 0; i < weight; i++ {
			scs = append(scs, subConn)
		}
	}

	return &wrrPicker{
		subConns: scs,
		// Start at a random index, as the same RR balancer rebuilds a new
		// picker when SubConn states change, and we don't want to apply excess
		// load to the first server in the list.
		next: rand.Intn(len(scs)),
	}
}

type wrrPicker struct {
	// subConns is the snapshot of the roundrobin balancer when this picker was
	// created. The slice is immutable. Each Get() will do a round robin
	// selection from it and return the selected SubConn.
	subConns []balancer.SubConn

	mu   sync.Mutex
	next int
}

// 选出一个Connection
func (p *wrrPicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	p.mu.Lock()
	sc := p.subConns[p.next]
	p.next = (p.next + 1) % len(p.subConns)
	p.mu.Unlock()
	return balancer.PickResult{SubConn: sc}, nil
}
