package hertz

import (
	"context"
	"crypto/tls"
	"net/url"
	"strings"
	"time"

	"github.com/apus-run/sea-kit/log"
	hertz "github.com/cloudwego/hertz/pkg/app/server"

	"github.com/apus-run/gaia/middleware"
	"github.com/apus-run/gaia/transport"
	thttp "github.com/apus-run/gaia/transport/http"
)

var (
	_ transport.Server     = (*Server)(nil)
	_ transport.Endpointer = (*Server)(nil)
)

type Server struct {
	*hertz.Hertz

	tlsConf *tls.Config
	timeout time.Duration
	addr    string

	err error

	filters []thttp.FilterFunc
	ms      []middleware.Middleware
}

func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		timeout: 1 * time.Second,
	}

	srv.init(opts...)

	return srv
}

func (s *Server) init(opts ...ServerOption) {
	for _, o := range opts {
		o(s)
	}

	s.Hertz = hertz.Default(hertz.WithHostPorts(s.addr), hertz.WithTLS(s.tlsConf))
}

func (s *Server) Endpoint() (*url.URL, error) {
	addr := s.addr

	prefix := "http://"
	if s.tlsConf == nil {
		if !strings.HasPrefix(addr, "http://") {
			prefix = "http://"
		}
	} else {
		if !strings.HasPrefix(addr, "https://") {
			prefix = "https://"
		}
	}
	addr = prefix + addr

	var endpoint *url.URL
	endpoint, s.err = url.Parse(addr)

	return endpoint, s.err
}

func (s *Server) Start(ctx context.Context) error {
	log.Infof("[hertz] server listening on: %s", s.addr)

	return s.Run()
}

func (s *Server) Stop(ctx context.Context) error {
	log.Info("[hertz] server stopping")
	return s.Shutdown(ctx)
}
