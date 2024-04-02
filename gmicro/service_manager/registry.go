package service_manager

import (
	"context"
)

type ServiceRegistry interface {
	Register(context.Context, *ServiceInstance) error
	Deregister(context.Context, *ServiceInstance) error
}
