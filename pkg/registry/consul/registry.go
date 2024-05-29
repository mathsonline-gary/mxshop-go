package consul

import (
	"context"
	"net/url"
	"strconv"

	"github.com/hashicorp/consul/api"
	"github.com/zycgary/mxshop-go/pkg/registry"
)

var _ registry.Registrar = (*registrar)(nil)

type Option func(registrar *registrar)

type registrar struct {
	cli    *api.Client
	checks api.AgentServiceChecks
}

func New(client *api.Client, opts ...Option) registry.Registrar {
	chs := make(api.AgentServiceChecks, 0, 1)
	r := &registrar{
		cli:    client,
		checks: chs,
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

func WithCheck(check *api.AgentServiceCheck) Option {
	return func(r *registrar) {
		r.checks = append(r.checks, check)
	}
}

func (r registrar) Register(_ context.Context, instance *registry.Instance) error {
	addresses := make(map[string]api.ServiceAddress, len(instance.Endpoints))
	for _, endpoint := range instance.Endpoints {
		raw, err := url.Parse(endpoint)
		if err != nil {
			return err
		}
		port, _ := strconv.ParseUint(raw.Port(), 10, 16)
		addresses[raw.Scheme] = api.ServiceAddress{Address: endpoint, Port: int(port)}
	}

	registration := &api.AgentServiceRegistration{
		Name:            instance.Name,
		ID:              instance.ID,
		Tags:            instance.Tags,
		Meta:            instance.Metadata,
		TaggedAddresses: addresses,
	}

	if len(r.checks) > 0 {
		registration.Checks = r.checks
	}

	if err := r.cli.Agent().ServiceRegister(registration); err != nil {
		return err
	}

	return nil
}

func (r registrar) Deregister(_ context.Context, instance *registry.Instance) error {
	return r.cli.Agent().ServiceDeregister(instance.ID)
}
