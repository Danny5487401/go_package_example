package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
)

func main() {

	dsn := "clickhouse://default:K5cCigzWXk@my-clickhouse.clickhouse.svc.cluster.local:9000/my_database?dial_timeout=1000ms&max_execution_time=60"
	options, err := clickhouse.ParseDSN(dsn)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	// 配置连接参数
	db, err := clickhouse.Open(options)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(ctx); err != nil {
		var exception *clickhouse.Exception
		if errors.As(err, &exception) {
			fmt.Printf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		} else {
			log.Fatal(err)
		}
		return
	}
	if err = db.Exec(ctx, "DROP TABLE IF EXISTS example"); err != nil {
		log.Fatal(err)
	}

	// 创建表: Memory 引擎以未压缩的形式将数据存储在 RAM 中
	err = db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS example (
			country FixedString(2),
			os_id        UInt8,
			browser_id   UInt8,
			categories   Array(Int16),
			action_day   Date,
			action_time  DateTime
		) engine=Memory
	`)

	if err != nil {
		log.Fatal(err)
	}

	// 插入数据
	stmt, err := db.PrepareBatch(ctx, "INSERT INTO example (country, os_id, browser_id, categories, action_day, action_time) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 100; i++ {
		if err = stmt.Append(
			"RU",
			uint8(10+i),
			uint8(100+i),
			[]int16{1, 2, 3},
			time.Now(),
			time.Now(),
		); err != nil {
			log.Fatal(err)
		}
	}

	if err = stmt.Send(); err != nil {
		log.Fatal(err)
	}

	// 查询数据
	rows, err := db.Query(ctx, "SELECT country, os_id, browser_id, categories, action_day, action_time FROM example")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var e = Example{}
		if err = rows.ScanStruct(&e); err != nil {
			log.Fatal(err)
		}
		log.Printf("country: %s, os: %d, browser: %d, categories: %v, action_day: %s, action_time: %s",
			e.Country, e.Os, e.Browser, e.Categories, e.ActionDay, e.ActionTime)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	// 删除表
	if err = db.Exec(ctx, "DROP TABLE example"); err != nil {
		log.Fatal(err)
	}
}

type Example struct {
	Country    string    `ch:"country"`
	Os         uint8     `ch:"os_id"`
	Browser    uint8     `ch:"browser_id"`
	Categories []int16   `ch:"categories"`
	ActionDay  time.Time `ch:"action_day"`
	ActionTime time.Time `ch:"action_time"`
}
