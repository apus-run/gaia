package http

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/url"

	"github.com/apus-run/sea-kit/log"

	"github.com/apus-run/gaia/internal/endpoint"
	"github.com/apus-run/gaia/internal/host"
	"github.com/apus-run/gaia/middleware"
	"github.com/apus-run/gaia/transport"
)

var (
	_ transport.Server     = (*Server)(nil)
	_ transport.Endpointer = (*Server)(nil)
	_ http.Handler         = (*Server)(nil)
)

// NewServer creates an HTTP server by options.
func NewServer(opts ...ServerOption) *Server {
	srv := defaultServer()
	// apply options
	for _, o := range opts {
		o(srv)
	}

	// TODO: must set server
	if srv.tlsConf != nil {
		t, err := srv.tlsConf.Config()
		if err != nil {
			log.Errorf("TLS Config Error - %v", err)
		}
		if err == nil {
			srv.Server = &http.Server{
				Handler:      srv,
				ReadTimeout:  srv.readTimeout,
				WriteTimeout: srv.writeTimeout,
				TLSConfig:    t,
			}
		}
	}

	srv.Server = &http.Server{
		Handler:      srv,
		ReadTimeout:  srv.readTimeout,
		WriteTimeout: srv.writeTimeout,
	}

	return srv
}

// Use uses a service middleware with selector.
// selector:
//   - '/*'
//   - '/helloworld.v1.Greeter/*'
//   - '/helloworld.v1.Greeter/SayHello'
func (s *Server) Use(selector string, m ...middleware.Middleware) {
	s.middleware.Add(selector, m...)
}

func (s *Server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	s.Handler.ServeHTTP(res, req)
}

// Endpoint return a real address to registry endpoint.
// examples:
//
//	https://127.0.0.1:8000
//	Legacy: http://127.0.0.1:8000?isSecure=false
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

	var err error
	if s.tlsConf != nil {
		log.Infof("[HTTPS] server is listening on: %s", s.lis.Addr().String())
		err = s.ServeTLS(s.lis, s.tlsConf.Cert, s.tlsConf.Key)
	} else {
		log.Infof("[HTTP] server is listening on: %s", s.lis.Addr().String())
		err = s.Serve(s.lis)
	}
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

// Stop stop the HTTP server.
func (s *Server) Stop(ctx context.Context) error {
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
