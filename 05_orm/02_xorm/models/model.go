package models

import (
	"time"
)

type MasterSlaveTable2 struct {
	Id          int64     `xorm:"id notnull pk autoincr" `   // 如果field名称为Id而且类型为int64并且没有定义tag，则会被xorm视为主键，并且拥有自增属性。
	Description string    `xorm:"description comment('描述')"` // string类型默认映射为varchar(255)
	Name        string    `xorm:"'usr_name' notnull varchar(25) comment('用户名')" `
	CreatedAt   time.Time `xorm:"'created'"` // 注意双引号里面加单引号
	UpdatedAt   time.Time `xorm:"'updated'"`
}

// 自定义表名
func (m *MasterSlaveTable2) TableName() string {
	return "tableName"
}

type MasterSlaveTable struct {
	Id          int64     `xorm:"id notnull pk autoincr" `   // 如果field名称为Id而且类型为int64并且没有定义tag，则会被xorm视为主键，并且拥有自增属性。
	Description string    `xorm:"description comment('描述')"` // string类型默认映射为varchar(255)
	Name        string    `xorm:"'usr_name' notnull varchar(25) comment('用户名')" `
	CreatedAt   int64     `xorm:"'created'"` // 注意双引号里面加单引号
	UpdatedAt   time.Time `xorm:"'updated'"`
}

type ServerInfo struct {
	ServerInfoId       string    `xorm:"varchar(32) pk server_info_id"`
	CreatedAt          LocalTime `xorm:"timestamp created"`
	UpdatedAt          LocalTime `xorm:"timestamp updated"`
	DeletedAt          LocalTime `xorm:"timestamp deleted index"`
	OrgId              string    `xorm:"varchar(100) org_id" json:"orgId"`
	ServerIp           string    `xorm:"varchar(128) server_ip" json:"serverIp"`
	ServerNameDesc     string    `xorm:"varchar(500) server_name_desc" json:"serverNameDesc"`
	ServerTimeNow      LocalTime `xorm:"timestamp server_time" json:"serverTime"`
	DataReceiveTime    LocalTime `xorm:"timestamp data_receive_time" sql:"DEFAULT:current_timestamp" json:"dataRecvTime"`
	LastUploadDataTime LocalTime `xorm:"timestamp last_upload_data_time" json:"lastUploadDataTime"`
	LastCheckTime      LocalTime `xorm:"timestamp last_check_time" json:"lastCheckTime"`
	LastErrorTime      LocalTime `xorm:"timestamp last_error_time" json:"lastErrorTime"`
}

//既有LocalTime类型的，又有*LocalTime类型的，*LocalTime是考虑到有时候数据值可能为NULL，即字段值可能为空的情况。
