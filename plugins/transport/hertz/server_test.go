package hertz

import (
	"context"
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
)

func TestServer(t *testing.T) {
	ctx := context.Background()

	srv := NewServer(
		WithAddress("127.0.0.1:8800"),
	)

	srv.GET("/login/*param", func(ctx context.Context, c *app.RequestContext) {
		if len(c.Params.ByName("param")) > 1 {
			c.AbortWithStatus(404)
			return
		}
		c.String(200, "Hello World!")
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
