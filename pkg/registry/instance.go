package registry

import "fmt"

// Instance is an instance of a service in a discovery system.
type Instance struct {
	// ID is the unique instance ID as registered.
	ID string `json:"id"`
	// Name is the service name as registered.
	Name string `json:"name"`
	// Tags are the tags of the service instance for further filtering or categorisation.
	Tags []string `json:"tags"`
	// Metadata is the kv pair metadata associated with the service instance.
	Metadata map[string]string `json:"metadata"`
	// Endpoints are endpoint addresses of the service instance.
	// schema:
	//   http://127.0.0.1:8000?isSecure=false
	//   v1://127.0.0.1:9000?isSecure=false
	Endpoints []string `json:"endpoints"`
}

// String returns the instance as a string.
func (i *Instance) String() string {
	return fmt.Sprintf("%s-%s", i.Name, i.ID)
}
