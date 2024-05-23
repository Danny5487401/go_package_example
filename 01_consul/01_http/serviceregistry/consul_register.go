package serviceregistry

import (
	"errors"
	"fmt"
	"github.com/Danny5487401/go_package_example/01_consul/01_http/instance"
	"github.com/hashicorp/consul/api"

	"strconv"
	"unsafe"
)

type consulServiceRegistry struct {
	serviceInstances     map[string]map[string]instance.ServiceInstance
	client               api.Client
	localServiceInstance instance.ServiceInstance
}

func (c consulServiceRegistry) GetInstances(serviceId string) ([]instance.ServiceInstance, error) {
	catalogService, _, _ := c.client.Catalog().Service(serviceId, "", nil)
	if len(catalogService) > 0 {
		result := make([]instance.ServiceInstance, len(catalogService))
		for index, sever := range catalogService {
			s := instance.DefaultServiceInstance{
				InstanceId: sever.ServiceID,
				ServiceId:  sever.ServiceName,
				Host:       sever.Address,
				Port:       sever.ServicePort,
				Metadata:   sever.ServiceMeta,
			}
			result[index] = s
		}
		return result, nil
	}
	return nil, nil
}

func (c consulServiceRegistry) GetServices() ([]string, error) {
	services, _, _ := c.client.Catalog().Services(nil)
	result := make([]string, unsafe.Sizeof(services))
	index := 0
	for serviceName, _ := range services {
		result[index] = serviceName
		index++
	}
	return result, nil
}

func (c consulServiceRegistry) Register(serviceInstance instance.ServiceInstance) bool {
	// 创建注册到consul的服务到
	registration := new(api.AgentServiceRegistration)
	registration.ID = serviceInstance.GetInstanceId()
	registration.Name = serviceInstance.GetServiceId()
	registration.Port = serviceInstance.GetPort()
	var tags []string
	if serviceInstance.IsSecure() {
		tags = append(tags, "secure=true")
	} else {
		tags = append(tags, "secure=false")
	}
	//if serviceInstance.GetMetadata() != nil {
	//	for key, value := range serviceInstance.GetMetadata() {
	//		tags = append(tags, key+"="+value)
	//	}
	//}
	registration.Tags = tags
	registration.Meta = serviceInstance.GetMetadata()

	registration.Address = serviceInstance.GetHost()

	// 增加consul健康检查回调函数
	check := new(api.AgentServiceCheck)

	//schema := "http"
	//if serviceInstance.IsSecure() {
	//	schema = "https"
	//}
	check.HTTP = fmt.Sprintf("https://%s/actuator/health", registration.Address) // 我使用的内网穿透，不配置port，
	check.Timeout = "2s"
	check.Interval = "5s"
	check.DeregisterCriticalServiceAfter = "30s" // 故障检查失败30s后 consul自动将注册服务删除
	check.TLSSkipVerify = false
	check.Name = serviceInstance.GetServiceId()
	check.CheckID = serviceInstance.GetInstanceId()
	registration.Check = check

	// 注册服务到consul
	err := c.client.Agent().ServiceRegister(registration)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	if c.serviceInstances == nil {
		c.serviceInstances = map[string]map[string]instance.ServiceInstance{}
	}

	services := c.serviceInstances[serviceInstance.GetServiceId()]

	if services == nil {
		services = map[string]instance.ServiceInstance{}
	}

	services[serviceInstance.GetInstanceId()] = serviceInstance

	c.serviceInstances[serviceInstance.GetServiceId()] = services

	c.localServiceInstance = serviceInstance

	return true
}

// deregister a service
func (c consulServiceRegistry) Deregister() {
	if c.serviceInstances == nil {
		return
	}

	services := c.serviceInstances[c.localServiceInstance.GetServiceId()]

	if services == nil {
		return
	}

	delete(services, c.localServiceInstance.GetInstanceId())

	if len(services) == 0 {
		delete(c.serviceInstances, c.localServiceInstance.GetServiceId())
	}

	_ = c.client.Agent().ServiceDeregister(c.localServiceInstance.GetInstanceId())

	c.localServiceInstance = nil
}

// new a consulServiceRegistry instance
// token is optional
func NewConsulServiceRegistry(host string, port int, token string) (*consulServiceRegistry, error) {
	if len(host) < 3 {
		return nil, errors.New("check host")
	}

	if port <= 0 || port > 65535 {
		return nil, errors.New("check port, port should between 1 and 65535")
	}

	config := api.DefaultConfig()
	config.Address = host + ":" + strconv.Itoa(port)
	config.Token = token
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &consulServiceRegistry{client: *client}, nil
}
