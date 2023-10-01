package fasthttp

import (
	"context"
	"testing"

	"github.com/valyala/fasthttp"
)

func TestServer(t *testing.T) {
	ctx := context.Background()

	srv := NewServer(
		WithAddress(":8800"),
	)

	srv.GET("/login/*param", func(c *fasthttp.RequestCtx) {
		_, _ = c.WriteString("Hello World!")
	})

	if err := srv.Start(ctx); err != nil {
		panic(err)
	}

	defer func() {
		if err := srv.Stop(ctx); err != nil {
			t.Errorf("expected nil got %v", err)
		}
	}()
}
