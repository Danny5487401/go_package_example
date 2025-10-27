package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Users struct {
	Name    string   `yaml:"name"`
	Age     int8     `yaml:"age"`
	Address string   `yaml:"address"`
	Hobby   []string `yaml:"hobby"`
}

func main() {

	file, err := os.ReadFile("21_viper/04_yaml/test.yaml")
	if err != nil {
		log.Fatal(err)
	}
	var data [7]Users
	err = yaml.Unmarshal(file, &data)

	if err != nil {
		log.Fatal(err)
	}
	for _, v := range data {
		fmt.Println(v)
	}
}
