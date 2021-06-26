package handler

import (
	"go_grpc_example/23_micro/proto/sum"
	"go_grpc_example/23_micro/sum-srv/service"

	"context"
)

type handler struct {
}

func (h handler) GetSum(ctx context.Context, req *sum.SumRequest, rsp *sum.SumResponse) error {
	inputs := make([]int64, 0)
	var i int64 = 0
	for ; i <= req.Input; i++ {
		inputs = append(inputs, i)
	}
	rsp.Output = service.GetSum(inputs...)
	return nil
}

// 私有处理
func Handler() sum.SumHandler {
	return handler{}
}
