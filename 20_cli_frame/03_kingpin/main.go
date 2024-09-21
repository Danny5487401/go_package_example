package main

import (
	"os"
	"strings"

	"github.com/alecthomas/kingpin/v2"
)

var (
	app = kingpin.New("chat", "A command-line chat application.")
	//bool类型参数，可以通过 --debug使该值为true
	debug = app.Flag("debug", "Enable debug mode.").Bool()
	//识别 ./cli register
	register = app.Command("register", "Register a new user.")
	// ./cli register之后的参数，可通过./cli register danny 123456 传入name为danny pwd为123456 参数类型为字符串
	registerName = register.Arg("name", "Name for user.").Required().String()
	registerPwd  = register.Arg("pwd", "pwd of user.").Required().String()
	//识别 ./cli post
	post = app.Command("post", "Post a message to a channel.")
	//可以通过 ./cli post --image file  或者 ./cli post -i file 传入文件
	postImage = post.Flag("image", "Image to post.").Short('i').String()
	//可以通过./cli post txt 传入字符串，有默认值"hello world"
	postText = post.Arg("text", "Text to post.").Default("hello world").Strings()
)

func main() {
	//从os接收参数传给kingpin处理
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case register.FullCommand():
		println("name:" + *registerName)
		println("pwd:" + *registerPwd)
	case post.FullCommand():
		println((*postImage))
		text := strings.Join(*postText, " ")
		println("Post:", text)
	}
	if *debug == true {
		println("debug")
	}
}
