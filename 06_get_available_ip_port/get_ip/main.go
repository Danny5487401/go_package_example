package main

import (
	"fmt"
	"net"
	"strings"
)

// 获取对外ip
func GetOutboundIP() (ip string, err error) {
	// 8.8.8.8:google, 114.114.114.114是国内移动、电信和联通通用的DNS
	//conn, err := net.Dial("udp", "8.8.8.8:80")
	conn, err := net.Dial("udp", "114.114.114.114:80")
	if err != nil {
		return
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	fmt.Println(localAddr.String())
	fmt.Println("对外Ip是", localAddr.IP, "对外Port是", localAddr.Port)
	ip = strings.Split(localAddr.IP.String(), ":")[0]
	return
}

func main() {
	ip, err := GetOutboundIP()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(ip)
}
