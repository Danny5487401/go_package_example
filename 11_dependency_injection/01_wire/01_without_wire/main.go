package main

import (
	"fmt"
)

/*
未使用依赖注入
*/

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

// 开始调用
func main() {
	// 开始一大堆的初始化
	message := NewMessage()
	greeter := NewGreeter(message)
	event := NewEvent(greeter)

	event.Start()

}
