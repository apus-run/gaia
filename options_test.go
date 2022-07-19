package gaia

import (
	"reflect"
	"testing"
)

func TestID(t *testing.T) {
	o := &options{}
	i := "123"
	WithID(i)(o)

	if !reflect.DeepEqual(i, o.id) {
		t.Fatalf("o.id:%s is not equal to v:%s", o.id, i)
	}
}

func TestName(t *testing.T) {
	o := &options{}
	n := "user-srv"
	WithName(n)(o)
	if !reflect.DeepEqual(n, o.name) {
		t.Fatalf("o.name:%s is not equal to v:%s", o.name, n)
	}
}

func TestVersion(t *testing.T) {
	o := &options{}
	v := "v1.0.0"
	WithVersion(v)(o)
	if !reflect.DeepEqual(v, o.version) {
		t.Fatalf("o.version:%s is not equal to v:%s", o.version, v)
	}
}

func TestNewOptions(t *testing.T) {
	i := "123"
	n := "user-srv"
	v := "v1.0.0"

	opts := newOptions(WithID(i), WithName(n), WithVersion(v))
	t.Logf("options: %v \n", opts)
}
