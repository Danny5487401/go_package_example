package main

import (
	"fmt"
	"go_package_example/05_rds/02_xorm/models"
	"go_package_example/05_rds/02_xorm/util"
	// 注意引入，否则会空指针
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	var err error
	var eg = util.GetEngineGroup()

	//查询数据
	//查询条件
	var data2 = new(models.MasterSlaveTable)
	//_, err = eg.Alias("o").Where("o.usr_name=?", "master").Get(data2)
	_, err = eg.Where("usr_name= ?", "master").Get(data2)

	if err != nil {
		fmt.Println("错误信息是", err)
	}
	fmt.Printf("返回data2:%+v", *data2)

	// 排序
	dataSlice := make([]*models.MasterSlaveTable, 0)
	//eg.Asc("id").Find(&dataSlice) // 升序
	_ = eg.Desc("id").Limit(2, 0).Find(&dataSlice) // 返回多条数据按id降序 限制一条
	fmt.Printf("返回dataSlice:%+v\n", dataSlice)

	// 选择部分字段
	var data3 = new(models.MasterSlaveTable)
	data3.Name = "master1"
	ok, err := eg.Select("description,updated").Get(data3) //返回一条数据
	if err != nil {
		fmt.Println("选择部分字段错误", err.Error())
	}
	if ok {
		fmt.Printf("返回data3部分字段:%+v\n", *data3)
	}
	// iterate 迭代--直接多单条数据
	var data4 = new(models.MasterSlaveTable)
	eg.Select("description,updated").Iterate(data4, func(idx int, bean interface{}) error {
		user := bean.(*models.MasterSlaveTable)
		fmt.Println("data4数据", user)
		return nil
	})

	// 选中id为2，3的数据
	dataSlice1 := make([]models.MasterSlaveTable, 0)
	eg.In("id", "2", 3).Find(&dataSlice1) // 参数可以不同类型
	fmt.Printf("返回dataSlice1:%+v\n", dataSlice1)

}
