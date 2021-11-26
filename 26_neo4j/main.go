package main

import (
	"github.com/gin-gonic/gin"
	"go_grpc_example/26_neo4j/apis"
)

func main() {
	engine := gin.Default()
	router := engine.RouterGroup
	router.GET("hello", apis.GetAomObj)
	engine.Run(":8080")
}
