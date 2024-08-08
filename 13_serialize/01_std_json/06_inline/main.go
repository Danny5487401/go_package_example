package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type PageQuery struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type UserQuery struct {
	PageQuery `json:",inline"`
	Name      string `json:"name"`
}

func main() {
	jsonStr := `{"page":1,"limit":10,"name":"John"}`

	var userQuery UserQuery
	if err := json.Unmarshal([]byte(jsonStr), &userQuery); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("UserQuery: %+v\n", userQuery)
}
