package rmq

import (
	"errors"
	"github.com/streadway/amqp"
)

//初始化连接
func initConnect() (err error) {
	for _, v := range _Cfg.Connects {
		if err = CreateConnect(v); err != nil {
			return err
		}
	}
	return nil
}

//创建连接
func CreateConnect(v Connect) (err error) {
	var connect *amqp.Connection
	if connect, err = amqp.Dial(v.Addr); err != nil {
		return err
	} else {
		if _, ok := _ConnectPool[v.Name]; !ok {
			_ConnectPool[v.Name] = connect
		} else {
			return errors.New("连接已存在\n")
		}
	}
	return nil
}
