<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [sqlx是 Go 的软件包，它在出色的内置database/sql软件包的基础上提供了一组扩展。](#sqlx%E6%98%AF-go-%E7%9A%84%E8%BD%AF%E4%BB%B6%E5%8C%85%E5%AE%83%E5%9C%A8%E5%87%BA%E8%89%B2%E7%9A%84%E5%86%85%E7%BD%AEdatabasesql%E8%BD%AF%E4%BB%B6%E5%8C%85%E7%9A%84%E5%9F%BA%E7%A1%80%E4%B8%8A%E6%8F%90%E4%BE%9B%E4%BA%86%E4%B8%80%E7%BB%84%E6%89%A9%E5%B1%95)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# sqlx是 Go 的软件包，它在出色的内置database/sql软件包的基础上提供了一组扩展。

该库兼容sql原生包，同时又提供了更为强大的、优雅的查询、插入函数。

该库提供四个处理类型，分别是：

- sqlx.DB – 类似原生的 sql.DB,表示数据库；
- sqlx.Tx – 类似原生的 sql.Tx,表示transaction；
- sqlx.Stmt – 类似原生的 sql.Stmt, t准备 SQL 语句(prepared statemen)操作；
- sqlx.NamedStmt – 对特定参数命名并绑定生成 SQL 语句操作。
提供两个游标类型，分别是：

- sqlx.Rows – 类似原生的 sql.Rows, 从 Queryx 返回；
- sqlx.Row  – 类似原生的 sql.Row, 从 QueryRowx 返回