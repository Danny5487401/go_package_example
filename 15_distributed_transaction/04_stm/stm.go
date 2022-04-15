package _4_stm

import (
	"fmt"
	"github.com/spf13/cast"
	v3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

func txnStmTransfer(cli *v3.Client, from, to string, amount uint64) error {
	// NewSTM 创建了一个原子事务的上下文，业务代码作为一个函数传进去
	_, err := concurrency.NewSTM(cli, func(stm concurrency.STM) error {
		// stm.Get 封装了事务的读操作
		senderNum := toUint64(stm.Get(from))
		receiverNum := toUint64(stm.Get(to))
		if senderNum < amount {
			return fmt.Errorf("余额不足")
		}
		// 事务的写操作
		stm.Put(to, fromUint64(receiverNum+amount))
		stm.Put(from, fromUint64(senderNum-amount))
		return nil
	})
	return err
}

func fromUint64(u uint64) string {
	return cast.ToString(u)
}

func toUint64(get string) uint64 {
	return cast.ToUint64(get)
}
