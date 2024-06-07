package config

import "time"

type Registry struct {
	Driver string        `mapstructure:"driver"`
	Host   string        `mapstructure:"host"`
	Port   uint32        `mapstructure:"port"`
	Schema string        `mapstructure:"schema"`
	Name   string        `mapstructure:"name"`
	Tags   []string      `mapstructure:"tags"`
	Check  RegistryCheck `mapstructure:"check"`
}

type RegistryCheck struct {
	Protocol        string        `mapstructure:"protocol"`
	Endpoint        string        `mapstructure:"endpoint"`
	Interval        time.Duration `mapstructure:"interval"`
	Timeout         time.Duration `mapstructure:"timeout"`
	DeregisterAfter time.Duration `mapstructure:"deregister_after"`
}
