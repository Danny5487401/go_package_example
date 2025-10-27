package main

import (
	"fmt"
	"net/http"
	"time"
)

func sseHandler(w http.ResponseWriter, r *http.Request) {
	// 设置SSE所需的HTTP头
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// 检查是否支持流式响应
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	// 模拟实时数据推送
	for {
		fmt.Fprintf(w, "data: %s\n\n", time.Now().Format(time.RFC3339))
		flusher.Flush() // 立即将数据发送给客户端
		time.Sleep(2 * time.Second)
	}
}

func main() {
	http.HandleFunc("/events", sseHandler)
	http.ListenAndServe(":8090", nil)
}
