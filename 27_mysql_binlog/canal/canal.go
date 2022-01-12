package main

import (
	"fmt"
	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/schema"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/json-iterator/go"
	"log"
	"os"
	"os/signal"
	"reflect"
	"runtime/debug"
	"strings"
	"syscall"
	"time"
)

func main() {
	go binLogListener()
	// placeholder for your handsome code
	//合建chan
	c := make(chan os.Signal)
	//监听所有信号
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGUSR1, syscall.SIGUSR2)
	//阻塞直到有信号传入
	s := <-c
	fmt.Println("退出信号", s)

}

type MasterSlaveTable struct {
	Id          int64     `xorm:"id notnull pk autoincr" `   // 如果field名称为Id而且类型为int64并且没有定义tag，则会被xorm视为主键，并且拥有自增属性。
	Description string    `xorm:"description comment('描述')"` // string类型默认映射为varchar(255)
	Name        string    `xorm:"'usr_name' notnull varchar(25) comment('用户名')" `
	CreatedAt   int64     `xorm:"created"` // 记住重复写created,第一个为column标签并且加单引号，不加单引号为tag，添加数据会自动更新
	UpdatedAt   time.Time `xorm:"'updated_at' updated"`
}

// 表名，大小写不敏感
func (MasterSlaveTable) TableName() string {
	return "master_slave_table"
}

// 数据库名称，大小写不敏感
func (MasterSlaveTable) SchemaName() string {
	return "masterSlaveDB"
}
func binLogListener() {
	c, err := getDefaultCanal()
	if err != nil {
		log.Println("启动失败", err.Error())
		return
	}

	// 获取主机master位置 SHOW MASTER STATUS
	mysqlPos, err := c.GetMasterPos()
	if err == nil {
		log.Println("当前位置", mysqlPos)
		// 设置处理函数,需要在启动canal前注册
		c.SetEventHandler(&binlogHandler{})
		c.RunFrom(mysqlPos)
	}

}
func getDefaultCanal() (*canal.Canal, error) {
	connStr := "root:chuanzhi@tcp(106.14.35.115:3307)/masterSlaveDB?charset=utf8mb4"
	dsn, err := mysqlDriver.ParseDSN(connStr)
	if err != nil {
		return nil, err
	}

	cfg := canal.NewDefaultConfig()
	cfg.HeartbeatPeriod = 20 * time.Millisecond
	cfg.ReadTimeout = time.Second * 30
	cfg.Addr = dsn.Addr
	cfg.User = dsn.User
	cfg.Password = dsn.Passwd
	cfg.Flavor = "mysql"

	// 指定database
	cfg.Dump.TableDB = dsn.DBName
	// 指定表
	cfg.IncludeTableRegex = []string{"^user"}
	cfg.Charset = "utf8mb4"

	// FLUSH TABLES WITH READ LOCK简称(FTWRL)，该命令主要用于备份工具获取一致性备份(数据与binlog位点匹配)。
	// 由于FTWRL总共需要持有两把全局的MDL锁，并且还需要关闭所有表对象，因此这个命令的杀伤性很大，执行命令时容易导致库hang住
	cfg.Dump.SkipMasterData = true

	return canal.NewCanal(cfg)
}

// 自定义二进制处理的handler
type binlogHandler struct {
	canal.DummyEventHandler // Dummy handler from external lib
	BinlogParser            // Our custom helper
}

func (h *binlogHandler) OnRow(e *canal.RowsEvent) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Print(r, " ", string(debug.Stack()))
		}
	}()

	// base value for canal.DeleteAction or canal.InsertAction
	// 删除和插入时的基本值
	var n = 0
	var k = 1

	// 更新时候的基本值
	if e.Action == canal.UpdateAction {
		n = 1
		k = 2
	}

	for i := n; i < len(e.Rows); i += k {
		// 所有监听的数据库+表名
		key := strings.ToLower(e.Table.Schema + "." + e.Table.Name)

		// 需要业务监听的数据库+表名
		key2 := strings.ToLower(MasterSlaveTable{}.SchemaName() + "." + MasterSlaveTable{}.TableName())
		switch key {
		case key2:
			masterSlaveTableInfo := MasterSlaveTable{}
			// 写入到结构体数据
			h.GetBinLogData(&masterSlaveTableInfo, e, i)
			switch e.Action {
			case canal.UpdateAction:
				// 获取旧数据
				oldMasterSlaveTableInfo := MasterSlaveTable{}
				h.GetBinLogData(&oldMasterSlaveTableInfo, e, i-1)
				fmt.Printf("masterSlaveTableInfo %d name changed from %s to %s\n", masterSlaveTableInfo.Id, oldMasterSlaveTableInfo.Name, masterSlaveTableInfo.Name)
			case canal.InsertAction:
				fmt.Printf("masterSlaveTableInfo %d is created with name %s\n", masterSlaveTableInfo.Id, masterSlaveTableInfo.Name)
			case canal.DeleteAction:
				fmt.Printf("masterSlaveTableInfo %d is deleted with name %s\n", masterSlaveTableInfo.Id, masterSlaveTableInfo.Name)
			default:
				fmt.Printf("Unknown action")
			}
		}

	}
	return nil
}

