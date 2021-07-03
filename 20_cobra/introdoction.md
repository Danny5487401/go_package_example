# 初始化项目
```shell script
cobra init --pkg-name go_grpc_example/20_cobra
```
## 想要完成的功能
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
添加命令  
```shell script
cobra add show
show created at /Users/python/Desktop/go_grpc_example/20_cobra

```
```shell script
$ go build -o time main.go
```   
执行命令，我们可以看到time命令的帮助信息，在usage中我们可以看到，当前有两个command一个是默认的help另一个就是我们刚才创建的show，
Flags也有两个，这是初始化的时候自带的，后面我们可以自己进行修改.  
# 为 Command 添加选项(flags)  
选项的作用范围，可以把选项分为两类：  
* persistent 既可以设置给该 Command，又可以设置给该 Command 的子 Command
* local  
```shell script
ls --help
Usage: ls [OPTION]... [FILE]...
```  
其中的 OPTION 对应本文中介绍的 flags，以 - 或 – 开头；而 FILE 则被称为参数(arguments)或位置参数。一般的规则是参数在所有选项的后面，上面的 … 表示可以指定多个选项和多个参数
# 添加子命令
```shell script
cobra add parse -p showCmd
parse created at /Users/python/Desktop/go_grpc_example/20_cobra
```