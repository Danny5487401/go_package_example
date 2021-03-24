package main

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"

)
var (
	trans ut.Translator
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
//loginForm.user ->user
func removeTopStruct(fields map[string]string)(map[string]string){
	rsp := make(map[string]string)
	for field,errMsg := range fields{
		rsp[field[strings.Index(field,".")+1:]] = errMsg
	}
	return rsp
}


func InitTrans (locale string)(err error){
	// 修改Gin 的validator引擎，实现定制

	// 类型转换成*validator.Validate
	if v,ok := binding.Validator.Engine().(*validator.Validate);ok{
		// 注册一个获取json的tag自定义方法 loginForm.User ->loginForm.user  转小写
		v.RegisterTagNameFunc(func (field reflect.StructField)string{
			name := strings.SplitN(field.Tag.Get("json"),",",2)[0]
			if name == "_"{
				return ""
			}
			return name
		})

		zhT := zh.New()  // 中文翻译器
		enT := en.New()	 // 英文翻译器
		uni := ut.New(enT,zhT,enT) //第一个为备用的语言环境，后面的参数是应该支持的语言翻译器
		trans, ok = uni.GetTranslator(locale)
		if !ok{
			return fmt.Errorf("uni.GetTranslator failed:%s",locale)
		}
		switch locale {
		case "en":
			en_translations.RegisterDefaultTranslations(v, trans)
		case "zh":
			zh_translations.RegisterDefaultTranslations(v,trans)
		default:
			en_translations.RegisterDefaultTranslations(v, trans)

		}
		return
	}
	return
}



func main()  {
	// 初始化翻译器
	if err := InitTrans("zh");err != nil{
		fmt.Println("初始化翻译器错误")
		return
	}
	router := gin.Default()
	router.POST("/signup", func(context *gin.Context) {
		var signForm SignUpForm
		if err:= context.ShouldBind(&signForm);err !=nil{
			fmt.Println(err.Error())
			context.JSON(http.StatusBadRequest,gin.H{
				"err":err.Error(),
			})
			return
		}
		context.JSON(http.StatusOK,gin.H{
			"msg":"注册成功",
		})
	})

	router.POST("/loginJSON", func(c *gin.Context) {

		var loginForm LoginForm
		if err := c.ShouldBind(&loginForm); err != nil {
			errs, ok := err.(validator.ValidationErrors) //错误转换
			if !ok{
				c.JSON(http.StatusOK,gin.H{
					"msg":err.Error(),
				})
				return
			}
			c.JSON(http.StatusBadRequest,gin.H{
				"err":removeTopStruct(errs.Translate(trans)),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"msg": "登录成功",
		})
	})
	_ = router.Run(":8090")
}

