//+build wireinject

// 注意以上需要空一行
package wire

import (
	"fmt"
	"github.com/google/wire" //引入wire包
)

// 声明injector注入器
func InitializeEvent() Event {
	wire.Build(NewEvent, NewGreeter, NewMessage)

	// 返回零值给编译器，即使加了值，编译器也会忽略
	return Event{}
}

type Message string

type Greeter struct {
	Message Message // <- adding a Message field
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
	return Greeter{Message: m}
}

// 绑定欢迎
func (g Greeter) Greet() Message {
	return g.Message
}

// 创建一场宴会
func NewEvent(g Greeter) Event {
	return Event{Greeter: g}
}

// 宴会开始
func (e Event) Start() {
	msg := e.Greeter.Greet()
	fmt.Println(msg)
}
