// thrift --gen go example.thrift
// thrift version 0.22.0
namespace go thrift.example   // 定义所使用的命名空间

struct Person {             // 定义一个结构体
  1: required string name,  // 姓名字段
  2: optional i32 age       // 年龄字段
}

service ExampleService {    // 定义一个服务接口
  void sayHello(1: string name) // sayHello方法，接收一个姓名参数
}
