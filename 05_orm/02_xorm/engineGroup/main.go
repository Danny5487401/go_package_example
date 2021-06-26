package main

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"go_grpc_example/05_orm/02_xorm/model"
	"go_grpc_example/05_orm/02_xorm/util"
)

func main() {
	//var err error
	var eg = util.GetEngineGroup()
	// 获取数据库表的结构信息
	schemeTables, _ := eg.DBMetas()
	fmt.Println("表的数量", len(schemeTables))
	for _, tableInfo := range schemeTables {
		fmt.Printf("%+v\n", *tableInfo)
	}
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
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}

	// 添加数据
	//var data = masterSlaveTable{
	//	Description: "我是master数据库",
	//	Name:        "master",
	//	CreatedAt:   time.Now(),
	//	UpdatedAt:   JsonTime(time.Now()),
	//}
	//Id, err := eg.Insert(&data)
	//if err != nil {
	//	fmt.Println("错误信息是", err)
	//}
	//fmt.Println("返回Id", Id)

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

type masterSlaveTable struct {
	Id          int64     `xorm:"id notnull pk autoincr" `   // 如果field名称为Id而且类型为int64并且没有定义tag，则会被xorm视为主键，并且拥有自增属性。
	Description string    `xorm:"description comment('描述')"` // string类型默认映射为varchar(255)
	Name        string    `xorm:"'usr_name' notnull varchar(25) comment('用户名')" `
	CreatedAt   time.Time `xorm:"'created'"` // 注意双引号里面加单引号
	UpdatedAt   JsonTime  `xorm:"'updated'"`
}

// 自定义表名
func (m *masterSlaveTable) TableName() string {
	return "masterSlaveTable2"
}

type JsonTime time.Time

func (j JsonTime) MarshalJSON() ([]byte, error) {
	//当返回时间为空时，需特殊处理

	return []byte(`"` + time.Time(j).Format("2006-01-02 03:04:05") + `"`), nil
}

type ServerInfo struct {
	ServerInfoId       string            `xorm:"varchar(32) pk server_info_id"`
	CreatedAt          models.LocalTime  `xorm:"timestamp created"`
	UpdatedAt          models.LocalTime  `xorm:"timestamp updated"`
	DeletedAt          *models.LocalTime `xorm:"timestamp deleted index"`
	OrgId              string            `xorm:"varchar(100) org_id" json:"orgId"`
	ServerIp           string            `xorm:"varchar(128) server_ip" json:"serverIp"`
	ServerNameDesc     string            `xorm:"varchar(500) server_name_desc" json:"serverNameDesc"`
	ServerTimeNow      models.LocalTime  `xorm:"timestamp server_time" json:"serverTime"`
	DataReceiveTime    models.LocalTime  `xorm:"timestamp data_receive_time" sql:"DEFAULT:current_timestamp" json:"dataRecvTime"`
	LastUploadDataTime *models.LocalTime `xorm:"timestamp last_upload_data_time" json:"lastUploadDataTime"`
	LastCheckTime      *models.LocalTime `xorm:"timestamp last_check_time" json:"lastCheckTime"`
	LastErrorTime      *models.LocalTime `xorm:"timestamp last_error_time" json:"lastErrorTime"`
}

//既有LocalTime类型的，又有*LocalTime类型的，*LocalTime是考虑到有时候数据值可能为NULL，即字段值可能为空的情况。
