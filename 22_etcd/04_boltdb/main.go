package main

import (
	"fmt"
	"github.com/etcd-io/bbolt"
	"log"
)

var world = []byte("greeting")

func main() {
	// 1. 连接数据库，打开boltdb文件，获取db对象
	db, err := bbolt.Open("22_etcd/04_boltdb/bolt.db", 0644, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	key := []byte("hello")
	value := []byte("Hello World!")

	// 2. store some data
	err = db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(world)
		if err != nil {
			return err
		}

		err = bucket.Put(key, value)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	// 3. retrieve the data,内部封装事务
	err = db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(world)
		if bucket == nil {
			return fmt.Errorf("bucket %s not found", world)
		}

		val := bucket.Get(key)
		fmt.Println(string(val))

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

}
