package http

import (
	"net"
	"time"

	"github.com/apus-run/gaia/log"
	"github.com/apus-run/gaia/middleware"
)

// ServerOption is an HTTP server option.
type ServerOption func(*Server)

// Network with server network.
func Network(network string) ServerOption {
	return func(s *Server) {
		s.network = network
	}
}

// Address with server address.
func Address(addr string) ServerOption {
	return func(s *Server) {
		s.address = addr
	}
}

// WriteTimeout with server timeout.
func WriteTimeout(timeout time.Duration) ServerOption {
	return func(s *Server) {
		s.writeTimeout = timeout
	}
}

// ReadTimeout with server timeout.
func ReadTimeout(timeout time.Duration) ServerOption {
	return func(s *Server) {
		s.readTimeout = timeout
	}
}

// Middleware with service middleware option.
func Middleware(m ...middleware.Middleware) ServerOption {
	return func(o *Server) {
		o.ms = m
	}
}

// Listener with server lis
func Listener(lis net.Listener) ServerOption {
	return func(s *Server) {
		s.lis = lis
	}
}

// Logger with server logger.
func Logger(logger log.Logger) ServerOption {
	return func(s *Server) {
		s.log = log.NewHelper(logger)
	}
}
