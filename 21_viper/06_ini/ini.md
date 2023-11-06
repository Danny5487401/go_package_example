<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [ini (Initialization file 初始化文件)](#ini-initialization-file-%E5%88%9D%E5%A7%8B%E5%8C%96%E6%96%87%E4%BB%B6)
  - [格式简介](#%E6%A0%BC%E5%BC%8F%E7%AE%80%E4%BB%8B)
    - [parameter](#parameter)
    - [sections](#sections)
    - [comments](#comments)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# ini (Initialization file 初始化文件)

ini 是 Windows 上常用的配置文件格式。MySQL 的 Windows 版就是使用 ini 格式存储配置的

## 格式简介

案例
```ini
; 通用配置,文件后缀.ini
[common]

application.directory = APPLICATION_PATH  "/application"
application.dispatcher.catchException = TRUE


; 数据库配置
resources.database.master.driver = "pdo_mysql"
resources.database.master.hostname = "127.0.0.1"
resources.database.master.port = 3306
resources.database.master.database = "database"
resources.database.master.username = "username"
resources.database.master.password = "password"
resources.database.master.charset = "UTF8"


; 生产环境配置
[product : common]

; 开发环境配置
[develop : common]

resources.database.slave.driver = "pdo_mysql"
resources.database.slave.hostname = "127.0.0.1"
resources.database.slave.port = 3306
resources.database.slave.database = "test"
resources.database.slave.username = "root"
resources.database.slave.password = "123456"
resources.database.slave.charset = "UTF8"

; 测试环境配置
[test : common]

```

### parameter
INI所包含的最基本的“元素”就是parameter；每一个parameter都有一个name和一个value，如下所示：
```ini
name = value
```

### sections

所有的parameters都是以sections为单位结合在一起的。所有的section名称都是独占一行，并且sections名字都被方括号包围着（[ and ])。
在section声明后的所有parameters都是属于该section。对于一个section没有明显的结束标志符，一个section的开始就是上一个section的结束，或者是end of the file。

### comments
在INI文件中注释语句是以分号“；”开始的。所有的所有的注释语句不管多长都是独占一行直到结束的。在分号和行结束符之间的所有内容都是被忽略的。
```ini
;comments text

```


