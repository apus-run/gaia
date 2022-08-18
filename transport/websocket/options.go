package websocket

import (
	"crypto/tls"
	"net"
	"time"

	"github.com/apus-run/gaia/encoding"
	"github.com/apus-run/gaia/log"
)

type ServerOption func(o *Server)

func Network(network string) ServerOption {
	return func(s *Server) {
		s.network = network
	}
}

func Address(addr string) ServerOption {
	return func(s *Server) {
		s.address = addr
	}
}

func Timeout(timeout time.Duration) ServerOption {
	return func(s *Server) {
		s.timeout = timeout
	}
}

func Path(path string) ServerOption {
	return func(s *Server) {
		s.path = path
	}
}

func ConnectHandle(h ConnectHandler) ServerOption {
	return func(s *Server) {
		s.connectHandler = h
	}
}

func Logger(logger log.Logger) ServerOption {
	return func(s *Server) {
		s.log = log.NewHelper(logger, log.WithMessageKey("websocket"))
	}
}

func TLSConfig(c *tls.Config) ServerOption {
	return func(o *Server) {
		o.tlsConf = c
	}
}

func Listener(lis net.Listener) ServerOption {
	return func(s *Server) {
		s.lis = lis
	}
}

func Codec(c encoding.Codec) ServerOption {
	return func(s *Server) {
		s.Codec = c
	}
}
