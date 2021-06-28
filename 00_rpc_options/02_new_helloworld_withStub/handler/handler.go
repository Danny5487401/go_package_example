package handler

// 名称冲突
const HelloServiceName = "handler/HelloService"

// 我们关心的不是HelloService名字，而是结构体中的方法
// 解耦：鸭子类型 接口
type HelloService struct {
}

func (s *HelloService) Hello(request string, reply *string) error {
	*reply = "hello, " + request
	return nil
}
