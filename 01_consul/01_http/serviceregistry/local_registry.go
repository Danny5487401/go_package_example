package serviceregistry

import (
	"errors"
	"github.com/Danny5487401/go_package_example/01_consul/01_http/instance"
)

type localServiceRegistry struct {
	ServiceInstances     map[string]map[string]instance.ServiceInstance
	localServiceInstance instance.ServiceInstance
}

func NewLocalServiceRegistry() *localServiceRegistry {
	return &localServiceRegistry{ServiceInstances: map[string]map[string]instance.ServiceInstance{}}
}

func (d localServiceRegistry) Description() string {
	return "localServiceRegistry"
}

func (d localServiceRegistry) GetInstances(serviceId string) ([]instance.ServiceInstance, error) {
	if ret, ok := d.ServiceInstances[serviceId]; ok {
		var result []instance.ServiceInstance
		for _, value := range ret {
			result = append(result, value)
		}
		return result, nil
	} else {
		return nil, errors.New("no data")
	}
}

func (d localServiceRegistry) GetServices() ([]string, error) {
	var result []string
	for key, _ := range d.ServiceInstances {
		result = append(result, key)
	}
	return result, nil
}

func (d localServiceRegistry) Register(serviceInstance instance.DefaultServiceInstance) bool {
	if d.ServiceInstances == nil {
		d.ServiceInstances = map[string]map[string]instance.ServiceInstance{}
	}

	services := d.ServiceInstances[serviceInstance.GetServiceId()]

	if services == nil {
		services = map[string]instance.ServiceInstance{}
	}

	services[serviceInstance.InstanceId] = serviceInstance

	d.ServiceInstances[serviceInstance.GetServiceId()] = services

	d.localServiceInstance = serviceInstance

	return true
}

func (d localServiceRegistry) Deregister() bool {
	if d.ServiceInstances == nil {
		return true
	}

	if d.localServiceInstance == nil {
		return true
	}

	services := d.ServiceInstances[d.localServiceInstance.GetServiceId()]

	if services == nil {
		return true
	}

	delete(services, d.localServiceInstance.GetInstanceId())

	if len(services) == 0 {
		delete(d.ServiceInstances, d.localServiceInstance.GetServiceId())
	}

	d.localServiceInstance = nil

	return true
}
