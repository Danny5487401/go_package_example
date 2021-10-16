package mongo_conn_struct_for_pool

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"context"
	"fmt"
	"log"
	"time"
)

type MongoConn struct {
	clientOptions *options.ClientOptions
	client        *mongo.Client
	collections   *mongo.Collection
}

var mongoConn *MongoConn

func InitMongoConn(url, user, password, dbname string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Set client options
	// construct url: mongodb://username:password@127.0.0.1:27017/dbname
	mongoUrl := "mongodb://" + user + ":" + password + "@" + url + "/" + dbname
	mongoConn.clientOptions = options.Client().ApplyURI(mongoUrl)

	// Connect to MongoDB
	var err error
	mongoConn.client, err = mongo.Connect(ctx, mongoConn.clientOptions)
	if err != nil {
		log.Fatalf("connect to mongodb error: %v", err)
	}

	// Check the connection
	err = mongoConn.client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("check the connection to mongo error: %v", err)
	}

	mongoConn.collections = mongoConn.client.Database(dbname).Collection("tests")
	return nil
}

func CloseMongoConn() {
	err := mongoConn.client.Disconnect(context.TODO())
	if err != nil {
		log.Fatalf("disconnect mongo connect is error: %v", err)
		return
	}
	fmt.Printf("connection to MongoDB closed.")
}
