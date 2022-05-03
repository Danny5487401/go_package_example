package main

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

func main() {
	// 初始化一个 文件/目录 监听器
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close() // 最后结束程序时关闭watcher

	done := make(chan bool)
	go func() { // 启动一个协程来单独处理watcher发来的事件
		for {
			select {
			case event, ok := <-watcher.Events: // 正常的事件的处理逻辑
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
				}
			case err, ok := <-watcher.Errors: // 发生错误时的处理逻辑
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add("21_viper/02_fsnotify/foo.yaml") // 使watcher监控21_viper/02_fsnotify/foo.yaml
	if err != nil {
		log.Fatal(err)
	}
	<-done // 使主协程不退出
}
