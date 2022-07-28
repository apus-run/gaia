package grpc

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials"
	grpcInsecure "google.golang.org/grpc/credentials/insecure"

	"github.com/apus-run/gaia/log"
	"github.com/apus-run/gaia/middleware"
	"github.com/apus-run/gaia/registry"
	"github.com/apus-run/gaia/transport/grpc/resolver/discovery"
)

// Client is gRPC Client
type Client struct {
	endpoint     string
	timeout      time.Duration
	tlsConf      *tls.Config
	discovery    registry.Discovery
	ms           []middleware.Middleware
	ints         []grpc.UnaryClientInterceptor
	grpcOpts     []grpc.DialOption
	balancerName string

	log *log.Helper
}

// defaultClient return a default config server
func defaultClient() *Client {
	return &Client{
		timeout:      2000 * time.Millisecond,
		balancerName: roundrobin.Name,
		log:          log.NewHelper(log.GetLogger()),
	}
}

// Dial returns a GRPC connection.
func Dial(ctx context.Context, opts ...ClientOption) (*grpc.ClientConn, error) {
	return dial(ctx, false, opts...)
}

// DialInsecure returns an insecure GRPC connection.
func DialInsecure(ctx context.Context, opts ...ClientOption) (*grpc.ClientConn, error) {
	return dial(ctx, true, opts...)
}

func dial(ctx context.Context, insecure bool, opts ...ClientOption) (*grpc.ClientConn, error) {
	c := defaultClient()

	for _, o := range opts {
		o(c)
	}

	ints := []grpc.UnaryClientInterceptor{
		c.unaryClientInterceptor(c.ms, c.timeout),
	}
	if len(c.ints) > 0 {
		ints = append(ints, c.ints...)
	}
	grpcOpts := []grpc.DialOption{
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"loadBalancingConfig": [{"%s":{}}]}`, c.balancerName)),
		grpc.WithChainUnaryInterceptor(ints...),
	}
	if c.discovery != nil {
		grpcOpts = append(grpcOpts, grpc.WithResolvers(discovery.NewBuilder(c.discovery, discovery.WithInsecure(insecure))))
	}
	if insecure {
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(grpcInsecure.NewCredentials()))
	}
	if c.tlsConf != nil {
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(credentials.NewTLS(c.tlsConf)))
	}
	if len(c.grpcOpts) > 0 {
		grpcOpts = append(grpcOpts, c.grpcOpts...)
	}

	return grpc.DialContext(ctx, c.endpoint, grpcOpts...)
}
