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

	// 获取数据库表的结构信息
	schemeTables, _ := eg.DBMetas()
	fmt.Println("表的数量", len(schemeTables))
	for _, tableInfo := range schemeTables {
		fmt.Printf("%+v\n", *tableInfo)
	}

	masterTableInfo := new(masterSlaveTable)
	table, err := eg.TableInfo(masterTableInfo)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\n", *table)
	// 创建表
	eg.Charset("utf8")
	eg.StoreEngine("ISAM")
	err = eg.CreateTables(masterTableInfo)
	if err != nil {
		fmt.Println(err)
		return
	}

}

type masterSlaveTable struct {
	Id          int64  // 如果field名称为Id而且类型为int64并且没有定义tag，则会被xorm视为主键，并且拥有自增属性。
	description string `xorm:"description comment('描述')"` // string类型默认映射为varchar(255)
	name        string `xorm:"'usr_name' notnull varchar(25)" `
}

func (m *masterSlaveTable) TableName() string {
	return "masterSlaveTable2"
}