func (h *binlogHandler) String() string {
	return "binlogHandler"
}

type BinlogParser struct{}

func (m *BinlogParser) GetBinLogData(element interface{}, e *canal.RowsEvent, n int) error {
	var columnName string
	var ok bool
	v := reflect.ValueOf(element)
	s := reflect.Indirect(v)
	t := s.Type()
	num := t.NumField()
	for k := 0; k < num; k++ {
		parsedTag := parseTagSetting(t.Field(k).Tag)
		name := s.Field(k).Type().Name()

		if columnName, ok = parsedTag["COLUMN"]; !ok || columnName == "COLUMN" {
			continue
		}

		switch name {
		case "bool":
			s.Field(k).SetBool(m.boolHelper(e, n, columnName))
		case "int":
			s.Field(k).SetInt(m.intHelper(e, n, columnName))
		case "string":
			s.Field(k).SetString(m.stringHelper(e, n, columnName))
		case "Time":
			timeVal := m.dateTimeHelper(e, n, columnName)
			s.Field(k).Set(reflect.ValueOf(timeVal))
		case "float64":
			s.Field(k).SetFloat(m.floatHelper(e, n, columnName))
		default:
			if _, ok := parsedTag["FROMJSON"]; ok {

				newObject := reflect.New(s.Field(k).Type()).Interface()
				json := m.stringHelper(e, n, columnName)

				jsoniter.Unmarshal([]byte(json), &newObject)

				s.Field(k).Set(reflect.ValueOf(newObject).Elem().Convert(s.Field(k).Type()))
			}
		}
	}
	return nil
}

// 解析时间
func (m *BinlogParser) dateTimeHelper(e *canal.RowsEvent, n int, columnName string) time.Time {

	columnId := m.getBinlogIdByName(e, columnName)
	if e.Table.Columns[columnId].Type != schema.TYPE_TIMESTAMP {
		panic("Not dateTime type")
	}
	t, _ := time.Parse("2006-01-02 15:04:05", e.Rows[n][columnId].(string))

	return t
}

func (m *BinlogParser) intHelper(e *canal.RowsEvent, n int, columnName string) int64 {

	columnId := m.getBinlogIdByName(e, columnName)
	if e.Table.Columns[columnId].Type != schema.TYPE_NUMBER {
		return 0
	}

	switch e.Rows[n][columnId].(type) {
	case int8:
		return int64(e.Rows[n][columnId].(int8))
	case int32:
		return int64(e.Rows[n][columnId].(int32))
	case int64:
		return e.Rows[n][columnId].(int64)
	case int:
		return int64(e.Rows[n][columnId].(int))
	case uint8:
		return int64(e.Rows[n][columnId].(uint8))
	case uint16:
		return int64(e.Rows[n][columnId].(uint16))
	case uint32:
		return int64(e.Rows[n][columnId].(uint32))
	case uint64:
		return int64(e.Rows[n][columnId].(uint64))
	case uint:
		return int64(e.Rows[n][columnId].(uint))
	}
	return 0
}

func (m *BinlogParser) floatHelper(e *canal.RowsEvent, n int, columnName string) float64 {

	columnId := m.getBinlogIdByName(e, columnName)
	if e.Table.Columns[columnId].Type != schema.TYPE_FLOAT {
		panic("Not float type")
	}

	switch e.Rows[n][columnId].(type) {
	case float32:
		return float64(e.Rows[n][columnId].(float32))
	case float64:
		return float64(e.Rows[n][columnId].(float64))
	}
	return float64(0)
}

func (m *BinlogParser) boolHelper(e *canal.RowsEvent, n int, columnName string) bool {

	val := m.intHelper(e, n, columnName)
	if val == 1 {
		return true
	}
	return false
}

func (m *BinlogParser) stringHelper(e *canal.RowsEvent, n int, columnName string) string {

	columnId := m.getBinlogIdByName(e, columnName)
	if e.Table.Columns[columnId].Type == schema.TYPE_ENUM {

		values := e.Table.Columns[columnId].EnumValues
		if len(values) == 0 {
			return ""
		}
		if e.Rows[n][columnId] == nil {
			//Если в енум лежит нуул ставим пустую строку
			return ""
		}

		return values[e.Rows[n][columnId].(int64)-1]
	}

	value := e.Rows[n][columnId]

	switch value := value.(type) {
	case []byte:
		return string(value)
	case string:
		return value
	}
	return ""
}

func (m *BinlogParser) getBinlogIdByName(e *canal.RowsEvent, name string) int {
	for id, value := range e.Table.Columns {
		if value.Name == name {
			return id
		}
	}
	panic(fmt.Sprintf("There is no column %s in table %s.%s", name, e.Table.Schema, e.Table.Name))
}

func parseTagSetting(tags reflect.StructTag) map[string]string {
	settings := map[string]string{}
	for _, str := range []string{tags.Get("sql"), tags.Get("gorm")} {
		tags := strings.Split(str, ";")
		for _, value := range tags {
			v := strings.Split(value, ":")
			k := strings.TrimSpace(strings.ToUpper(v[0]))
			if len(v) >= 2 {
				settings[k] = strings.Join(v[1:], ":")
			} else {
				settings[k] = k
			}
		}
	}
	return settings
}
