package main

import (
	"go.mongodb.org/mongo-driver/bson"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"context"
	"fmt"
	"log"

	"go_grpc_example/05_orm/03_mongo/model"
	"go_grpc_example/05_orm/03_mongo/util"
)

func main() {
	var (
		client     = util.GetMgoCli()
		db         *mongo.Database
		collection *mongo.Collection
		cursor     *mongo.Cursor
		err        error
	)
	// 选择数据库
	db = client.Database("test2")
	// 选择表
	collection = db.Collection("collection1")
	//// 构建数据
	//lr := &model.LogRecord{
	//	JobName: "Python",
	//	Command: "echo 4",
	//	Err:     "",
	//	Content: "4",
	//	Tp: model.TimePrint{
	//		StartTime: time.Now().Unix(),
	//		EndTime:   time.Now().Unix() + 10,
	//	},
	//}
	////插入某一条数据
	//iResult, err := collection.InsertOne(context.Background(), lr)
	//if err != nil {
	//	fmt.Println("错误", err)
	//	return
	//}
	////_id:默认生成一个全局唯一ID
	//id := iResult.InsertedID.(primitive.ObjectID)
	//fmt.Println("自增ID", id.Hex())

	// 查询数据
	//如果直接使用 LogRecord{JobName: "job10"}是查不到数据的，因为其他字段有初始值0或者“”
	cond := model.FindByJobName{JobName: "Python"}

	//按照jobName字段进行过滤jobName="job10",翻页参数0-2
	if cursor, err = collection.Find(context.TODO(), cond, options.Find().SetSkip(0), options.Find().SetLimit(2)); err != nil {
		fmt.Println(err)
		return
	}
	//延迟关闭游标
	defer func() {
		if err = cursor.Close(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

	//遍历方式一：游标获取结果数据
	//for cursor.Next(context.Background()) {
	//	var lr model.LogRecord
	//	//反序列化Bson到对象
	//	if cursor.Decode(&lr) != nil {
	//		fmt.Print(err)
	//		return
	//	}
	//	//打印结果数据
	//	fmt.Println(lr)
	//}
	// 遍历方式二
	var results []model.LogRecord
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}

	for _, result := range results {
		fmt.Println(result)
	}
	/*
		背景
			使用文档前面的方法进行查询显然是很麻烦的，我们不可能每次查询都定义一个新的struct，
			是否有一种通用的struct来帮助我们作为过滤条件呢，这时候就需要使用到BSON包
		Go驱动程序中有两大类型表示BSON数据：D和Raw
			D：一个BSON文档。这种类型应该在顺序重要的情况下使用，比如MongoDB命令。
			M：一张无序的map。它和D是一样的，只是它不保持顺序。
			A：一个BSON数组。
			E：D里面的一个元素
	*/

	//按照jobName分组,countJob中存储每组的数目
	groupStage := mongo.Pipeline{bson.D{
		{"$group", bson.D{
			{"_id", "$jobName"},
			{"countJob", bson.D{
				{"$sum", 1},
			}},
		}},
	}}
	if cursor, err = collection.Aggregate(context.TODO(), groupStage); err != nil {
		log.Fatal(err)
	}
	//延迟关闭游标
	defer func() {
		if err = cursor.Close(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()
	//遍历游标
	var results2 []bson.M
	if err = cursor.All(context.TODO(), &results2); err != nil {
		log.Fatal(err)
	}
	for _, result := range results2 {
		fmt.Println(result)
	}

	// 更新数据
	var uResult *mongo.UpdateResult
	filter := bson.M{"jobName": "job10"}
	update := bson.M{"$set": bson.M{"command": "ByBsonM"}}
	//update := bson.M{"$set": model.UpdateByJobName{Command: "byModel", Content: "model"}}
	//update := bson.M{"$set": model.LogRecord{JobName:"job10",Command:"byModel"}}
	if uResult, err = collection.UpdateMany(context.TODO(), filter, update); err != nil {
		log.Fatal(err)
	}
	//uResult.MatchedCount表示符合过滤条件的记录数，即更新了多少条数据。
	log.Println(uResult.MatchedCount)

	// 删除数据
	var uResult2 *mongo.DeleteResult
	filter2 := bson.M{"content": "100"}
	//3.删除开始时间早于当前时间的数据
	if uResult2, err = collection.DeleteMany(context.TODO(), filter2); err != nil {
		log.Fatal(err)
	}
	log.Println(uResult2.DeletedCount)

	// 删除小于这时间
	var delCond *model.DeleteCond
	var uResult3 *mongo.DeleteResult
	delCond = &model.DeleteCond{BeforeCond: model.TimeBeforeCond{BeforeTime: time.Now().Unix()}}
	if uResult3, err = collection.DeleteMany(context.TODO(), delCond); err != nil {
		log.Fatal(err)
	}
	log.Println(uResult3.DeletedCount)
}
