package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/websocket"
)

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	u := url.URL{Scheme: "ws", Host: "127.0.0.1:8090", Path: "/ws"}
	fmt.Printf("Connecting to %s\n", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println("dial:", err)
		return
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				fmt.Println("read:", err)
				return
			}
			fmt.Printf("Received: %s\n", message)
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter message: ")
		scanner.Scan()
		text := scanner.Text()

		err := c.WriteMessage(websocket.TextMessage, []byte(text))
		if err != nil {
			fmt.Println("write:", err)
			return
		}

		select {
		case <-done:
			fmt.Println("done")
			return
		case <-interrupt:
			fmt.Println("interrupt")
			return
		}
	}
	fmt.Println("bye")
}
