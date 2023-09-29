package http

import (
	"github.com/apus-run/gaia/internal/matcher"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/apus-run/gaia/internal/tls"
	"github.com/apus-run/gaia/middleware"
)

// Server is an HTTP server wrapper.
type Server struct {
	*http.Server
	lis          net.Listener
	tlsConf      *tls.Config
	network      string
	address      string
	readTimeout  time.Duration
	writeTimeout time.Duration
	endpoint     *url.URL

	filters    []FilterFunc
	middleware matcher.Matcher

	err error
}

// defaultServer return a default config server
func defaultServer() *Server {
	return &Server{
		network:      "tcp",
		address:      ":0",
		readTimeout:  1 * time.Second,
		writeTimeout: 1 * time.Second,
		middleware:   matcher.New(),
	}
}

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

// Timeout with server timeout.
func Timeout(timeout time.Duration) ServerOption {
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
		o.middleware.Use(m...)
	}
}

// Filter with HTTP middleware option.
func Filter(filters ...FilterFunc) ServerOption {
	return func(o *Server) {
		o.filters = filters
	}
}

// Listener with server lis
func Listener(lis net.Listener) ServerOption {
	return func(s *Server) {
		s.lis = lis
	}
}

// TLSConfig with TLS config.
func TLSConfig(c *tls.Config) ServerOption {
	return func(o *Server) {
		o.tlsConf = c
	}
}
