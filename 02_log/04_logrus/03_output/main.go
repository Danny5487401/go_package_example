package main

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
)

func main() {
	// 同时将日志写到bytes.Buffer、标准输出和文件中：
	writer1 := &bytes.Buffer{}
	writer2 := os.Stdout
	writer3, err := os.OpenFile("02_log/04_logrus/03_output/log.txt", os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Fatalf("create file log.txt failed: %v", err)
	}

	log.SetOutput(io.MultiWriter(writer1, writer2, writer3))
	log.Info("info msg")

	log.SetLevel(log.TraceLevel)
	log.SetFormatter(&log.JSONFormatter{})
	log.Trace("trace msg")

}
