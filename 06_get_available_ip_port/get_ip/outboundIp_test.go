package outboundIp

import (
	"fmt"
	"testing"
)

func TestOutBoundIp(t *testing.T) {
	outIp, outPort, err := GetOutboundIP()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("对外Ip是:%v, 对外Port是:%v", outIp, outPort)
}
