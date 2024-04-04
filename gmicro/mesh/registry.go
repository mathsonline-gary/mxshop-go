package mesh

import (
	"context"
)

type Registry interface {
	Register(context.Context, *ServiceInstance) error
	Deregister(context.Context, *ServiceInstance) error
}
