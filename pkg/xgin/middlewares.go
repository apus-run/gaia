package xgin

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/apus-run/gaia/middleware"
)

// Middlewares return middlewares wrapper
func Middlewares(m ...middleware.Middleware) gin.HandlerFunc {
	chain := middleware.Chain(m...)
	return func(c *gin.Context) {
		next := func(ctx context.Context, req interface{}) (interface{}, error) {
			c.Next()
			var err error
			if c.Writer.Status() >= http.StatusBadRequest {
				err = fmt.Errorf("error: code = %d msg = %s data = %s ", c.Writer.Status(), "", "")
			}
			return c.Writer, err
		}
		next = chain(next)
		ctx := NewGinContext(c.Request.Context(), c)
		c.Request = c.Request.WithContext(ctx)
		next(c.Request.Context(), c.Request)
	}
}
