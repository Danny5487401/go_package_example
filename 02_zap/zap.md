
#zap源码分析
![](zap_structure.png)
![](zap_structure2.png)
通过 zap 打印一条结构化的日志大致包含5个过程：

    1.分配日志 Entry: 创建整个结构体，此时虽然没有传参(fields)进来，但是 fields 参数其实创建了
    
    2.检查级别，添加core: 如果 logger 同时配置了 hook，则 hook 会在 core check 后把自己添加到 cores 中
    
    3.根据选项添加 caller info 和 stack 信息: 只有大于等于级别的日志才会创建checked entry
    
    4.Encoder 对 checked entry 进行编码: 创建最终的 byte slice，将 fields 通过自己的编码方式(append)编码成目标串
    
    5.Write 编码后的目标串，并对剩余的 core 执行操作， hook也会在这时被调用



##日志
日志有两个概念：字段和消息。字段用来结构化输出错误相关的上下文环境，而消息简明扼要的阐述错误本身
```go
//用户不存在的错误消息可以这么打印
log.Error("User does not exist", zap.Int("uid", uid))
```
User does not exist 是消息， 而 uid 是字段
