package registry

import "context"

type Registrar interface {
	Register(ctx context.Context, instance *Instance) error
	Deregister(ctx context.Context, instance *Instance) error
}
