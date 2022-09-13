package util

import (
	"errors"
	"fmt"
	"net"
)

func FindFirstNonLoopbackHostInfo() {

}

func FindFirstNonLoopbackAddress() {

}

func FindFirstNonLoopbackIP() (ipv4 string, err error) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			adders, _ := netInterfaces[i].Addrs()

			for _, address := range adders {
				if inet, ok := address.(*net.IPNet); ok && !inet.IP.IsLoopback() {
					fmt.Println(inet)
					if inet.IP.To4() != nil {
						return inet.IP.String(), nil
					}
				}
			}
		}
	}

	return "", errors.New("no find")
}
