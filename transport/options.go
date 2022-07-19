package transport

import "fmt"

type Option func(si *ServiceInfo)

// ServiceInfo represents service info
type ServiceInfo struct {
	Name    string `json:"name"`
	Scheme  string `json:"scheme"`
	Address string `json:"address"`
}

// Label ...
func (si *ServiceInfo) Label() string {
	return fmt.Sprintf("%s://%s", si.Scheme, si.Address)
}

func ApplyOptions(opts ...Option) ServiceInfo {
	si := defaultServiceInfo()
	for _, o := range opts {
		o(&si)
	}
	return si
}

// Scheme with ServiceInfo scheme.
func WithScheme(scheme string) Option {
	return func(c *ServiceInfo) {
		c.Scheme = scheme
	}
}

// Address with ServiceInfo address.
func WithAddress(address string) Option {
	return func(c *ServiceInfo) {
		c.Address = address
	}
}

// Name with ServiceInfo name.
func WithName(name string) Option {
	return func(c *ServiceInfo) {
		c.Name = name
	}
}

func defaultServiceInfo() ServiceInfo {
	// 设置默认值
	return ServiceInfo{}
}
