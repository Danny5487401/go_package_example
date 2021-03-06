package util

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql" //别忘记倒入
	"xorm.io/xorm"
)

var eg *xorm.EngineGroup

func GetEngineGroup() *xorm.EngineGroup {
	if eg == nil {
		initEngine()
	}
	return eg
}

func initEngine() {

	// 格式:sql.Open("mysql", "user:password@/dbname")
	master, err := xorm.NewEngine("mysql", "root:chuanzhi@tcp(ali.danny.games:3307)/masterSlaveDB")
	if err != nil {
		return
	}
	slave1, err := xorm.NewEngine("mysql", "root:123456@tcp(ali.danny.games:3308)/masterSlaveDB")
	if err != nil {
		return
	}
	slave2, err := xorm.NewEngine("mysql", "root:123456@tcp(ali.danny.games:3309)/masterSlaveDB")
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
	eg.ShowSQL(true) // 调试显示sql语句
	//设置连接池的空闲数大小
	eg.SetMaxIdleConns(5)
	//设置最大打开连接数
	eg.SetMaxOpenConns(5)
	fmt.Println("连接成功")
}
