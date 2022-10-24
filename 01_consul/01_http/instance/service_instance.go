package instance

import (
	"go_package_example/01_consul/01_http/util"
	"math/rand"
	"strconv"
	"time"
)

/**
 * Represents an instance of a service in a discovery system.
 */
type ServiceInstance interface {

	// return The unique instance ID as registered.
	GetInstanceId() string

	// return The service ID as registered.
	GetServiceId() string

	// return The hostname of the registered service instance.
	GetHost() string

	// return The port of the registered service instance.
	GetPort() int

	// return Whether the port of the registered service instance uses HTTPS.
	IsSecure() bool

	// return The key / value pair metadata associated with the service instance.
	GetMetadata() map[string]string
}

/**
 *
 */
type DefaultServiceInstance struct {
	InstanceId string
	ServiceId  string
	Host       string
	Port       int
	Secure     bool
	Metadata   map[string]string
}

/**
 * default
 */
func NewDefaultServiceInstance(serviceId string, host string, port int, secure bool,
	metadata map[string]string, instanceId string) (*DefaultServiceInstance, error) {

	if len(host) == 0 {
		localIP, err := util.FindFirstNonLoopbackIP()
		if err != nil {
			return nil, err
		}
		host = localIP
	}

	if len(instanceId) == 0 {
		instanceId = serviceId + "-" + strconv.FormatInt(time.Now().Unix(), 10) + "-" + strconv.Itoa(rand.Intn(9000)+1000)
	}

	return &DefaultServiceInstance{InstanceId: instanceId, ServiceId: serviceId, Host: host, Port: port, Secure: secure, Metadata: metadata}, nil
}

func (serviceInstance DefaultServiceInstance) GetInstanceId() string {
	return serviceInstance.InstanceId
}

func (serviceInstance DefaultServiceInstance) GetServiceId() string {
	return serviceInstance.ServiceId
}

func (serviceInstance DefaultServiceInstance) GetHost() string {
	return serviceInstance.Host
}

func (serviceInstance DefaultServiceInstance) GetPort() int {
	return serviceInstance.Port
}

func (serviceInstance DefaultServiceInstance) IsSecure() bool {
	return serviceInstance.Secure
}

func (serviceInstance DefaultServiceInstance) GetMetadata() map[string]string {
	return serviceInstance.Metadata
}
