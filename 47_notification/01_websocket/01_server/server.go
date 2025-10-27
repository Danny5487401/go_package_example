package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4196, // //指定读缓存大小
	WriteBufferSize: 1124,
	CheckOrigin:     checkOrigin,
}

// 检测请求来源
func checkOrigin(r *http.Request) bool {
	if r.Method != "GET" {
		fmt.Println("method is not GET")
		return false
	}
	if r.URL.Path != "/ws" {
		fmt.Println("path error")
		return false
	}
	return true
}

// ServerHTTP 用于升级协议
func ServerHTTP(w http.ResponseWriter, r *http.Request) {
	// 收到http请求之后升级协议
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error during connection upgrade:", err)
		return
	}
	defer conn.Close()

	for {
		// 服务端读取客户端请求
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error during message reading:", err)
			break
		}
		log.Printf("Received:%s", message)

		// 开启关闭连接监听
		conn.SetCloseHandler(func(code int, text string) error {
			fmt.Println(code, text) // 断开连接时将打印code和text
			return nil
		})

		//服务端给客户端返回请求
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			log.Println("Error during message writing:", err)
			return
		}

	}
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Index Page")
}

func main() {
	http.HandleFunc("/ws", ServerHTTP)
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe("localhost:8090", nil))
}
