package main

import (
	"fmt"
	"log"
	"time"

	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/jmoiron/sqlx"
)

func main() {
	var (
		// 多个主机，使用逗号分割
		host1    = "tencent.danny.games:9000"
		host2    = "tencent.danny.games:9000"
		username = "root"
		password = "chuanzhi"
		database = "default"
		tcpInfo  = "tcp://%s?username=%s&password=%s&database=%s&read_timeout=5&write_timeout=5&debug=true&compress=true&alt_hosts=%s"
	)
	tcpInfo = fmt.Sprintf(tcpInfo, host1, username, password, database, host2)
	connect, err := sqlx.Open("clickhouse", tcpInfo)
	if err != nil {
		log.Fatal(err)
	}
	var items []struct {
		CountryCode string    `db:"country_code"`
		OsID        uint8     `db:"os_id"`
		BrowserID   uint8     `db:"browser_id"`
		Categories  []int16   `db:"categories"`
		ActionTime  time.Time `db:"action_time"`
	}

	if err := connect.Select(&items, "SELECT country_code, os_id, browser_id, categories, action_time FROM example"); err != nil {
		log.Fatal(err)
	}

	for _, item := range items {
		log.Printf("country: %s, os: %d, browser: %d, categories: %v, action_time: %s", item.CountryCode, item.OsID, item.BrowserID, item.Categories, item.ActionTime)
	}
}
