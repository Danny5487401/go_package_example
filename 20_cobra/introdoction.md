# Cobra
Cobra 是一个非常实用(流行)的golang包，很多优秀的开源应用都在使用它，包括 Docker 和 Kubernetes 等，
它提供了简单的接口来创建命令行程序。同时，Cobra 也是一个应用程序，用来生成应用框架，从而开发以 Cobra 为基础的应用。
![](.introdoction_images/cobra_menu.png)  

在cobra中，所有的命令会组成一个树的结构，必然有一个根命令，我们应用的每次执行，都是从这个根命令开始的，官方文档也说过，基于cobra的应用的main包中的代码是很简单的，几乎没有额外的操作，仅有的操作其实就是执行我们的根命令。
## 主要功能
* 简易的子命令行模式，如 app server， app fetch 等等
* 完全兼容 posix 命令行模式
* 嵌套子命令 subcommand
* 支持全局，局部，串联 flags
* 使用 cobra 很容易的生成应用程序和命令，使用 cobra create - appname 和 cobra add cmdname
* 如果命令输入错误，将提供智能建议，如 app srver，将提示 srver 没有，是不是 app server
* 自动生成 commands 和 flags 的帮助信息
* 自动生成详细的 help 信息，如 app help
* 自动识别帮助 flag -h，--help
* 自动生成应用程序在 bash 下命令自动完成功能
* 自动生成应用程序的 man 手册
* 命令行别名
* 自定义 help 和 usage 信息
* 可选的与 viper apps 的紧密集成

## 基本操作
需求：
- （1）show 查看当前时间
- （2）parse 指定时间格式 --format，parse为show的子命令。

### 1. 通过命令初始化项目
```shell script
cobra init --pkg-name go_grpc_example/20_cobra
```
### 2. 通过命令生成动作
```shell script
cobra add serve
cobra add config
cobra add create -p 'configCmd'
```
cobra的三个概念：
* commands 行为
* arguments 数值
* flags 对行为的改变
执行命令行程序时的一般格式为:  
`[appName] [command] [arguments] --[flag]`  
  eg:
```shell script
cobra add show
# 终端返回： show created at /Users/python/Desktop/go_grpc_example/20_cobra

```
  
执行命令，我们可以看到time命令的帮助信息，在usage中我们可以看到，当前有两个command一个是默认的help另一个就是我们刚才创建的show，
Flags也有两个，这是初始化的时候自带的，后面我们可以自己进行修改.  
### 3. 为 Command 添加选项(flags)  

选项的作用范围，可以把选项分为两类：  
* persistent 既可以设置给该 Command，又可以设置给该 Command 的子 Command
* local  
```shell script
ls --help
Usage: ls [OPTION]... [FILE]...
```  
解释：其中的 OPTION 对应本文中介绍的 flags，以 - 或 – 开头；
而 FILE 则被称为参数(arguments)或位置参数。一般的规则是参数在所有选项的后面，上面的 … 表示可以指定多个选项和多个参数
#### 4. 添加子命令
```shell script
cobra add parse -p showCmd  
# 终端返回： parse created at /Users/python/Desktop/go_grpc_example/20_cobra
```
解释： 给showCmd命令添加子命令parse


#### 5. 编译
```shell script
$ go build -o time main.go
``` 


## 执行流程
```go
func Execute() {
  rootCmd.Execute()
}
func (c *Command) Execute() error {
	_, err := c.ExecuteC()
	return err
}
```