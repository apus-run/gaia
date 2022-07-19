package gaia

import (
	"context"
	"net/url"
	"os"
	"syscall"
	"time"

	"github.com/apus-run/gaia/v1/log"
	"github.com/apus-run/gaia/v1/registry"
	"github.com/apus-run/gaia/v1/transport"
)

// Option is an application option.
type Option func(o *options)

// options is an application options.
type options struct {
	id        string
	name      string
	version   string
	metadata  map[string]string
	endpoints []*url.URL

	ctx  context.Context
	sigs []os.Signal

	registry        registry.Registry
	registryTimeout time.Duration
	stopTimeout     time.Duration
	servers         []transport.Server

	logger log.Logger
}

// defaultOptions 初始化默认值
func defaultOptions() options {
	return options{
		ctx:             context.Background(),
		sigs:            []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
		registryTimeout: 10 * time.Second,
		stopTimeout:     10 * time.Second,
	}
}

// newOptions returns a new options
func newOptions(opts ...Option) options {
	opt := defaultOptions()

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// WithID with app id
func WithID(id string) Option {
	return func(o *options) {
		o.id = id
	}
}

// WithName .
func WithName(name string) Option {
	return func(o *options) {
		o.name = name
	}
}

// WithVersion with a version
func WithVersion(version string) Option {
	return func(o *options) {
		o.version = version
	}
}

// WithContext with a context
func WithContext(ctx context.Context) Option {
	return func(o *options) {
		o.ctx = ctx
	}
}

// WithSignal with some system signal
func WithSignal(sigs ...os.Signal) Option {
	return func(o *options) {
		o.sigs = sigs
	}
}

// WithMetadata with service metadata.
func WithMetadata(md map[string]string) Option {
	return func(o *options) {
		o.metadata = md
	}
}

// WithEndpoint with service endpoint.
func WithEndpoint(endpoints ...*url.URL) Option {
	return func(o *options) {
		o.endpoints = endpoints
	}
}

// WithRegistry with service registry.
func WithRegistry(r registry.Registry) Option {
	return func(o *options) {
		o.registry = r
	}
}

// WithLogger .
func WithLogger(logger log.Logger) Option {
	return func(o *options) {
		o.logger = logger
	}
}

// WithServer with a server , http or grpc
func WithServer(srv ...transport.Server) Option {
	return func(o *options) {
		o.servers = srv
	}
}

// WithRegistryTimeout with registrar timeout.
func WithRegistryTimeout(t time.Duration) Option {
	return func(o *options) {
		o.registryTimeout = t
	}
}

// WithStopTimeout with app stop timeout.
func WithStopTimeout(t time.Duration) Option {
	return func(o *options) {
		o.stopTimeout = t
	}
}
