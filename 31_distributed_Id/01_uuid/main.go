package main

import (
	"fmt"
	"github.com/google/uuid"
)

func main() {

	// v4 版本
	uuidV4 := uuid.Must(uuid.NewRandom())
	fmt.Println("UUID v4:", uuidV4) // 53f497f8-0d9c-45c6-8d86-d3cb1d5a4a13

	// v7 版本
	uuidV7 := uuid.Must(uuid.NewV7())
	fmt.Println("UUID v7:", uuidV7) // 0193e8e4-99e9-76fd-a660-2561cee76c39

}
