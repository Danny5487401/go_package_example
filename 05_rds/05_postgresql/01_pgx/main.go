package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"time"
)

func main() {
	// urlExample := "postgres://username:password@localhost:5432/database_name"
	url := "postgres://postgres:postgres@localhost:5432/postgres"
	// 使用并发安全的连接池
	dbpool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	var greeting string
	err = dbpool.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(greeting)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// 定期打印连接池状态
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			printPoolStats(dbpool)
		case <-ctx.Done():
			return
		}
	}
}

func printPoolStats(pool *pgxpool.Pool) {
	stat := pool.Stat()
	now := time.Now().Format("15:04:05.000")

	// 计算关键衍生指标
	_ = float64(stat.AcquireCount()) / stat.AcquireDuration().Seconds()
	idleRatio := float64(stat.IdleConns()) / float64(stat.TotalConns())

	fmt.Printf("[%s] 连接池状态:\n", now)
	fmt.Printf("  活跃连接: %d/%d (使用率: %.0f%%)\n",
		stat.AcquiredConns(), stat.MaxConns(),
		float64(stat.AcquiredConns())/float64(stat.MaxConns())*100)
	fmt.Printf("  空闲连接: %d (占比: %.0f%%)\n",
		stat.IdleConns(), idleRatio*100)
	fmt.Printf("  累计创建连接: %d (复用率: %.2f)\n",
		stat.NewConnsCount(),
		float64(stat.AcquireCount())/float64(stat.NewConnsCount()))
	fmt.Printf("  平均获取耗时: %.2fms\n",
		float64(stat.AcquireDuration().Microseconds())/float64(stat.AcquireCount())/1000)

}
