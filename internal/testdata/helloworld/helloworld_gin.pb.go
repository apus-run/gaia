// Code generated protoc-gen-go-gin. DO NOT EDIT.
// protoc-gen-go-gin v1.10.1

package helloworld

import (
	context "context"
	errcode "github.com/apus-run/gaia/pkg/errcode"
	ginx "github.com/apus-run/gaia/pkg/ginx"
	gin "github.com/gin-gonic/gin"
	metadata "google.golang.org/grpc/metadata"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the eagle package it is being compiled against.

// context.
// metadata.
// gin.ginx.errcode.

type GreeterHTTPServer interface {
	SayHello(context.Context, *HelloRequest) (*HelloReply, error)
	SayHelloPost(context.Context, *HelloRequest) (*HelloReply, error)
}

func RegisterGreeterHTTPServer(r gin.IRouter, srv GreeterHTTPServer) {
	s := &Greeter{
		server: srv,
		router: r,
	}
	s.RegisterService()
}

type Greeter struct {
	server GreeterHTTPServer
	router gin.IRouter
}

func (s *Greeter) SayHello_0_HTTP_Handler(ctx *ginx.Context) {
	var in HelloRequest

	if err := ctx.ShouldBindUri(&in); err != nil {
		e := errcode.ErrInvalidParam.WithDetails(err.Error())
		ctx.Error(e)
		return
	}

	md := metadata.New(nil)
	for k, v := range ctx.Request.Header {
		md.Set(k, v...)
	}
	newCtx := metadata.NewIncomingContext(ctx, md)
	out, err := s.server.(GreeterHTTPServer).SayHello(newCtx, &in)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.Success(out)
}

func (s *Greeter) SayHelloPost_0_HTTP_Handler(ctx *ginx.Context) {
	var in HelloRequest

	if err := ctx.ShouldBindJSON(&in); err != nil {
		e := errcode.ErrInvalidParam.WithDetails(err.Error())
		ctx.Error(e)
		return
	}

	md := metadata.New(nil)
	for k, v := range ctx.Request.Header {
		md.Set(k, v...)
	}
	newCtx := metadata.NewIncomingContext(ctx, md)
	out, err := s.server.(GreeterHTTPServer).SayHelloPost(newCtx, &in)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.Success(out)
}

func (s *Greeter) RegisterService() {
	s.router.Handle("GET", "/hello/:name", ginx.Handle(s.SayHello_0_HTTP_Handler))
	s.router.Handle("POST", "/hello", ginx.Handle(s.SayHelloPost_0_HTTP_Handler))
}
