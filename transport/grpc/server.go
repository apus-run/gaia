package grpc

import (
	"context"
	"crypto/tls"
	"net"
	"net/url"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/apus-run/gaia/internal/endpoint"
	"github.com/apus-run/gaia/internal/host"
	"github.com/apus-run/gaia/log"
	"github.com/apus-run/gaia/middleware"
	"github.com/apus-run/gaia/transport"
)

var (
	_ transport.Server     = (*Server)(nil)
	_ transport.Endpointer = (*Server)(nil)
)

// Server is a gRPC server wrapper.
type Server struct {
	*grpc.Server
	ctx               context.Context
	tlsConf           *tls.Config
	lis               net.Listener
	err               error
	network           string
	address           string
	endpoint          *url.URL
	middleware        []middleware.Middleware
	unaryInterceptor  []grpc.UnaryServerInterceptor
	streamInterceptor []grpc.StreamServerInterceptor
	grpcOpts          []grpc.ServerOption
	health            *health.Server

	log *log.Helper
}

// defaultServer return a default config server
func defaultServer() *Server {
	return &Server{
		ctx:     context.Background(),
		network: "tcp",
		address: ":0",
		health:  health.NewServer(),
		log:     log.NewHelper(log.GetLogger()),
	}
}

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
		grpcOpts = append(grpcOpts, grpc.Creds(credentials.NewTLS(srv.tlsConf)))
	}
	if len(srv.grpcOpts) > 0 {
		grpcOpts = append(grpcOpts, srv.grpcOpts...)
	}

	srv.Server = grpc.NewServer(grpcOpts...)

	// listen and endpoint
	srv.err = srv.listenAndEndpoint()

	// see https://github.com/grpc/grpc/blob/master/doc/health-checking.md
	grpc_health_v1.RegisterHealthServer(srv.Server, srv.health)

	// register reflection and the interface can be debugged through the grpcurl tool
	// https://github.com/grpc/grpc-go/blob/master/Documentation/server-reflection-tutorial.md#enable-server-reflection
	// see https://github.com/fullstorydev/grpcurl
	reflection.Register(srv.Server)

	return srv
}

// Endpoint return a real address to registry endpoint.
// examples:
//   grpc://127.0.0.1:9000?isSecure=false
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
func (s *Server) Stop() error {
	s.health.Shutdown()
	s.Server.Stop()
	s.log.Info("[gRPC] server stopping")
	return nil
}

func (s *Server) GracefullyStop(ctx context.Context) error {
	s.health.Shutdown()
	s.Server.GracefulStop()
	s.log.Info("[gRPC] server graceful stopping")
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
