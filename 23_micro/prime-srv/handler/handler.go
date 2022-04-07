package handler

import (
	"context"
	"go_package_example/23_micro/prime-srv/service"

	"go_package_example/23_micro/proto/prime"
)

type handler struct {
}

func (h handler) GetPrime(ctx context.Context, req *prime.PrimeRequest, rsp *prime.PrimeResponse) error {
	inputs := make([]int64, 0)
	var i int64 = 0
	for ; i <= req.Input; i++ {
		inputs = append(inputs, i)
	}
	rsp.Output = service.GetPrime(inputs...)
	return nil

}

func Handler() prime.PrimeHandler {
	return handler{}
}
