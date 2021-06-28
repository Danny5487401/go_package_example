package rmq

import (
	"errors"
	"github.com/streadway/amqp"
)

type MSG struct {
	Body    []byte
	Tag     uint64
	Channel string
	Poper   string
}

func (m MSG) Ack(multiple bool) (err error) {
	if _, ok := _ChannelPool[m.Channel]; !ok {
		return errors.New("Ack失败,Channel无效")
	} else {
		// 手动调用ack
		_ChannelPool[m.Channel].Ack(m.Tag, multiple)
	}
	return nil
}

//向交换机推送一条消息
func Push(name string, key string, msg []byte) (err error) {
	if _, ok := _Pusher[name]; !ok {
		return errors.New("pusher不存在")
	}

	cfg := _Pusher[name]
	if key != "" {
		// 路由的key
		cfg.Key = key
	}
	if _, ok := _ChannelPool[cfg.Channel]; !ok {
		return errors.New("channel不存在")
	}
	if err = _ChannelPool[cfg.Channel].Publish(cfg.Exchange, cfg.Key, cfg.Mandtory, cfg.Immediate,
		amqp.Publishing{ContentType: cfg.ContentType, Body: msg}); err != nil {
		return err
	}
	return nil
}

//从队列获取消息 -- 推模式
func Pop(name string, callback func(MSG)) (err error) {
	if _, ok := _Poper[name]; !ok {
		return errors.New("Poper不存在")
	}
	cfg := _Poper[name]
	if _, ok := _ChannelPool[cfg.Channel]; !ok {
		return errors.New("Channel不存在")
	}
	var msgs <-chan amqp.Delivery
	if msgs, err = _ChannelPool[cfg.Channel].Consume(cfg.QName, cfg.Consumer,
		cfg.AutoACK, cfg.Exclusive, cfg.NoLocal, cfg.NoWait, nil); err != nil {
		return err
	}
	go handleMsg(msgs, callback, cfg.Channel, name)

	return nil
}

//处理消息(顺序处理,如果需要多线程可以在回调函数中做手脚)
func handleMsg(msgs <-chan amqp.Delivery, callback func(MSG), channel string, poperName string) {
	for d := range msgs {
		var msg MSG = MSG{
			Body:    d.Body,
			Tag:     d.DeliveryTag,
			Channel: channel,
			Poper:   poperName,
		}
		callback(msg)
	}
}
