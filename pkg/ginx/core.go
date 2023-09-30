package ginx

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"google.golang.org/grpc/status"

	"github.com/apus-run/gaia/middleware"
	"github.com/apus-run/gaia/pkg/errcode"
	thttp "github.com/apus-run/gaia/transport/http"
	httpStatus "github.com/apus-run/gaia/transport/http/status"
)

const (
	// RequestId
	requestIdFieldKey = "REQUEST_ID"
	// AcceptLanguageHeaderName represents the header name of accept language
	AcceptLanguageHeaderName = "Accept-Language"
	// ClientTimezoneOffsetHeaderName represents the header name of client timezone offset
	ClientTimezoneOffsetHeaderName = "X-Timezone-Offset"
)
const (
	// CodeOK means a successful response
	CodeOK = 0
	// CodeErr means a failure response
	CodeErr = 1
)

// Result defines HTTP JSON response
type Result struct {
	Code    int      `json:"code"`
	Msg     string   `json:"msg"`
	Data    any      `json:"data"`
	Details []string `json:"details,omitempty"`
}

// Context a wrapper of gin.Context
type Context struct {
	*gin.Context
}

// WrapContext returns a context wrapped by this file
func WrapContext(c *gin.Context) *Context {
	return &Context{
		Context: c,
	}
}

// HandlerFunc defines the handler to wrap gin.Context
type HandlerFunc func(c *Context)

// Handle convert HandlerFunc to gin.HandlerFunc
func Handle(h HandlerFunc) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		c := WrapContext(ginCtx)
		h(c)
	}
}

// JSON returns JSON response
// e.x. {"code":<code>, "msg":<msg>, "data":<data>, "details":<details>}
func (c *Context) JSON(httpStatus int, resp Result) {
	c.Context.JSON(httpStatus, resp)
}

// JSONOK returns JSON response with successful business code and data
// e.x. {"code":0, "msg":"成功", "data":<data>}
func (c *Context) JSONOK(msg string, data any) {
	j := new(Result)
	j.Code = CodeOK
	j.Msg = msg

	switch d := data.(type) {
	case error:
		j.Data = d.Error()
	case nil:
		j.Data = gin.H{}
	default:
		j.Data = data
	}

	c.Context.JSON(http.StatusOK, j)
	return
}

// JSONE returns JSON response with failure business code ,msg and data
// e.x. {"code":<code>, "msg":<msg>, "data":<data>}
func (c *Context) JSONE(code int, msg string, data any) {
	j := new(Result)
	j.Code = code
	j.Msg = msg
	switch d := data.(type) {
	case error:
		j.Data = d.Error()
	case nil:
		j.Data = gin.H{}
	default:
		j.Data = data
	}

	c.Context.JSON(http.StatusOK, j)
	return
}

func (c *Context) Success(data any) {
	j := new(Result)
	j.Code = errcode.Success.Code()
	j.Msg = errcode.Success.Msg()

	if data == nil {
		j.Data = gin.H{}
	} else {
		j.Data = data
	}

	c.Context.JSON(http.StatusOK, j)
	return
}

func (c *Context) Error(err error) {
	if err == nil {
		c.JSON(http.StatusOK, Result{
			Code: errcode.Success.Code(),
			Msg:  errcode.Success.Msg(),
			Data: gin.H{},
		})
		return
	}

	if v, ok := err.(*errcode.Error); ok {
		response := Result{
			Code:    v.Code(),
			Msg:     v.Msg(),
			Data:    gin.H{},
			Details: []string{},
		}
		details := v.Details()
		if len(details) > 0 {
			response.Details = details
		}
		c.JSON(errcode.ToHTTPStatusCode(v.Code()), response)
		return
	} else {
		// receive gRPC error
		if st, ok := status.FromError(err); ok {
			response := Result{
				Code:    int(st.Code()),
				Msg:     st.Message(),
				Data:    gin.H{},
				Details: []string{},
			}
			details := st.Details()
			if len(details) > 0 {
				for _, v := range details {
					response.Details = append(response.Details, cast.ToString(v))
				}
			}
			// https://httpstatus.in/
			// https://github.com/grpc-ecosystem/grpc-gateway/blob/master/runtime/errors.go#L15
			// https://github.com/googleapis/googleapis/blob/master/google/rpc/code.proto
			c.JSON(httpStatus.FromGRPCCode(st.Code()), response)
			return
		}
	}
}

// RouteNotFound 未找到相关路由
func (c *Context) RouteNotFound() {
	c.String(http.StatusNotFound, "the route not found")
}

// GetClientLocale returns the client locale name
func (c *Context) GetClientLocale() string {
	value := c.GetHeader(AcceptLanguageHeaderName)

	return value
}

// GetClientTimezoneOffset returns the client timezone offset
func (c *Context) GetClientTimezoneOffset() (int16, error) {
	value := c.GetHeader(ClientTimezoneOffsetHeaderName)
	offset, err := strconv.Atoi(value)

	if err != nil {
		return 0, err
	}

	return int16(offset), nil
}

// SetRequestId sets the given request id to context
func (c *Context) SetRequestId(requestId string) {
	c.Set(requestIdFieldKey, requestId)
}

// GetRequestId returns the current request id
func (c *Context) GetRequestId() string {
	requestId, exists := c.Get(requestIdFieldKey)

	if !exists {
		return ""
	}

	return requestId.(string)
}

// healthCheckResponse 健康检查响应结构体
type healthCheckResponse struct {
	Status   string `json:"status"`
	Hostname string `json:"hostname"`
}

// HealthCheck will return OK if the underlying BoltDB is healthy.
// At least healthy enough for demoing purposes.
func HealthCheck(c *gin.Context) {
	name, err := os.Hostname()
	if err != nil {
		name = "unknown"
	}
	c.JSON(http.StatusOK, healthCheckResponse{Status: "UP", Hostname: name})
}

type ginKey struct{}

// NewGinContext returns a new Context that carries gin.Context value.
func NewGinContext(ctx context.Context, c *gin.Context) context.Context {
	return context.WithValue(ctx, ginKey{}, c)
}

// FromGinContext returns the gin.Context value stored in ctx, if any.
func FromGinContext(ctx context.Context) (c *gin.Context, ok bool) {
	c, ok = ctx.Value(ginKey{}).(*gin.Context)
	return
}

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
		if ginCtx, ok := FromGinContext(ctx); ok {
			thttp.SetOperation(ctx, ginCtx.FullPath())
		}
		next(c.Request.Context(), c.Request)
	}
}
