package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	buf := &bytes.Buffer{}
	buf.WriteString("name = awesome web")
	buf.WriteByte('\n')
	buf.WriteString("version = 0.0.1")

	env, err := godotenv.Parse(buf)
	if err != nil {
		log.Fatal(err)
	}

	err = godotenv.Write(env, "44_godotenv/test.env")
	if err != nil {
		log.Fatal(err)
	}

	err = godotenv.Load("44_godotenv/test.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	name := os.Getenv("name")
	version := os.Getenv("version")
	fmt.Println(name, version)

	// now do something with s3 or whatever
}
