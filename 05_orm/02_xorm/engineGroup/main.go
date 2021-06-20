package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

var eg *xorm.EngineGroup

func main() {
	var err error
	// 格式:sql.Open("mysql", "user:password@/dbname")
	master, err := xorm.NewEngine("mysql", "root:chuanzhi@tcp(106.14.35.115:3307)/masterSlaveDB")
	if err != nil {
		return
	}
	slave1, err := xorm.NewEngine("mysql", "root:123456@tcp(106.14.35.115:3308)/masterSlaveDB")
	if err != nil {
		return
	}
	slave2, err := xorm.NewEngine("mysql", "root:123456@tcp(106.14.35.115:3309)/masterSlaveDB")
	if err != nil {
		return
	}
	slaves := []*xorm.Engine{slave1, slave2}
	eg, err = xorm.NewEngineGroup(master, slaves)
	eg = eg.SetPolicy(xorm.RandomPolicy())
	//连接测试
	if err := eg.Ping(); err != nil {
		fmt.Println(err)
		return
	}
	eg.ShowSQL(true)
	//设置连接池的空闲数大小
	eg.SetMaxIdleConns(5)
	//设置最大打开连接数
	eg.SetMaxOpenConns(5)
	fmt.Println("连接成功")

}
