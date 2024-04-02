package app

import "net/url"

type Option func(*options)

type options struct {
	id        string
	name      string
	endpoints []url.URL
}

func WithID(id string) Option {
	return func(o *options) {
		o.id = id
	}
}

func WithName(name string) Option {
	return func(o *options) {
		o.name = name
	}
}
func WithEndpoints(endpoints []url.URL) Option {
	return func(o *options) {
		o.endpoints = endpoints
	}
}
