package util

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 虽然mgo 十分好用且稳定, 但是由于mgo不再维护 不支持事务, 并且golang 推荐使用官方驱动 mongo driver. 所以更换成mongo driver.
var mgoCli *mongo.Client

func GetMgoCli() *mongo.Client {
	if mgoCli == nil {
		initEngine()
	}
	return mgoCli
}

func initEngine() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// 连接:在本地的时候mgo 的mongodburl 可以写成127.0.0.1,但是mongo driver 必须写成 mongodb://127.0.0.1
	//uri := "mongodb://ali.danny.games:27017"
	user := "python"
	password := "chuanzhi"
	url := "ali.danny.games:27017"
	dbname := "db1"
	auth := "authSource=admin"

	var numPool uint64 = 10
	//开启auth认证
	uri := "mongodb://" + user + ":" + password + "@" + url + "/" + dbname + "?" + auth
	var err error

	// primary （只主）只从 primary 节点读数据，这个是默认设置
	mode, err := readpref.ModeFromString("primary")
	if err != nil {
		return
	}
	rp, err := readpref.New(mode)

	if err != nil {
		return
	}
	opt := options.Client().SetReadPreference(rp).ApplyURI(uri)

	// 设置连接池,默认100
	opt.SetMaxPoolSize(numPool)

	opts := []*options.ClientOptions{opt}

	// 添加中间件
	opts = append(opts, options.Client().SetMonitor(getMonitor()))

	// 开始连接
	mgoCli, err = mongo.Connect(ctx, opts...)
	if err != nil {
		log.Fatal(err)
	}
	// 检查连接
	if err := mgoCli.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected and pinged.")
}

//打印日志
func getMonitor() *event.CommandMonitor {
	return &event.CommandMonitor{
		Started: func(ctx context.Context, startedEvent *event.CommandStartedEvent) {
			fmt.Printf("[mongo] started:%+v \n", *startedEvent)
		},
		Succeeded: func(ctx context.Context, succeededEvent *event.CommandSucceededEvent) {
			fmt.Printf("[mongo] success:%+v \n", *succeededEvent)
		},
		Failed: func(ctx context.Context, failedEvent *event.CommandFailedEvent) {
			fmt.Printf("[mongo] ERROR failed:%+v \n", *failedEvent)
		},
	}
}
