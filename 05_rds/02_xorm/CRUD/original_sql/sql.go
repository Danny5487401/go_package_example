package main

import (
	"fmt"
	"go_grpc_example/05_rds/02_xorm/util"
	// 注意引入，否则会空指针
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	var err error
	var eg = util.GetEngineGroup()

	// 原生查询
	rsp, err := eg.Query("select * from user_active")
	if err != nil {
		return
	}
	for _, v := range rsp {
		fmt.Printf("解析前数据是%+v\n", v)
	}

	// 原生插入数据
	affected, err := eg.Exec("insert into user_active(id,uid) values(11,100321346)")
	if err != nil {
		fmt.Println("插入数据有误", err.Error())
		return
	}
	fmt.Println(affected.RowsAffected())
}
