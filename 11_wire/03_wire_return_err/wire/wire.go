//+build wireinject

// 注意以上需要空一行
/*
初始化对象带error返回
*/
package wire

import (
	"errors"
	"fmt"
	"github.com/google/wire" //引入wire包
	"time"
)

// 声明injector注入器
func InitializeEvent() (Event, error) {
	wire.Build(NewEvent, NewGreeter, NewMessage)
	return Event{}, nil
}

type Message string

type Greeter struct {
	Message Message // <- adding a Message field
	Grumpy  bool
}

// 宴会
type Event struct {
	Greeter Greeter // <- adding a Greeter field
}

//创建消息
func NewMessage() Message {
	return Message("Hi there!")
}

// 创建招待人
func NewGreeter(m Message) Greeter {
	var grumpy bool
	if time.Now().Unix()%2 == 0 {
		grumpy = true
	}
	return Greeter{Message: m, Grumpy: grumpy}
}

// 绑定欢迎方法
func (g Greeter) Greet() Message {
	if g.Grumpy {
		return Message("Go away!")
	}
	return g.Message
}

// 创建一场宴会  带错误返回
func NewEvent(g Greeter) (Event, error) {
	if g.Grumpy {
		return Event{}, errors.New("could not create event: event greeter is grumpy暴躁的")
	}
	return Event{Greeter: g}, nil
}

// 宴会开始
func (e Event) Start() {
	msg := e.Greeter.Greet()
	fmt.Println(msg)
}
