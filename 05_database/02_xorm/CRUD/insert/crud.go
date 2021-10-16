package main

import (
	"fmt"
	"go_grpc_example/05_database/02_xorm/models"
	"go_grpc_example/05_database/02_xorm/util"
	// 注意引入，否则会空指针
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	var eg = util.GetEngineGroup()
	// 获取数据库表的结构信息
	schemeTables, _ := eg.DBMetas()
	fmt.Println("表的数量", len(schemeTables))
	for _, tableInfo := range schemeTables {
		fmt.Printf("%+v\n", *tableInfo)
	}

	// 添加数据
	// 插入一条
	var data = &models.MasterSlaveTable{
		Description: "我是master数据库6",
		Name:        "master6",
		// 不填写created ，默认为0
		//UpdatedAt: time.Now(), // 不填默认是空
	}
	affected, err := eg.Insert(data)
	if err != nil {
		fmt.Println("错误信息是", err)
	}
	fmt.Println("返回Id", data.Id)
	fmt.Println("影响的行数", affected)
	// 插入多条 数据库支持批量插入
	/*
		批量插入会自动生成Insert into table values (),(),()的语句，因此各个数据库对SQL语句有长度限制，
		因此这样的语句有一个最大的记录数，根据经验测算在150条左右。大于150条后，生成的sql语句将太长可能导致执行失败。
		因此在插入大量数据时，目前需要自行分割成每150条插入一次
	*/
	//multiData := make([]*models.MasterSlaveTable, 2)
	//multiData[0] = new(models.MasterSlaveTable)
	//multiData[0].Name = "name6"
	//multiData[1] = new(models.MasterSlaveTable)
	//multiData[1].Name = "name7"
	//affected2, err := eg.Insert(&multiData)
	//if err != nil {
	//	fmt.Println("错误信息是", err)
	//}
	//fmt.Println("返回影响的行号", affected2)
	//// 批量插入不返回字段的Id
	//fmt.Println("返回Id", multiData[0].Id)
	//fmt.Println("返回Id", multiData[1].Id)

	// 查询条件
	//var data2 = new(masterSlaveTable)
	//_, err = eg.Alias("o").Where("o.usr_name=?", "从").Get(data2)
	//_, err = eg.Where("usr_name=?", "从").Get(data2)
	//
	//if err != nil {
	//	fmt.Println("错误信息是", err)
	//}
	//fmt.Printf("返回data2:%+v", *data2)

	// 排序
	//dataSlice := make([]masterSlaveTable, 0)
	////eg.Asc("id").Find(&dataSlice) // 升序
	//_ = eg.Desc("id").Limit(1, 0).Find(&dataSlice) // 降序序限制一条
	//fmt.Printf("返回data2:%+v", dataSlice)

}
