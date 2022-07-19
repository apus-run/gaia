package transport

import (
	"context"
	"net/url"
)

// Server ...
type Server interface {
	Start(context.Context) error
	Stop() error
	GracefullyStop(context.Context) error
}

// Endpointer is registry endpoint.
type Endpointer interface {
	Endpoint() (*url.URL, error)
}
