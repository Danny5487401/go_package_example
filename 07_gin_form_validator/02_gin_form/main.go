package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginForm struct {
	User     string `json:"user" binding:"required,min=3,max=10"`
	Password string `json:"password" binding:"required"`
}

type SignUpForm struct {
	Age        uint8  `json:"age" binding:"gte=1,lte=130"`
	Name       string `json:"name" binding:"required,min=3"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"` //跨字段
}

func main() {
	router := gin.Default()
	router.POST("/signup", func(context *gin.Context) {
		var signForm SignUpForm
		if err := context.ShouldBind(&signForm); err != nil {
			fmt.Println(err.Error())
			context.JSON(http.StatusBadRequest, gin.H{
				"err": err.Error(),
			})
			return
		}
		context.JSON(http.StatusOK, gin.H{
			"msg": "注册成功",
		})
	})

	router.POST("/loginJSON", func(c *gin.Context) {

		var loginForm LoginForm
		if err := c.ShouldBind(&loginForm); err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"msg": "登录成功",
		})
	})
	_ = router.Run(":8090")
}
