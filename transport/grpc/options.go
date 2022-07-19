package grpc

import (
	"crypto/tls"
	"google.golang.org/grpc"
	"net"
	"time"

	"github.com/apus-run/gaia/v1/log"
	"github.com/apus-run/gaia/v1/middleware"
)

// ServerOption is gRPC server option.
type ServerOption func(o *Server)

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

// Logger with server logger.
func Logger(logger log.Logger) ServerOption {
	return func(s *Server) {
		s.log = log.NewHelper(logger)
	}
}

// Middleware with server middleware.
func Middleware(m ...middleware.Middleware) ServerOption {
	return func(s *Server) {
		s.middleware = m
	}
}

// TLSConf with TLS config.
func TLSConf(c *tls.Config) ServerOption {
	return func(s *Server) {
		s.tlsConf = c
	}
}

// Listener with server lis
func Listener(lis net.Listener) ServerOption {
	return func(s *Server) {
		s.lis = lis
	}
}

// UnaryInterceptor returns a ServerOption that sets the UnaryServerInterceptor for the server.
func UnaryInterceptor(in ...grpc.UnaryServerInterceptor) ServerOption {
	return func(s *Server) {
		s.unaryInterceptor = in
	}
}

// StreamInterceptor returns a ServerOption that sets the StreamServerInterceptor for the server.
func StreamInterceptor(in ...grpc.StreamServerInterceptor) ServerOption {
	return func(s *Server) {
		s.streamInterceptor = in
	}
}

// GrpcOptions with grpc options.
func GrpcOptions(opts ...grpc.ServerOption) ServerOption {
	return func(s *Server) {
		s.grpcOpts = opts
	}
}

// ClientOption is gRPC client option.
type ClientOption func(o *Client)

// WithEndpoint ...
func WithEndpoint(endpoint string) ClientOption {
	return func(c *Client) {
		c.endpoint = endpoint
	}
}

// WithGrpcOptions with gRPC options.
func WithGrpcOptions(opts ...grpc.DialOption) ClientOption {
	return func(c *Client) {
		c.grpcOpts = opts
	}
}

// WithMiddleware with client middleware.
func WithMiddleware(ms ...middleware.Middleware) ClientOption {
	return func(c *Client) {
		c.ms = ms
	}
}

// WithTLSConfig with TLS config.
func WithTLSConfig(conf *tls.Config) ClientOption {
	return func(c *Client) {
		c.tlsConf = conf
	}
}

// WithUnaryInterceptor returns a DialOption that specifies the interceptor for unary RPCs.
func WithUnaryInterceptor(in ...grpc.UnaryClientInterceptor) ClientOption {
	return func(c *Client) {
		c.ints = in
	}
}

// WithBalancerName with balancer name
func WithBalancerName(name string) ClientOption {
	return func(c *Client) {
		c.balancerName = name
	}
}

// WithLogger with server logger.
func WithLogger(logger log.Logger) ClientOption {
	return func(c *Client) {
		c.log = log.NewHelper(logger)
	}
}

// WithTimeout with client timeout.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.timeout = timeout
	}
}
