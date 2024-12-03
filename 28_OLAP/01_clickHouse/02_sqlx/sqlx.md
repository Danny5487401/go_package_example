<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [sqlx](#sqlx)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# sqlx

该库兼容sql原生包，在出色的内置database/sql软件包的基础上提供了一组扩展,提供了更为强大的、优雅的查询、插入函数。

该库提供四个处理类型，分别是：

- sqlx.DB – 类似原生的 sql.DB,表示数据库；
- sqlx.Tx – 类似原生的 sql.Tx,表示transaction；
- sqlx.Stmt – 类似原生的 sql.Stmt, t准备 SQL 语句(prepared statemen)操作；
- sqlx.NamedStmt – 对特定参数命名并绑定生成 SQL 语句操作。

提供两个游标类型，分别是：

- sqlx.Rows – 类似原生的 sql.Rows, 从 Queryx 返回；
- sqlx.Row  – 类似原生的 sql.Row, 从 QueryRowx 返回