package rmq

import (
	"encoding/json"
	"errors"
	"github.com/streadway/amqp"
	"io/ioutil"
	"os"
)

var _Cfg *tCfg = new(tCfg)                                                       //配置文件对象
var _ConnectPool map[string]*amqp.Connection = make(map[string]*amqp.Connection) //连接名称:连接对象
var _ChannelPool map[string]*amqp.Channel = make(map[string]*amqp.Channel)       //信道名称:信道对象
var _ExchangePool map[string]string = make(map[string]string)                    //交换机名称:所属信道名称
var _QueuePool map[string]string = make(map[string]string)                       //队列名称:所属信道名称
var _Pusher map[string]Pusher = make(map[string]Pusher)                          //Pusher名称:Pusher配置
var _Poper map[string]Poper = make(map[string]Poper)                             //Poper名称:Poper配置

//初始化
func Init(path string) (err error) {
	if err = loadCfg(path); err != nil {
		return err
	}
	if err = initConnect(); err != nil {
		return err
	}
	if err = initChannel(); err != nil {
		return err
	}
	if err = initExchange(); err != nil {
		return err
	}
	if err = initQueue(); err != nil {
		return err
	}
	if err = initPusher(); err != nil {
		return err
	}
	if err = initPoper(); err != nil {
		return err
	}
	return err
}

//读取配置文件
func loadCfg(path string) (err error) {
	var fp *os.File
	if fp, err = os.Open(path); err != nil {
		return err
	}
	var data []byte
	if data, err = ioutil.ReadAll(fp); err != nil {
		return err
	}
	if err = fp.Close(); err != nil {
		return err
	}
	if err = json.Unmarshal(data, _Cfg); err != nil {
		return err
	}
	return nil
}

//初始化信道
func initChannel() (err error) {
	for _, v := range _Cfg.Channels {
		if err = CreateChannel(v); err != nil {
			return err
		}
	}
	return nil
}

//创建信道
func CreateChannel(v Channel) (err error) {
	if _, ok := _ConnectPool[v.Connect]; !ok {
		return errors.New("连接不存在\n")
	}
	var channel *amqp.Channel
	if channel, err = _ConnectPool[v.Connect].Channel(); err != nil {
		return err
	} else {
		if _, ok := _ChannelPool[v.Name]; !ok {
			_ChannelPool[v.Name] = channel
		} else {
			return errors.New("信道已存在\n")
		}

	}
	if err = channel.Qos(v.QosCount, v.QosSize, false); err != nil {
		return nil
	}
	return nil
}

//初始化交换机
func initExchange() (err error) {
	for _, v := range _Cfg.Exchanges {
		if err = CreateExchange(v); err != nil {
			return err
		}
	}
	return nil
}

//创建交换机
func CreateExchange(v Exchange) (err error) {
	if _, ok := _ChannelPool[v.Channel]; !ok {
		return errors.New("信道不存在\n")
	}
	if err = _ChannelPool[v.Channel].ExchangeDeclare(v.Name, v.Type,
		v.Durable, v.AutoDeleted, v.Internal, v.NoWait, v.Args); err != nil {
		return err
	} else {
		if _, ok := _ExchangePool[v.Name]; !ok {
			_ExchangePool[v.Name] = v.Channel
		} else {
			return errors.New("交换机已存在")
		}
	}
	for _, b := range v.Bind {
		if err = _ChannelPool[v.Channel].ExchangeBind(b.Destination, b.Key, v.Name, b.NoWait, nil); err != nil {
			return err
		}
	}
	return nil
}

//初始化队列
func initQueue() (err error) {
	for _, v := range _Cfg.Queue {
		if err = CreateQueue(v); err != nil {
			return err
		}
	}
	return nil
}

//创建队列
func CreateQueue(v Queue) (err error) {
	if _, ok := _ChannelPool[v.Channel]; !ok {
		return errors.New("信道不存在\n")
	}

	//处理x-message-ttl的类型，json里面写的是int，go读出来的是double
	if _, ok := v.Args["x-message-ttl"]; ok {
		t := int32(v.Args["x-message-ttl"].(float64))
		delete(v.Args, "x-message-ttl")
		v.Args["x-message-ttl"] = t
	}
	// 申明队列
	if _, err = _ChannelPool[v.Channel].QueueDeclare(v.Name, v.Durable,
		v.AutoDelete, v.Exclusive, v.NoWait, v.Args); err != nil {
		return err
	} else {
		if _, ok := _QueuePool[v.Name]; !ok {
			_QueuePool[v.Name] = v.Channel
		} else {
			return errors.New("队列已存在")
		}
	}
	// 绑定交换机
	for _, b := range v.Bind {
		if err = _ChannelPool[v.Channel].QueueBind(v.Name, b.Key, b.ExchangeName, b.NoWait, nil); err != nil {
			return err
		}
	}
	return nil
}

//创建Pusher
func CreatePusher(v Pusher) (err error) {
	if _, ok := _Pusher[v.Name]; !ok {
		_Pusher[v.Name] = v
	} else {
		return errors.New("Pusher已存在")
	}
	return nil
}

//删除Pusher
func DeletePusher(name string) (err error) {
	if _, ok := _Poper[name]; ok {
		delete(_Poper, name)
	}
	return nil
}

//初始化Pusher
func initPusher() (err error) {
	for _, v := range _Cfg.Pusher {
		if err = CreatePusher(v); err != nil {
			return err
		}
	}
	return err
}

//创建Poper
func CreatePoper(v Poper) (err error) {
	if _, ok := _Poper[v.Name]; !ok {
		_Poper[v.Name] = v
	} else {
		return errors.New("Poper已存在")
	}
	return err
}

//删除Poper
func DeletePoper(name string) (err error) {
	if _, ok := _Poper[name]; ok {
		delete(_Poper, name)
	}
	return nil
}

//初始化Poper
func initPoper() (err error) {
	for _, v := range _Cfg.Poper {
		if err = CreatePoper(v); err != nil {
			return err
		}
	}
	return nil
}

//关闭
func Fini() (err error) {
	for _, conn := range _ConnectPool {
		for _, ch := range _ChannelPool {
			if err = ch.Close(); err != nil {
				return err
			}
		}
		if err = conn.Close(); err != nil {
			return err
		}
	}
	//清空所有缓存
	_Cfg = new(tCfg)                                 //配置文件对象
	_ConnectPool = make(map[string]*amqp.Connection) //连接名称:连接对象
	_ChannelPool = make(map[string]*amqp.Channel)    //信道名称:信道对象
	_ExchangePool = make(map[string]string)          //交换机名称:所属信道名称
	_QueuePool = make(map[string]string)             //队列名称:所属信道名称
	_Pusher = make(map[string]Pusher)                //Pusher名称:Pusher配置
	_Poper = make(map[string]Poper)                  //Poper名称:Poper配置

	return nil
}
