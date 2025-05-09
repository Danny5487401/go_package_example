package outboundIp

import (
	"fmt"
	"net"
)

// 获取对外ip
func GetOutboundIP() (ip string, port int, err error) {
	// 8.8.8.8:google, 114.114.114.114是国内移动、电信和联通通用的DNS
	//conn, err := net.Dial("udp", "8.8.8.8:80")
	conn, err := net.Dial("udp", "114.114.114.114:80")
	if err != nil {
		return
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	fmt.Println("IP加端口完整信息：", localAddr.String())
	ip = localAddr.IP.String()
	port = localAddr.Port
	return
}
