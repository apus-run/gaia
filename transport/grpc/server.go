package grpc

import (
	"context"
	"net"
	"net/url"

	"github.com/apus-run/sea-kit/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	hapi "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/apus-run/gaia/internal/endpoint"
	"github.com/apus-run/gaia/internal/host"
	"github.com/apus-run/gaia/middleware"
	"github.com/apus-run/gaia/transport"
	"google.golang.org/grpc/admin"
)

var (
	_ transport.Server     = (*Server)(nil)
	_ transport.Endpointer = (*Server)(nil)
)

// NewServer creates a gRPC server by options.
func NewServer(opts ...ServerOption) *Server {
	srv := defaultServer()

	for _, o := range opts {
		o(srv)
	}
	unaryInterceptor := []grpc.UnaryServerInterceptor{
		srv.unaryServerInterceptor(),
	}
	streamInterceptor := []grpc.StreamServerInterceptor{
		srv.streamServerInterceptor(),
	}
	if len(srv.unaryInterceptor) > 0 {
		unaryInterceptor = append(unaryInterceptor, srv.unaryInterceptor...)
	}
	if len(srv.streamInterceptor) > 0 {
		streamInterceptor = append(streamInterceptor, srv.streamInterceptor...)
	}
	grpcOpts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(unaryInterceptor...),
		grpc.ChainStreamInterceptor(streamInterceptor...),
	}
	if srv.tlsConf != nil {
		t, err := srv.tlsConf.Config()
		if err != nil {
			log.Errorf("TLS Config Error - %v", err)
		}
		grpcOpts = append(grpcOpts, grpc.Creds(credentials.NewTLS(t)))
	}
	if len(srv.grpcOpts) > 0 {
		grpcOpts = append(grpcOpts, srv.grpcOpts...)
	}

	srv.Server = grpc.NewServer(grpcOpts...)
	// internal register
	if !srv.customHealth {
		hapi.RegisterHealthServer(srv.Server, srv.health)
	}

	// register reflection and the interface can be debugged through the grpcurl tool
	// https://github.com/grpc/grpc-go/blob/master/Documentation/server-reflection-tutorial.md#enable-server-reflection
	// see https://github.com/fullstorydev/grpcurl
	reflection.Register(srv.Server)
	// admin register
	srv.adminClean, _ = admin.Register(srv.Server)
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

// Endpoint return a real address to registry endpoint.
// examples:
//
//	grpc://127.0.0.1:9000?isSecure=false
func (s *Server) Endpoint() (*url.URL, error) {
	if err := s.listenAndEndpoint(); err != nil {
		return nil, s.err
	}
	return s.endpoint, nil
}

// Start start the gRPC server.
func (s *Server) Start(ctx context.Context) error {
	if err := s.listenAndEndpoint(); err != nil {
		return s.err
	}
	s.ctx = ctx
	log.Infof("[gRPC] server listening on: %s", s.lis.Addr().String())
	s.health.Resume()
	return s.Serve(s.lis)
}

// Stop stop the gRPC server.
func (s *Server) Stop(ctx context.Context) error {
	if s.adminClean != nil {
		s.adminClean()
	}
	s.health.Shutdown()
	s.GracefulStop()
	log.Info("[gRPC] server stopping")
	return nil
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
		s.endpoint = endpoint.NewEndpoint(endpoint.Scheme("grpc", s.tlsConf != nil), addr)
	}
	return s.err
}
