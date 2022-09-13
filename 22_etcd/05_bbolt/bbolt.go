package main

import (
	"bytes"

	"log"

	bolt "go.etcd.io/bbolt"
)

func main() {

	db, err := bolt.Open("my.db", 0600, nil)

	if err != nil {

		log.Fatal(err)

	}

	defer db.Close()

	if err := db.Update(func(tx *bolt.Tx) error {

		b, err := tx.CreateBucket([]byte("b1"))

		if err != nil {

			log.Fatal(err)

		}

		if err := b.Put([]byte("foo"), []byte("bar")); err != nil {

			log.Fatal(err)

		}

		if v := b.Get([]byte("foo")); !bytes.Equal(v, []byte("bar")) {

			log.Fatalf("unexpected value: %v", v)

		}

		return nil

	}); err != nil {

		log.Fatal(err)

	}

}
