package main

import (
	"fmt"
	"go_grpc_example/05_rds/02_xorm/models"
	"go_grpc_example/05_rds/02_xorm/util"
	// 注意引入，否则会空指针
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	var err error
	var eg = util.GetEngineGroup()

	// 更新数据

	// 结构体方式
	var data5 = new(models.MasterSlaveTable)
	data5.Description = "我是Joy哥哥"

	// Update方法将返回两个参数，第一个为 更新的记录数
	//eg.Cols("description").Update(data5) // 这是更新所有数据的Description字段
	var affected int64
	affected, err = eg.Where("usr_name=?", "joy").Cols("description").Update(data5)
	if err != nil {
		fmt.Println("选择部分字段错误", err.Error())
		return
	}
	fmt.Println("影响的行数", affected)

	// map方式

}
