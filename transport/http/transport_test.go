package http

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"github.com/apus-run/gaia/transport"
)

func TestTransport_Kind(t *testing.T) {
	o := &Transport{}
	if !reflect.DeepEqual(transport.KindHTTP, o.Kind()) {
		t.Errorf("expect %v, got %v", transport.KindHTTP, o.Kind())
	}
}

func TestTransport_Endpoint(t *testing.T) {
	v := "hello"
	o := &Transport{endpoint: v}
	if !reflect.DeepEqual(v, o.Endpoint()) {
		t.Errorf("expect %v, got %v", v, o.Endpoint())
	}
}

func TestTransport_Operation(t *testing.T) {
	v := "hello"
	o := &Transport{operation: v}
	if !reflect.DeepEqual(v, o.Operation()) {
		t.Errorf("expect %v, got %v", v, o.Operation())
	}
}

func TestTransport_Request(t *testing.T) {
	v := &http.Request{}
	o := &Transport{request: v}
	if !reflect.DeepEqual(v, o.Request()) {
		t.Errorf("expect %v, got %v", v, o.Request())
	}
}

func TestTransport_PathTemplate(t *testing.T) {
	v := "template"
	o := &Transport{pathTemplate: v}
	if !reflect.DeepEqual(v, o.PathTemplate()) {
		t.Errorf("expect %v, got %v", v, o.PathTemplate())
	}
}

func TestSetOperation(t *testing.T) {
	tr := &Transport{}
	ctx := transport.NewServerContext(context.Background(), tr)
	SetOperation(ctx, "gaia")
	if !reflect.DeepEqual(tr.operation, "gaia") {
		t.Errorf("expect %v, got %v", "gaia", tr.operation)
	}
}
