package serviceregistry

import "go_package_example/01_consul/01_http/instance"

type ServiceRegistry interface {
	Register(serviceInstance instance.ServiceInstance) bool

	Deregister()
}
