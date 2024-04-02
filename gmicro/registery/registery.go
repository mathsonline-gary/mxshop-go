package registery

import "net/url"

type ServiceInstance struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Version   string            `json:"version"`
	Endpoints []url.URL         `json:"endpoints"`
	MetaData  map[string]string `json:"meta_data"`
}
