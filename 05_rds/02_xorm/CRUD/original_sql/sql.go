package main

import (
	"fmt"
	"go_grpc_example/05_rds/02_xorm/util"
	"time"

	// 注意引入，否则会空指针
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	var err error
	var eg = util.GetEngineGroup()

	// 原生查询
	// Note：Query方法在下面不行
	//rsp, err := eg.QueryString("select * from user_active")
	//if err != nil {
	//	return
	//}
	//for _, v := range rsp {
	//	fmt.Printf("解析前数据是%+v\n", v)
	//}
	//data := make([]models.UserActive, 0)
	//err = mapstructure.WeakDecode(rsp, &data)
	//if err != nil {
	//	fmt.Println("获取格式错误", err.Error())
	//	return
	//}
	//fmt.Printf("结构是%+v\n", data)

	// 原生插入数据
	uid := 10032100
	totalDays := 100
	createdAt := time.Now().Unix()
	// insert ignore会忽略数据库中已经存在的数据(根据主键或者唯一索引判断)，如果数据库没有数据，就插入新的数据，如果有数据的话就跳过这条数据.但是主键还是自动增长
	affected, err := eg.Exec("insert ignore into user_active(uid,total_days,created_at) values(?,?,?)", uid, totalDays, createdAt)
	if err != nil {
		fmt.Println("插入数据有误", err.Error())
		return
	}
	fmt.Println(affected.RowsAffected())
}
