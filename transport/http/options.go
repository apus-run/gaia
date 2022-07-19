package http

import (
	"net"
	"time"

	"github.com/apus-run/gaia/log"
	"github.com/apus-run/gaia/middleware"
)

// ServerOption is an HTTP server option.
type ServerOption func(*Server)

// WithNetwork with server network.
func WithNetwork(network string) ServerOption {
	return func(s *Server) {
		s.network = network
	}
}

// WithAddress with server address.
func WithAddress(addr string) ServerOption {
	return func(s *Server) {
		s.address = addr
	}
}

// WithWriteTimeout with server timeout.
func WithWriteTimeout(timeout time.Duration) ServerOption {
	return func(s *Server) {
		s.writeTimeout = timeout
	}
}

// WithReadTimeout with server timeout.
func WithReadTimeout(timeout time.Duration) ServerOption {
	return func(s *Server) {
		s.readTimeout = timeout
	}
}

// WithMiddleware with service middleware option.
func WithMiddleware(m ...middleware.Middleware) ServerOption {
	return func(o *Server) {
		o.ms = m
	}
}

// WithListener with server lis
func WithListener(lis net.Listener) ServerOption {
	return func(s *Server) {
		s.lis = lis
	}
}

// WithLogger with server logger.
func WithLogger(logger log.Logger) ServerOption {
	return func(s *Server) {
		s.log = log.NewHelper(logger)
	}
}
