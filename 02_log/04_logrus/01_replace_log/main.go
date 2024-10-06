package main

// 替代 import "log"
import log "github.com/sirupsen/logrus"

func main() {
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
