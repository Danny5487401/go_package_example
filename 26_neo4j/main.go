package main

import (
	"github.com/Danny5487401/go_package_example/26_neo4j/apis"
	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()
	router := engine.RouterGroup
	router.GET("getStudent", apis.GetAomObj)
	router.GET("creatStudent", apis.CreateObj)
	engine.Run(":8080")
}
