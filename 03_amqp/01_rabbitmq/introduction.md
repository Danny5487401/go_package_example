# docker安装
```shell script
docker run -d --hostname my-rabbit --name rmq -p 15672:15672 -p 5672:5672 -p 25672:25672 -e RABBITMQ_DEFAULT_USER=用户名 -e RABBITMQ_DEFAULT_PASS=密码 rabbitmq:3-management
```

通过命令可以看出，一共映射了三个端口，简单说下这三个端口是干什么的。
5672：连接生产者、消费者的端口。
15672：WEB管理页面的端口。
25672：分布式集群的端口

借助备用交换机、TTL+DLX代替mandatory、immediate方案：
1、P发送msg给Ex，Ex无法把msg路由到Q，则会把路由转发给ErrEx。
2、msg暂存在Q上之后，如果C不能及时消费msg，则msg会转发到DlxEx。
3、TTL为msg在Q上的暂存时间，单位为毫秒。

通过设置参数，可以设置Ex的备用交换器ErrEx
创建Exchange时，指定Ex的Args – “alternate-exchange”:”ErrEx”。
其中ErrEx为备用交换器名称

通过设置参数，可以设置Q的DLX交换机DlxEX
创建Queue时，指定Q的Args参数：
“x-message-ttl”:0 //msg超时时间，单位毫秒
“x-dead-letter-exchange”:”dlxExchange” //DlxEx名称
“x-dead-letter-routing-key”:”dlxQueue” //DlxEx路由键
