package http

import (
	"crypto/tls"
	"log"
	"net"
	"reflect"
	"testing"

	xlog "github.com/apus-run/gaia/log"
	"github.com/apus-run/gaia/middleware"
)

func TestNetwork(t *testing.T) {
	o := &Server{}
	v := "abc"
	Network(v)(o)
	if !reflect.DeepEqual(v, o.network) {
		t.Errorf("expected %v got %v", v, o.network)
	}
}

func TestAddress(t *testing.T) {
	o := &Server{}
	v := "abc"
	Address(v)(o)
	if !reflect.DeepEqual(v, o.address) {
		t.Errorf("expected %v got %v", v, o.address)
	}
}

func TestLogger(t *testing.T) {
	o := &Server{}
	v := xlog.NewStdLogger(log.Writer())
	Logger(v)(o)
	if !reflect.DeepEqual(v, o.log) {
		t.Fatalf("o.logger:%v is not equal to xlog.NewHelper(v):%v", o.log, xlog.NewHelper(v))
	}
}

func TestMiddleware(t *testing.T) {
	o := &Server{}
	v := []middleware.Middleware{
		func(middleware.Handler) middleware.Handler { return nil },
	}
	Middleware(v...)(o)
	if !reflect.DeepEqual(v, o.ms) {
		t.Errorf("expected %v got %v", v, o.ms)
	}
}

func TestTLSConfig(t *testing.T) {
	o := &Server{}
	v := &tls.Config{}
	TLSConfig(v)(o)
	if !reflect.DeepEqual(v, o.tlsConf) {
		t.Errorf("expected %v got %v", v, o.tlsConf)
	}
}

func TestListener(t *testing.T) {
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}
	s := &Server{}
	Listener(lis)(s)
	if !reflect.DeepEqual(s.lis, lis) {
		t.Errorf("expected %v got %v", lis, s.lis)
	}
	if e, err := s.Endpoint(); err != nil || e == nil {
		t.Errorf("expected not empty")
	}
}
