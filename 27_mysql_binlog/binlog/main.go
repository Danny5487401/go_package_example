package main

import (
	"context"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
	"os"
	"time"
)

func main() {

	// Create a binlog syncer with a unique server id, the server id must be different from other MySQL's.
	// flavor is mysql or mariadb
	// 创建配置
	cfg := replication.BinlogSyncerConfig{
		ServerID: 1000000,
		Flavor:   "mysql",
		Host:     "106.14.35.115",
		Port:     3307,
		User:     "root",
		Password: "chuanzhi",
	}
	syncer := replication.NewBinlogSyncer(cfg)

	// 指定的文件及偏移位置
	streamer, _ := syncer.StartSync(mysql.Position{Name: "binlog.000002", Pos: 4})

	// or you can start a gtid replication like
	// streamer, _ := syncer.StartSyncGTID(gtidSet)
	// the mysql GTID set likes this "de278ad0-2106-11e4-9f8e-6edd0ca20947:1-2"
	// the mariadb GTID set likes this "0-1-100"

	// or we can use a timeout context
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		ev, err := streamer.GetEvent(ctx)
		cancel()

		if err == context.DeadlineExceeded {
			// meet timeout
			continue
		}

		ev.Dump(os.Stdout)
	}

}
