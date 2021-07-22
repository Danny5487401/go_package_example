package main

import (
	"fmt"
	"go_grpc_example/05_orm/02_xorm/models"

	_ "github.com/go-sql-driver/mysql"

	"go_grpc_example/05_orm/02_xorm/util"
)

func main() {
	var err error
	var eg = util.GetEngineGroup()
	// 获取数据库表的结构信息
	//schemeTables, _ := eg.DBMetas()
	//fmt.Println("表的数量", len(schemeTables))
	//for _, tableInfo := range schemeTables {
	//	fmt.Printf("%+v\n", *tableInfo)
	//}
	// 自己构建表结构信息
	//masterTableInfo := new(masterSlaveTable)
	//table, err := eg.TableInfo(masterTableInfo)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Printf("%+v\n", *table)
	// 创建表
	//eg.Charset("utf8")
	////eg.StoreEngine("ISAM")
	//err = eg.CreateTables(masterSlaveTable{})
	//err = eg.CreateTables(ServerInfo{})
	// 优先级Table()最大
	// 方式一
	//err = eg.Table("table").CreateTable(models.MasterSlaveTable2{})
	// 方式二
	//err = eg.CreateTables(models.MasterSlaveTable{}, models.MasterSlaveTable2{}, models.ServerInfo{})
	// 方式三 推荐方式 有创建索引
	err = eg.Sync2(models.MasterSlaveTable{})

	if err != nil {
		fmt.Println(err)
		return
	}

}
