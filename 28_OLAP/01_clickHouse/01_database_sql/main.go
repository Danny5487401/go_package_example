package main

import (
	"context"
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
	db := clickhouse.OpenDB(options)
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(time.Hour)

	if err := db.PingContext(ctx); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		} else {
			fmt.Printf("err : %v \n", err)
		}
		return
	}
	if _, err := db.ExecContext(ctx, "DROP TABLE IF EXISTS example"); err != nil {
		log.Fatal(err)
	}
	// 创建表: Memory 引擎以未压缩的形式将数据存储在 RAM 中
	_, err = db.ExecContext(ctx, `
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
	var (
		tx, _   = db.Begin()
		stmt, _ = tx.Prepare("INSERT INTO example (country, os_id, browser_id, categories, action_day, action_time) VALUES (?, ?, ?, ?, ?, ?)")
	)
	defer stmt.Close()

	for i := 0; i < 100; i++ {
		if _, err := stmt.Exec(
			"RU",
			10+i,
			100+i,
			[]int16{1, 2, 3},
			time.Now(),
			time.Now(),
		); err != nil {
			log.Fatal(err)
		}
	}

	if err = tx.Commit(); err != nil {
		log.Fatal(err)
	}

	// 查询数据
	rows, err := db.Query("SELECT country, os_id, browser_id, categories, action_day, action_time FROM example")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			country               string
			os, browser           uint8
			categories            []int16
			actionDay, actionTime time.Time
		)
		if err = rows.Scan(&country, &os, &browser, &categories, &actionDay, &actionTime); err != nil {
			log.Fatal(err)
		}
		log.Printf("country: %s, os: %d, browser: %d, categories: %v, action_day: %s, action_time: %s", country, os, browser, categories, actionDay, actionTime)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	// 删除表
	if _, err = db.Exec("DROP TABLE example"); err != nil {
		log.Fatal(err)
	}
}
