package main

import (
	"fmt"
	"net"
	"strings"
)
// 获取对外ip
func GetOutboundIP()(ip string,err error){
	conn,err :=net.Dial("udp","8.8.8.8:80")
	if err !=nil{
		return
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	fmt.Println(localAddr.String())
	fmt.Println(localAddr.IP)
	fmt.Println(localAddr.Port)
	ip = strings.Split(localAddr.IP.String(),":")[0]
	return
}

func main(){
	ip,err :=GetOutboundIP()
	if err != nil{
		fmt.Println(err.Error())
	}
	fmt.Println(ip)
}
