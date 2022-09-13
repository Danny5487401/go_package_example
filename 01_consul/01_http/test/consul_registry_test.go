package test

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go_package_example/01_consul/01_http/instance"
	"go_package_example/01_consul/01_http/serviceregistry"
	"go_package_example/01_consul/01_http/util"
	"math/rand"
	"testing"
	"time"
)

func TestConsulServiceRegistry(t *testing.T) {
	host := "tencent.danny.games"
	port := 8500
	token := ""
	registryDiscoveryClient, err := serviceregistry.NewConsulServiceRegistry(host, port, token)

	ip, err := util.FindFirstNonLoopbackIP()
	if err != nil {
		t.Error(err)
		panic(err)
	}

	fmt.Println(ip)
	rand.Seed(time.Now().UnixNano())

	hostSrv := "aa96-154-86-159-40.ngrok.io"
	si, _ := instance.NewDefaultServiceInstance("go-user-server", hostSrv, 8000,
		false, map[string]string{"user": "danny"}, "")

	registryDiscoveryClient.Register(si)

	r := gin.Default()
	r.GET("/actuator/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	err = r.Run("192.168.16.111:2222")
	if err != nil {
		registryDiscoveryClient.Deregister()
	}
}

func TestConsulServiceDiscovery(t *testing.T) {
	host := "127.0.0.1"
	port := 8500
	token := ""
	registryDiscoveryClient, err := serviceregistry.NewConsulServiceRegistry(host, port, token)
	if err != nil {
		panic(err)
	}

	t.Log(registryDiscoveryClient.GetServices())

	t.Log(registryDiscoveryClient.GetInstances("ecm-monitor"))
}
