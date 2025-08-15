package main

// 替代 import "log"
import log "github.com/sirupsen/logrus"

func main() {
	// 在输出日志中添加文件名和方法信息
	log.SetReportCaller(true)
	log.Print("Print")
	log.Printf("Printf: %s", "print")
	log.Println("Println")

	log.Fatal("Fatal")
	log.Fatalf("Fatalf: %s", "fatal")
	log.Fatalln("Fatalln")

	log.Panic("Panic")
	log.Panicf("Panicf: %s", "panic")
	log.Panicln("Panicln")
}

/*
time="2025-08-08T11:18:45+08:00" level=info msg=Print func=main.main file="/Users/python/Downloads/git_download/go_package_example/02_log/04_logrus/01_replace_log/main.go:9"
time="2025-08-08T11:18:45+08:00" level=info msg="Printf: print" func=main.main file="/Users/python/Downloads/git_download/go_package_example/02_log/04_logrus/01_replace_log/main.go:10"
time="2025-08-08T11:18:45+08:00" level=info msg=Println func=main.main file="/Users/python/Downloads/git_download/go_package_example/02_log/04_logrus/01_replace_log/main.go:11"
time="2025-08-08T11:18:45+08:00" level=fatal msg=Fatal func=main.main file="/Users/python/Downloads/git_download/go_package_example/02_log/04_logrus/01_replace_log/main.go:13"
*/
