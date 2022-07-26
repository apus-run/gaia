package http

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/apus-run/gaia/internal/endpoint"
	"github.com/apus-run/gaia/internal/host"
	"github.com/apus-run/gaia/log"
	"github.com/apus-run/gaia/middleware"
	"github.com/apus-run/gaia/transport"
)

var (
	_ transport.Server     = (*Server)(nil)
	_ transport.Endpointer = (*Server)(nil)
	_ http.Handler         = (*Server)(nil)
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
	ms           []middleware.Middleware
	err          error

	log *log.Helper
}

// defaultServer return a default config server
func defaultServer() *Server {
	return &Server{
		network:      "tcp",
		address:      ":0",
		readTimeout:  1 * time.Second,
		writeTimeout: 1 * time.Second,

		log: log.NewHelper(log.GetLogger()),
	}
}

// NewServer creates an HTTP server by options.
func NewServer(opts ...ServerOption) *Server {
	srv := defaultServer()
	// apply options
	for _, o := range opts {
		o(srv)
	}

	// NOTE: must set server
	srv.Server = &http.Server{
		Handler:      srv,
		ReadTimeout:  srv.readTimeout,
		WriteTimeout: srv.writeTimeout,
		TLSConfig:    srv.tlsConf,
	}
	return srv
}

func (s *Server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	s.Handler.ServeHTTP(res, req)
}

// Endpoint return a real address to registry endpoint.
// examples:
//   https://127.0.0.1:8000
//   Legacy: http://127.0.0.1:8000?isSecure=false
func (s *Server) Endpoint() (*url.URL, error) {
	if err := s.listenAndEndpoint(); err != nil {
		return nil, err
	}
	return s.endpoint, nil
}

// Start start the HTTP server.
func (s *Server) Start(ctx context.Context) error {
	if err := s.listenAndEndpoint(); err != nil {
		return err
	}
	s.BaseContext = func(net.Listener) context.Context {
		return ctx
	}
	log.Infof("[HTTP] server is listening on: %s", s.lis.Addr().String())
	var err error
	if s.tlsConf != nil {
		err = s.ServeTLS(s.lis, "", "")
	} else {
		err = s.Serve(s.lis)
	}
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

// Stop stop the HTTP server.
func (s *Server) Stop() error {
	log.Infof("[HTTP] server is stopping")
	return s.Close()
}

func (s *Server) GracefullyStop(ctx context.Context) error {
	log.Infof("[HTTP] server is stopping")
	return s.Shutdown(ctx)
}

// Health 心跳检测
func (s *Server) Health() bool {
	if s.lis == nil {
		return false
	}

	conn, err := s.lis.Accept()
	if err != nil {
		return false
	}

	er := conn.Close()
	if er != nil {
		return false
	}
	return true
}

func (s *Server) listenAndEndpoint() error {
	if s.lis == nil {
		lis, err := net.Listen(s.network, s.address)
		if err != nil {
			s.err = err
			return err
		}
		s.lis = lis
	}
	if s.endpoint == nil {
		addr, err := host.Extract(s.address, s.lis)
		if err != nil {
			s.err = err
			return err
		}
		s.endpoint = endpoint.NewEndpoint(endpoint.Scheme("http", s.tlsConf != nil), addr)
	}
	return s.err
}
