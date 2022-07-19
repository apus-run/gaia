package grpc

import (
	"context"
	"google.golang.org/grpc"
	"time"

	"github.com/apus-run/gaia/v1/middleware"
	ic "github.com/apus-run/gaia/v1/pkg/context"
)

// wrappedStream is rewrite grpc stream's context
type wrappedStream struct {
	grpc.ServerStream
	ctx context.Context
}

func NewWrappedStream(ctx context.Context, stream grpc.ServerStream) grpc.ServerStream {
	return &wrappedStream{
		ServerStream: stream,
		ctx:          ctx,
	}
}

func (w *wrappedStream) Context() context.Context {
	return w.ctx
}

// unaryServerInterceptor is a gRPC unary server interceptor
func (s *Server) unaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx, cancel := ic.Merge(ctx, s.ctx)
		defer cancel()

		// md, ok := metadata.FromIncomingContext(ctx)

		h := func(ctx context.Context, req interface{}) (interface{}, error) {
			return handler(ctx, req)
		}
		if len(s.middleware) > 0 {
			h = middleware.Chain(s.middleware...)(h)
		}

		reply, err := h(ctx, req)

		return reply, err
	}
}

// streamServerInterceptor is a gRPC stream server interceptor
func (s *Server) streamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx, cancel := ic.Merge(ss.Context(), s.ctx)
		defer cancel()

		// md, ok := metadata.FromIncomingContext(ctx)

		ws := NewWrappedStream(ctx, ss)

		err := handler(srv, ws)

		return err
	}
}

// unaryClientInterceptor client unary interceptor
func (c *Client) unaryClientInterceptor(ms []middleware.Middleware, timeout time.Duration) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if timeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, timeout)
			defer cancel()
		}

		h := func(ctx context.Context, req interface{}) (interface{}, error) {
			return reply, invoker(ctx, method, req, reply, cc, opts...)
		}

		if len(ms) > 0 {
			h = middleware.Chain(ms...)(h)
		}

		_, err := h(ctx, req)

		return err
	}
}
