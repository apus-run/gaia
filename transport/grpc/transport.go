package grpc

import (
	"github.com/apus-run/gaia/transport"
)

var _ transport.Transporter = (*Transport)(nil)

// Transport is an HTTP transport.
type Transport struct {
	endpoint  string
	operation string
}

// Kind returns the transport kind.
func (tr *Transport) Kind() transport.Kind {
	return transport.KindHTTP
}

// Endpoint returns the transport endpoint.
func (tr *Transport) Endpoint() string {
	return tr.endpoint
}

// Operation returns the transport operation.
func (tr *Transport) Operation() string {
	return tr.operation
}
