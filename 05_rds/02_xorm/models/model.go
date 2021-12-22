package models

import (
	"time"
)

type MasterSlaveTable2 struct {
	Id          int64     `xorm:"id notnull pk autoincr" `   // 如果field名称为Id而且类型为int64并且没有定义tag，则会被xorm视为主键，并且拥有自增属性。
	Description string    `xorm:"description comment('描述')"` // string类型默认映射为varchar(255)，varchar要注明长度
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
	CreatedAt   int64     `xorm:"created"` // 记住重复写created,第一个为column标签并且加单引号，不加单引号为tag，添加数据会自动更新
	UpdatedAt   time.Time `xorm:"'updated_at' updated"`
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
	Money              float64   `xorm:"money decimal"`
}

//既有LocalTime类型的，又有*LocalTime类型的，*LocalTime是考虑到有时候数据值可能为NULL，即字段值可能为空的情况。

// user-active-record用户活跃记录表
type UserActiveRecord struct {
	Id            int64  `xorm:"'id' notnull pk autoincr" `
	Uid           int64  `xorm:"'uid' int(10) notnull index comment('用户Id')"`
	Platform      int    `xorm:"'platform' SMALLINT(3) notnull comment('安卓1 ios2') "`
	ActiveDate    int    `xorm:"'active_date' int(8)  notnull  comment('日期') "`
	DeviceId      string `xorm:"'device_id' varchar(100) notnull comment('数美Id') "`
	ClientIP      string `xorm:"'client_ip'  varchar(100) notnull comment('客户端请求的ip') "`
	Imei          string `xorm:"'imei'  varchar(64) notnull comment('客户端的imei') "`
	TrustId       string `xorm:"'trusted_id' varchar(100) notnull comment('数盟的可信id') "`
	Brand         string `xorm:"'brand' varchar(100) notnull comment('手机品牌') "`
	Model         string `xorm:"'model' varchar(100) notnull comment('手机型号') "`
	SystemVersion string `xorm:"'sys_ver' varchar(64) notnull comment('系统版本号') "`
	AppVersion    int    `xorm:"'app_ver' int(8)  comment('客户端的app版本') "`
	ExtendField   string `xorm:"'extend_field' json notnull comment('json扩展字符串') "`
	CreatedAt     int64  `xorm:"'created_at'  int(10) DEFAULT(0) notnull created"`
	UpdatedAt     int64  `xorm:"'updated_at'  int(10) DEFAULT(0) notnull updated"`
}

// UserActive 用户活跃天数表 user_active
type UserActive struct {
	Id         int64 `xorm:"'id' notnull pk autoincr" mapstructure:"id"`
	Uid        int64 `xorm:"'uid' bigint notnull DEFAULT(0) unique(uid) comment('用户Id')"` // UNIQUE KEY `UQE_user_active_uid` (`uid`)  UQE_表名_括号内文字(uid)
	TotalDays  int64 `xorm:"'total_days' int(5) notnull DEFAULT(0) comment('总天数')"`
	LatestDate int64 `xorm:"'latest_date' int(6) notnull DEFAULT(0) comment('上次更新的活跃日期')"`
	CreatedAt  int64 `xorm:"'created_at'  int(10) DEFAULT(0) notnull created"`
	UpdatedAt  int64 `xorm:"'updated_at'  int(10) DEFAULT(0) notnull updated"`
}
