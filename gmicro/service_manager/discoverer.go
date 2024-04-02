package service_manager

import "context"

type ServiceDiscoverer interface {
	GetServices(ctx context.Context, serviceName string) ([]*ServiceInstance, error)
	Watch(ctx context.Context, serviceName string) (ServiceWatcher, error)
}

type ServiceWatcher interface {
	Next() ([]*ServiceInstance, error)
	Stop() error
}
