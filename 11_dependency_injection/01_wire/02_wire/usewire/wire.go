//+build wireinject

// 注意以上需要空一行

package wire

import (
	"fmt"
	"github.com/google/wire"
)

// 1.声明injector注入器
//Injector 就是你最终想要的结果——最终的 App 对象的初始化函数
/*
func 来一袋垃圾食品() 一袋垃圾食品 {
    panic(wire.Build(来一份巨无霸套餐, 来一份双层鳕鱼堡套餐, 来一盒麦乐鸡, 垃圾食品打包))
}
*/
func InitializeEvent() Event {
	// 方式一
	//wire.Build(NewEvent, NewGreeter, NewMessage)
	// 方式二
	wire.Build(SomeProviderSet)

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

// 2.实现各个 Provider
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

/*
3。wire 里面还有个 ProviderSet 的概念，就是把一组 Provider 打包，因为通常你点单的时候很懒，不想这样点你的巨无霸套餐：我要一杯可乐，一包薯条，
一个巨无霸汉堡；你想直接戳一下就好了，来一份巨无霸套餐。这个套餐就是 ProviderSet，一组约定好的配方，
不然你的点单列表（injector 里的 Build）就会变得超级长，这样你很麻烦，服务员看着也很累
// 先定义套餐内容
var 巨无霸套餐 = wire.NewSet(来一杯可乐，来一包薯条，来一个巨无霸汉堡)

// 然后实现各个食品的做法
func 来一杯可乐() 一杯可乐 {}
func 来一包薯条() 一包薯条 {}
func 来一个巨无霸汉堡() 一个巨无霸汉堡 {}

*/
// 方式二
var SomeProviderSet = wire.NewSet(NewEvent, NewGreeter, NewMessage)
